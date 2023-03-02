package flip

import (
	"fmt"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Flips flips (mirrors) a image vertically or horizontally.
type Flip struct {
	Direction string `json:"direction,omitempty"`
}

// UnmarshalCaddyfile configures the Flip instance.
//
// Syntax:
//
//	flip <h|v>
//
// Parameters:
//
// h|v determines if the image flipped horizontally or vertically.
func (f *Flip) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() > 1 {
		return imagefilter.ErrTooManyArgs
	}
	if !d.NextArg() {
		return imagefilter.ErrTooFewArgs
	}

	f.Direction = d.Val()

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Flip) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	direction := repl.ReplaceAll(f.Direction, "")

	switch direction {
	case "h":
		return imaging.FlipH(img), nil
	case "v":
		return imaging.FlipV(img), nil
	default:
		return nil, fmt.Errorf("unknown flip direction %s", direction)
	}
}

// CaddyModule returns the Caddy module information.
func (Flip) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.flip",
		New: func() caddy.Module { return new(Flip) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Flip{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Flip)(nil)
	_ caddyfile.Unmarshaler = (*Flip)(nil)
)
