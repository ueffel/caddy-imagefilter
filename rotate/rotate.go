package rotate

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

// Rotate rotates a image 90, 180 or 270 degrees counter-clockwise.
type Rotate struct {
	Angle string `json:"angle,omitempty"`
}

// UnmarshalCaddyfile configures the Rotate instance.
//
// Syntax:
//
//	rotate <angle>
//
// Parameters:
//
// angle is one of the following: 0, 90, 180, 270 (0 is valid, but nothing will be done to the
// image).
func (f *Rotate) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() > 1 {
		return imagefilter.ErrTooManyArgs
	}
	if !d.NextArg() {
		return imagefilter.ErrTooFewArgs
	}

	f.Angle = d.Val()

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Rotate) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	angleRepl := repl.ReplaceAll(f.Angle, "")
	angle, err := strconv.Atoi(angleRepl)
	if err != nil {
		return img, fmt.Errorf("invalid angle: %w", err)
	}

	switch angle {
	case 0:
		return img, nil
	case 90:
		return imaging.Rotate90(img), nil
	case 180:
		return imaging.Rotate180(img), nil
	case 270:
		return imaging.Rotate270(img), nil
	default:
		return nil, errors.New("invalid angle (only 0, 90, 180, 270 allowed)")
	}
}

// CaddyModule returns the Caddy module information.
func (Rotate) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.rotate",
		New: func() caddy.Module { return new(Rotate) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Rotate{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Rotate)(nil)
	_ caddyfile.Unmarshaler = (*Rotate)(nil)
)
