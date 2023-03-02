package grayscale

import (
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Grayscale produces a grayscaled version of the image.
type Grayscale struct{}

// UnmarshalCaddyfile configures the Grayscale instance.
//
// Syntax:
//
//	grayscale
//
// no parameters.
func (*Grayscale) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() > 0 {
		return imagefilter.ErrTooManyArgs
	}

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Grayscale) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Grayscale(img), nil
}

// CaddyModule returns the Caddy module information.
func (Grayscale) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.grayscale",
		New: func() caddy.Module { return new(Grayscale) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Grayscale{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Grayscale)(nil)
	_ caddyfile.Unmarshaler = (*Grayscale)(nil)
)
