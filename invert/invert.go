package invert

import (
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Invert produces an inverted (negated) version of the image.
type Invert struct{}

// UnmarshalCaddyfile configures the Grayscale instance.
//
// Syntax:
//
//	invert
//
// no parameters.
func (*Invert) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() > 0 {
		return imagefilter.ErrTooManyArgs
	}

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Invert) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Invert(img), nil
}

// CaddyModule returns the Caddy module information.
func (Invert) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.invert",
		New: func() caddy.Module { return new(Invert) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Invert{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Invert)(nil)
	_ caddyfile.Unmarshaler = (*Invert)(nil)
)
