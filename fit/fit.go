package fit

import (
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Fit scales a image to fit to the specified maximum width and height using a linear filter, the
// image aspect ratio is preserved. If the image already fits inside the bounds, nothing will be
// done.
type Fit struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// UnmarshalCaddyfile configures the Fit instance.
//
// Syntax:
//
//	fit <width> <height>
//
// Parameters:
//
// width must be a positive integer and determines the maximum width.
//
// height must be a positive integer and determines the maximum height.
func (f *Fit) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
func (f *Fit) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

	if height <= 0 || width <= 0 {
		return img, fmt.Errorf("invalid width height combination %d %d", width, height)
	}

	return imaging.Fit(img, width, height, imaging.Linear), nil
}

// CaddyModule returns the Caddy module information.
func (Fit) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.fit",
		New: func() caddy.Module { return new(Fit) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Fit{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Fit)(nil)
	_ caddyfile.Unmarshaler = (*Fit)(nil)
)
