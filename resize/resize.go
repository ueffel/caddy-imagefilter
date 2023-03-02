package resize

import (
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Resize can downsize images. If upsizing of an image is detected, nothing will be done and
// the input image is returned unchanged.
type Resize struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// UnmarshalCaddyfile configures the Resize instance.
//
// Syntax:
//
//	resize <width> <height>
//
// Parameters:
//
// width must be a positive integer and determines the maximum width.
//
// height must be a positive integer and determines the maximum height.
//
// Either width or height can be 0, then the image aspect ratio is preserved.
func (f *Resize) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() < 2 {
		return imagefilter.ErrTooFewArgs
	}
	if d.CountRemainingArgs() > 2 {
		return imagefilter.ErrTooManyArgs
	}

	args := d.RemainingArgs()
	f.Width = args[0]
	f.Height = args[1]

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Resize) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var width int
	widthRepl := repl.ReplaceAll(f.Width, "")
	if widthRepl == "" {
		width = 0
	} else {
		width, err = strconv.Atoi(widthRepl)
		if err != nil {
			return img, fmt.Errorf("invalid width: %w", err)
		}
	}
	var height int
	heightRepl := repl.ReplaceAll(f.Height, "")
	if heightRepl == "" {
		height = 0
	} else {
		height, err = strconv.Atoi(heightRepl)
		if err != nil {
			return img, fmt.Errorf("invalid height: %w", err)
		}
	}

	if height < 0 || width < 0 || height == 0 && width == 0 {
		return img, fmt.Errorf("invalid width height combination %d %d", width, height)
	}

	// no upsizing
	if height == 0 && img.Bounds().Dx() <= width ||
		width == 0 && img.Bounds().Dy() <= height ||
		img.Bounds().Dx() <= width && img.Bounds().Dy() <= height {
		return img, nil
	}

	return imaging.Resize(img, width, height, imaging.Linear), nil
}

// CaddyModule returns the Caddy module information.
func (Resize) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.resize",
		New: func() caddy.Module { return new(Resize) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Resize{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Resize)(nil)
	_ caddyfile.Unmarshaler = (*Resize)(nil)
)
