package crop

import (
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Crop produces a cropped image as rectangular region of a specific size.
type Crop struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Anchor string `json:"anchor,omitempty"`
}

// UnmarshalCaddyfile configures Crop instance.
//
// Syntax:
//
//	crop <width> <height> [<anchor>]
//
// Parameters:
//
// width must be a positive integer and determines the width of the cropped image.
//
// height must be a positive integer and determines the height of the cropped image.
//
// anchor determines the anchor point of the rectangular region that is cut out. Possible values
// are: center, topleft, top, topright, left, right, bottomleft, bottom, bottomright.
// Default is center.
func (f *Crop) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if d.CountRemainingArgs() < 2 {
		return imagefilter.ErrTooFewArgs
	}
	if d.CountRemainingArgs() > 3 {
		return imagefilter.ErrTooManyArgs
	}

	args := d.RemainingArgs()
	f.Width = args[0]
	f.Height = args[1]

	if len(args) < 3 {
		f.Anchor = "center"
	} else {
		f.Anchor = args[2]
	}

	return nil
}

// Apply applies the image filter to an image and returns the new image.
func (f *Crop) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	widthRepl := repl.ReplaceAll(f.Width, "")
	width, err := strconv.Atoi(widthRepl)
	if err != nil {
		return img, fmt.Errorf("invalid width %s %w", widthRepl, err)
	}
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %d", width)
	}

	heightRepl := repl.ReplaceAll(f.Height, "")
	height, err := strconv.Atoi(heightRepl)
	if err != nil {
		return img, fmt.Errorf("invalid height %s %w", heightRepl, err)
	}
	if height <= 0 {
		return img, fmt.Errorf("invalid height %d", height)
	}

	var anchor imaging.Anchor
	anchorRepl := repl.ReplaceAll(f.Anchor, "")
	switch anchorRepl {
	case "center":
		anchor = imaging.Center
	case "topleft":
		anchor = imaging.TopLeft
	case "top":
		anchor = imaging.Top
	case "topright":
		anchor = imaging.TopRight
	case "left":
		anchor = imaging.Left
	case "right":
		anchor = imaging.Right
	case "bottomleft":
		anchor = imaging.BottomLeft
	case "bottom":
		anchor = imaging.Bottom
	case "bottomright":
		anchor = imaging.BottomRight
	default:
		return nil, fmt.Errorf("invalid anchor '%s'", anchorRepl)
	}

	return imaging.CropAnchor(img, width, height, anchor), nil
}

// CaddyModule returns the Caddy module information.
func (Crop) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.crop",
		New: func() caddy.Module { return new(Crop) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Crop{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Crop)(nil)
	_ caddyfile.Unmarshaler = (*Crop)(nil)
)
