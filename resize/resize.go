package resize

import (
	"encoding/json"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// ResizeFactory creates Resize instances.
type ResizeFactory struct{}

// Resize can downsize images. If upsizing of an image is detected, nothing will be done and
// the input image is returned unchanged.
type Resize struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff ResizeFactory) Name() string { return "resize" }

// New initialises and returns a ResizeFilter instance.
//
// Syntax:
//
//    resize <width> <height>
//
// Parameters:
//
// width must be a positive integer and determines the maximum width.
//
// height must be a positive integer and determines the maximum height.
//
// Either width or height can be 0, then the image aspect ratio is preserved.
func (ff ResizeFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 2 {
		return nil, imagefilter.ErrTooManyArgs
	}
	return Resize{Width: args[0], Height: args[1]}, nil
}

// Unmarshal decodes JSON data and returns a Resize instance.
func (ff ResizeFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Resize{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Resize) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

	if height < 0 || width < 0 || height == 0 && width == 0 {
		return img, fmt.Errorf("invalid width height combination %d %d", width, height)
	}

	// no upsizing
	if height == 0 && img.Bounds().Dx() <= width ||
		width == 0 && img.Bounds().Dy() <= height ||
		img.Bounds().Dx() <= width && img.Bounds().Dy() <= height {
		return img, nil
	}

	return imaging.Resize(img, width, height, imaging.Linear), nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(ResizeFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*ResizeFactory)(nil)
	_ imagefilter.Filter        = (*Resize)(nil)
)
