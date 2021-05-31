package fit

import (
	"encoding/json"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// FitFactory creates Fit instances.
type FitFactory struct{}

// Fit scales a image to fit to the specified maximum width and height using a linear filter, the
// image aspect ratio is preserved. If the image already fits inside the bounds, nothing will be
// done.
type Fit struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff FitFactory) Name() string { return "fit" }

// New initialises and returns a configured Fit instance.
//
// Syntax:
//
//    fit <width> <height>
//
// Parameters:
//
// width must be a positive integer and determines the maximum width.
//
// height must be a positive integer and determines the maximum height.
func (ff FitFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 2 {
		return nil, imagefilter.ErrTooManyArgs
	}

	return Fit{Width: args[0], Height: args[1]}, nil
}

// Unmarshal decodes JSON data and returns a Fit instance.
func (ff FitFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Fit{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Fit) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

// init registers the image filter.
func init() {
	imagefilter.Register(FitFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*FitFactory)(nil)
	_ imagefilter.Filter        = (*Fit)(nil)
)
