package imagefilter

import (
	"encoding/json"
	"fmt"
	"image"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
)

var (
	registeredFilterMu sync.RWMutex
	registeredFilter   = make(map[string]FilterFactory)
)

// ImageFilter is a caddy module that can apply image filters to images from the filesystem at
// runtime. It should be used together with a cache module, so filters don't have to be applied
// repeatedly because it's an expensive operation.
type ImageFilter struct {
	// Filters is a map of initialized image filters. Keys have the form
	// "<position>_<image filter name>", where <position> specifies the order in which the image
	// filters will be applied.
	Filters filters `json:"filters,omitempty"`

	// FilterOrder is a slice of strings in the form "<position>_<image filter name>". Each entry
	// should have a corresponding entry in the Filters map.
	FilterOrder []string `json:"filterOrder,omitempty"`

	// Root ise path to the root of the site. Default is `{http.vars.root}` if set, or current
	// working directory otherwise.
	Root string `json:"root,omitempty"`

	logger *zap.Logger
}

type filters map[string]Filter

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
			case "root":
				if !h.Args(&img.Root) {
					return nil, h.ArgErr()
				}

			default:
				factory, ok := registeredFilter[h.Val()]
				if !ok {
					return nil, h.Errf("unrecognized subdirective or filter '%s'", h.Val())
				}
				filter, err := factory.New(h.RemainingArgs()...)
				if err != nil {
					return nil, err
				}

				filterName := fmt.Sprintf("%04d_%s", filterIndex, factory.Name())
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
	img.logger = ctx.Logger(img)

	if img.Root == "" {
		img.Root = "{http.vars.root}"
	}
	return nil
}

// Validate validates the configuration of the image filter module.
func (img *ImageFilter) Validate() error {
	// this is just a very inefficient file_server otherwise
	if len(img.FilterOrder) == 0 {
		return fmt.Errorf("no image filters to apply configured")
	}

	for _, filterName := range img.FilterOrder {
		if _, ok := img.Filters[filterName]; !ok {
			return fmt.Errorf("no image filter '%s' configured", filterName)
		}
	}

	return nil
}

// ServeHTTP looks for the file in the current root directory and applys the configured filters.
func (img *ImageFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	root := repl.ReplaceAll(img.Root, ".")
	if root == "" {
		root = "."
	}

	uri := repl.ReplaceAll(r.URL.Path, "")
	filename := filepath.Join(root, filepath.Clean("/"+uri))

	_, err := os.Stat(filename)
	if err != nil {
		return next.ServeHTTP(w, r)
	}
	file, err := os.Open(filename)
	if err != nil {
		img.logger.Warn("decoding of image failed", zap.Error(err))
		return next.ServeHTTP(w, r)
	}
	defer file.Close()

	reqImg, formatName, err := image.Decode(file)
	if err != nil {
		img.logger.Warn("decoding of image failed", zap.Error(err))
		return next.ServeHTTP(w, r)
	}

	for _, filterName := range img.FilterOrder {
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
		img.logger.Info("not supported format, falling back to jpeg", zap.String("format", formatName))
		format = imaging.JPEG
		formatName = "jpg"
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

	err = imaging.Encode(w, reqImg, format)
	if err != nil {
		img.logger.Error("failed to encode image", zap.Error(err))
	}

	return nil
}

// FilterFactory generates instances of it's corresponding image filter.
type FilterFactory interface {
	// Name retrurns the name if the filter, which is also the directive used in the image filter
	// block. It should be in lower case.
	Name() string

	// New intitialises and returns the image filter instance.
	New(...string) (Filter, error)

	// Unmarshal decodes JSON configuration and returns the corresponding image filter instance.
	Unmarshal([]byte) (Filter, error)
}

// Filter is a image filter that can be applied to an image.
type Filter interface {
	// Apply applies the image filter the an image and returns the new image.
	Apply(*caddy.Replacer, image.Image) (image.Image, error)
}

// Interface guards
var (
	_ json.Unmarshaler            = (*filters)(nil)
	_ caddy.Provisioner           = (*ImageFilter)(nil)
	_ caddy.Validator             = (*ImageFilter)(nil)
	_ caddyhttp.MiddlewareHandler = (*ImageFilter)(nil)
	_ caddy.Module                = (*ImageFilter)(nil)
)