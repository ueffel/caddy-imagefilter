package smartcrop

import (
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/disintegration/imaging"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	imagefilter "github.com/ueffel/caddy-imagefilter/v2"
)

// Smartcrop finds good rectangular image crops of a specific size.
// It uses https://github.com/muesli/smartcrop
type Smartcrop struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// UnmarshalCaddyfile configures the Smartcrop instance.
//
// Syntax:
//
//	smartcrop <width> <height>
//
// Parameters:
//
// width must be a positive integer and determines the width of the cropped image.
//
// height must be a positive integer and determines the height of the cropped image.
func (f *Smartcrop) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
func (f *Smartcrop) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

	analyzer := smartcrop.NewAnalyzer(nfnt.NewResizer(resize.Bilinear))
	topCrop, err := analyzer.FindBestCrop(img, width, height)
	if err != nil {
		return img, fmt.Errorf("determining smartcrop %w", err)
	}

	cropped := imaging.Crop(img, topCrop)
	return imaging.Resize(cropped, width, height, imaging.Linear), nil
}

// CaddyModule returns the Caddy module information.
func (Smartcrop) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_filter.filter.smartcrop",
		New: func() caddy.Module { return new(Smartcrop) },
	}
}

// init registers the image filter.
func init() {
	caddy.RegisterModule(Smartcrop{})
}

// Interface guards.
var (
	_ imagefilter.Filter    = (*Smartcrop)(nil)
	_ caddyfile.Unmarshaler = (*Smartcrop)(nil)
)
