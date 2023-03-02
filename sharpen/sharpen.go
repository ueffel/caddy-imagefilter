package sharpen

import (
	"errors"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Sharpen produces a sharpened version of the image.
type Sharpen struct {
	Sigma string `json:"sigma,omitempty"`
}

// UnmarshalCaddyfile configures the Sharpen instance.
//
// Syntax:
//
//	sharpen [<sigma>]
//
// Parameters:
//
// sigma must be a positive floating point number and indicates how much the image will be
// sharpened. Default is 1.
func (f *Sharpen) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() > 1 {
		return imagefilter.ErrTooManyArgs
	}

	if d.NextArg() {
		f.Sigma = d.Val()
	}
	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Sharpen) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var sigma float64
	sigmaRepl := repl.ReplaceAll(f.Sigma, "")
	if sigmaRepl == "" {
		sigma = 1
	} else {
		sigma, err = strconv.ParseFloat(sigmaRepl, 64)
		if err != nil {
			return img, fmt.Errorf("invalid sigma: %w", err)
		}
	}

	if sigma <= 0 {
		return img, errors.New("invalid sigma: cannot be less or equal 0")
	}

	return imaging.Sharpen(img, sigma), nil
}

// CaddyModule returns the Caddy module information.
func (Sharpen) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.sharpen",
		New: func() caddy.Module { return new(Sharpen) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Sharpen{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Sharpen)(nil)
	_ caddyfile.Unmarshaler = (*Sharpen)(nil)
)
