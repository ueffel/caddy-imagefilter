package resize

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// ResizeFilterFactory creates ResizeFilter instances
type ResizeFilterFactory struct{}

// ResizeFilter can downsize images. If upsizing of an image is detected, nothing will be done and
// the input image is returned unchanged
type ResizeFilter struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// Name retrurns the name if the filter, which is also the directive used in the image filter block
func (ff ResizeFilterFactory) Name() string { return "resize" }

// New intitialises and returns a ResizeFilter instance.
func (ff ResizeFilterFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}
	if len(args) > 2 {
		return nil, errors.New("too many arguments")
	}
	return ResizeFilter{Width: args[0], Height: args[1]}, nil
}

// Unmarshal decodes JSON data and returns a ResizeFilter instance.
func (ff ResizeFilterFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := ResizeFilter{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply does the resizing of an image. If upsizing of an image is detected, nothing will be done
// and the input image is returned unchanged
func (f ResizeFilter) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var width int
	widthRepl := repl.ReplaceAll(f.Width, "")
	if widthRepl == "" {
		width = 0
	} else {
		width, err = strconv.Atoi(widthRepl)
		if err != nil {
			return img, fmt.Errorf("invalid width: %v", err)
		}
	}
	var height int
	heightRepl := repl.ReplaceAll(f.Height, "")
	if heightRepl == "" {
		height = 0
	} else {
		height, err = strconv.Atoi(heightRepl)
		if err != nil {
			return img, fmt.Errorf("invalid height: %v", err)
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

func init() {
	imagefilter.Register(ResizeFilterFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*ResizeFilterFactory)(nil)
	_ imagefilter.Filter        = (*ResizeFilter)(nil)
)
