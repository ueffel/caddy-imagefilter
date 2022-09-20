package imagefilter

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	_ "golang.org/x/image/webp"
	"golang.org/x/sync/semaphore"
)

var (
	registeredFilterMu sync.RWMutex
	registeredFilter   = make(map[string]FilterFactory)

	ErrTooFewArgs  = errors.New("too few arguments")
	ErrTooManyArgs = errors.New("too many arguments")
)

// ImageFilter is a caddy module that can apply image filters to images from the filesystem at
// runtime. It should be used together with a cache module, so filters don't have to be applied
// repeatedly because it's an expensive operation.
type ImageFilter struct {
	// The file system implementation to use. By default, Caddy uses the local disk file system.
	FileSystemRaw json.RawMessage `json:"file_system,omitempty" caddy:"namespace=caddy.fs inline_key=backend"`
	fileSystem    fs.StatFS

	// Filters is a map of initialized image filters. Keys have the form
	// "<position>_<image filter name>", where <position> specifies the order in which the image
	// filters will be applied.
	Filters filters `json:"filters,omitempty"`

	logger *zap.Logger

	concurrencySemaphore *semaphore.Weighted

	// Root is the path to the root of the site. Default is `{http.vars.root}` if set, or current
	// working directory otherwise.
	Root string `json:"root,omitempty"`

	// FilterOrder is a slice of strings in the form "<position>_<image filter name>". Each entry
	// should have a corresponding entry in the Filters map.
	FilterOrder []string `json:"filter_order,omitempty"`

	encodingOpts []imaging.EncodeOption

	// JpegQuality determines the quality of jpeg encoding after the filters are applied. It ranges
	// from 1 to 100 inclusive, higher is better. Default is 75.
	JpegQuality int `json:"jpeg_quality,omitempty"`

	// PngCompression determines the compression of png images. Possible values are:
	//   * 0: Default compression
	//   * -1: no compression
	//   * -2: fastest compression
	//   * -3: best compression
	PngCompression int `json:"png_compression,omitempty"`

	// MaxConcurrent determines how many request can be served concurrently. Default is 0, which
	// means unlimited
	MaxConcurrent int64 `json:"max_concurrent,omitempty"`
}

type filters map[string]Filter

// osFS is a simple fs.StatFS implementation that uses the local file system.
type osFS struct{}

func (osFS) Open(name string) (fs.File, error)     { return os.Open(name) }
func (osFS) Stat(name string) (fs.FileInfo, error) { return os.Stat(name) }

// Register registers a filter with it's FilterFactory which is used to create instances of
// the corresponding filter.
func Register(factory FilterFactory) {
	registeredFilterMu.Lock()
	defer registeredFilterMu.Unlock()
	if registeredFilter == nil {
		panic("registeredFilter are nil!")
	}
	name := factory.Name()
	if _, dup := registeredFilter[name]; dup {
		panic(fmt.Sprintf("filter already registered '%s'", name))
	}
	registeredFilter[name] = factory
}

// init registers the caddy module and the image_filter directive.
func init() {
	httpcaddyfile.RegisterHandlerDirective("image_filter", parseCaddyfile)
	caddy.RegisterModule(ImageFilter{})
}

// UnmarshalJSON unmarshals the Filter slice.
func (fs *filters) UnmarshalJSON(data []byte) error {
	var rawFilters map[string]json.RawMessage
	err := json.Unmarshal(data, &rawFilters)
	if err != nil {
		return err
	}
	result := filters{}
	for k, v := range rawFilters {
		filterType := k[5:]
		factory, ok := registeredFilter[filterType]
		if !ok {
			return fmt.Errorf("unrecognized filter '%s'", filterType)
		}
		filter, err := factory.Unmarshal(v)
		if err != nil {
			return err
		}
		result[k] = filter
	}
	*fs = result
	return nil
}

// CaddyModule returns the Caddy module information.
func (ImageFilter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter",
		New: func() caddy.Module { return new(ImageFilter) },
	}
}

// parseCaddyfile parses the caddyfile configuration and initialises the handler.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	img := new(ImageFilter)
	filters := make(map[string]Filter)
	var filterOrder []string
	filterIndex := 0
	for h.Next() {
		if len(h.RemainingArgs()) > 0 {
			return nil, h.ArgErr()
		}

		for h.NextBlock(0) {
			switch h.Val() {
			case "fs":
				if !h.NextArg() {
					return nil, h.ArgErr()
				}
				if img.FileSystemRaw != nil {
					return nil, h.Err("file system module already specified")
				}
				name := h.Val()
				modID := "caddy.fs." + name
				unm, err := caddyfile.UnmarshalModule(h.Dispenser, modID)
				if err != nil {
					return nil, err
				}
				statFS, ok := unm.(fs.StatFS)
				if !ok {
					return nil,
						h.Errf("module %s (%T) is not a supported file system implementation (requires fs.StatFS)",
							modID,
							unm)
				}
				img.FileSystemRaw = caddyconfig.JSONModuleObject(statFS, "backend", name, nil)

			case "root":
				if !h.Args(&img.Root) {
					return nil, h.ArgErr()
				}

			case "jpeg_quality":
				args := h.RemainingArgs()
				if len(args) != 1 {
					return nil, h.ArgErr()
				}
				q, err := strconv.Atoi(args[0])
				if err != nil {
					return nil, h.Errf("invalid jpeg_quality: %w", err)
				}
				img.JpegQuality = q

			case "png_compression":
				args := h.RemainingArgs()
				if len(args) != 1 {
					return nil, h.ArgErr()
				}
				q, err := strconv.Atoi(args[0])
				if err != nil {
					return nil, h.Errf("invalid png_compression: %w", err)
				}
				img.PngCompression = q

			case "max_concurrent":
				args := h.RemainingArgs()
				if len(args) != 1 {
					return nil, h.ArgErr()
				}
				mc, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, h.Errf("invalid max_concurrent: %w", err)
				}
				img.MaxConcurrent = mc

			default:
				factory, ok := registeredFilter[h.Val()]
				if !ok {
					return nil, h.Errf("unrecognized subdirective or filter '%s'", h.Val())
				}
				factoryName := factory.Name()
				filter, err := factory.New(h.RemainingArgs()...)
				if err != nil {
					return nil, h.Errf("%s: %w", factoryName, err)
				}
				filterName := fmt.Sprintf("%04d_%s", filterIndex, factoryName)
				filters[filterName] = filter
				filterOrder = append(filterOrder, filterName)
				filterIndex++
			}
		}
	}
	img.Filters = filters
	img.FilterOrder = filterOrder

	return img, nil
}

// Provision sets up image filter module.
func (img *ImageFilter) Provision(ctx caddy.Context) error {
	img.logger = ctx.Logger()

	// establish which file system (possibly a virtual one) we'll be using
	if len(img.FileSystemRaw) > 0 {
		mod, err := ctx.LoadModule(img, "FileSystemRaw")
		if err != nil {
			return fmt.Errorf("loading file system module: %v", err)
		}
		img.fileSystem = mod.(fs.StatFS)
	}
	if img.fileSystem == nil {
		img.fileSystem = osFS{}
	}

	if img.Root == "" {
		img.Root = "{http.vars.root}"
	}

	if img.JpegQuality == 0 {
		img.JpegQuality = jpeg.DefaultQuality
	}
	img.encodingOpts = append(img.encodingOpts, imaging.JPEGQuality(img.JpegQuality))

	img.encodingOpts = append(img.encodingOpts, imaging.PNGCompressionLevel(png.CompressionLevel(img.PngCompression)))

	if img.MaxConcurrent > 0 {
		img.concurrencySemaphore = semaphore.NewWeighted(img.MaxConcurrent)
	}

	return nil
}

// Validate validates the configuration of the image filter module.
func (img *ImageFilter) Validate() error {
	// this is just a very inefficient file_server otherwise
	if len(img.FilterOrder) == 0 {
		return errors.New("no image filters to apply configured")
	}

	for i, filterName := range img.FilterOrder {
		if _, ok := img.Filters[filterName]; !ok {
			return fmt.Errorf("no image filter '%s' configured", filterName)
		}
		if i >= 9999 {
			return fmt.Errorf("too many filters")
		}
	}

	if img.JpegQuality <= 0 || img.JpegQuality > 100 {
		return errors.New("jpeg_quality must be between 1 and 100")
	}

	if img.PngCompression > 0 || img.PngCompression < -3 {
		return errors.New("png_compression must be between -3 and 0")
	}

	if img.MaxConcurrent < 0 {
		return errors.New("max_concurrent must be greater or equal 0")
	}

	return nil
}

// ServeHTTP looks for the file in the current root directory and applys the configured filters.
func (img *ImageFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	if img.concurrencySemaphore != nil {
		err := img.concurrencySemaphore.Acquire(r.Context(), 1)
		if err != nil {
			return caddyhttp.Error(http.StatusInternalServerError, err)
		}
		defer img.concurrencySemaphore.Release(1)
	}

	root := repl.ReplaceAll(img.Root, ".")
	if root == "" {
		root = "."
	}

	uri := repl.ReplaceAll(r.URL.Path, "")
	filename := filepath.Join(root, filepath.Clean("/"+uri))

	_, err := img.fileSystem.Stat(filename)
	if err != nil {
		return caddyhttp.Error(http.StatusNotFound, err)
	}
	file, err := img.fileSystem.Open(filename)
	if err != nil {
		return caddyhttp.Error(http.StatusNotFound, err)
	}
	defer file.Close()

	reqImg, formatName, err := image.Decode(file)
	if err != nil {
		img.logger.Warn("decoding of image failed", zap.Error(err))
		return caddyhttp.Error(http.StatusUnsupportedMediaType, err)
	}
	file.Close()

	for _, filterName := range img.FilterOrder {
		if r.Context().Err() != nil {
			return r.Context().Err()
		}
		filter := img.Filters[filterName]
		newImg, err := filter.Apply(repl, reqImg)
		if err != nil {
			img.logger.Warn("error applying image filter: ", zap.String("name", filterName), zap.Error(err))
			continue
		}
		reqImg = newImg
	}

	format, err := imaging.FormatFromExtension(formatName)
	if err != nil {
		img.logger.Info("not supported format, falling back to png", zap.String("format", formatName))
		format = imaging.PNG
		formatName = "png"
	}

	if w.Header().Get("Content-Type") == "" {
		mtyp := mime.TypeByExtension("." + formatName)
		if mtyp == "" {
			// do not allow Go to sniff the content-type; see
			// https://www.youtube.com/watch?v=8t8JYpt0egE
			w.Header()["Content-Type"] = nil
		} else {
			w.Header().Set("Content-Type", mtyp)
		}
	}

	if r.Context().Err() != nil {
		return r.Context().Err()
	}

	err = imaging.Encode(w, reqImg, format, img.encodingOpts...)
	if err != nil {
		img.logger.Error("failed to encode image", zap.Error(err))
	}

	return nil
}

// FilterFactory generates instances of it's corresponding image filter.
type FilterFactory interface {
	// Name returns the name of the filter, which is also the directive used in the image filter
	// block. It should be in lower case.
	Name() string

	// New initialises and returns the image filter instance.
	New(...string) (Filter, error)

	// Unmarshal decodes JSON configuration and returns the corresponding image filter instance.
	Unmarshal([]byte) (Filter, error)
}

// Filter is a image filter that can be applied to an image.
type Filter interface {
	// Apply applies the image filter to an image and returns the new image.
	Apply(*caddy.Replacer, image.Image) (image.Image, error)
}

// Interface guards.
var (
	_ json.Unmarshaler            = (*filters)(nil)
	_ caddy.Provisioner           = (*ImageFilter)(nil)
	_ caddy.Validator             = (*ImageFilter)(nil)
	_ caddyhttp.MiddlewareHandler = (*ImageFilter)(nil)
	_ caddy.Module                = (*ImageFilter)(nil)
)
