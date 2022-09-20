package smartcrop

import (
	"encoding/json"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// SmartcropFactory creates Crop instances.
type SmartcropFactory struct{}

// Smartcrop finds good rectangular image crops of a specific size.
// It uses https://github.com/muesli/smartcrop
type Smartcrop struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff SmartcropFactory) Name() string { return "smartcrop" }

// New initialises and returns a configured Smartcrop instance.
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
func (ff SmartcropFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 3 {
		return nil, imagefilter.ErrTooManyArgs
	}

	return Smartcrop{Width: args[0], Height: args[1]}, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Smartcrop) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

// Unmarshal decodes JSON data and returns a Smartcrop instance.
func (ff SmartcropFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Smartcrop{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(SmartcropFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*SmartcropFactory)(nil)
	_ imagefilter.Filter        = (*Smartcrop)(nil)
)
