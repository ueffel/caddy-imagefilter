package sharpen

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

// SharpenFactory creates Sharpen instances.
type SharpenFactory struct{}

// Sharpen produces a sharpened version of the image.
type Sharpen struct {
	Sigma string `json:"sigma,omitempty"`
}

// Name returns the name if the filter, which is also the directive used in the image filter block.
func (ff SharpenFactory) Name() string { return "sharpen" }

// New initialises and returns a Sharpen instance.
//
// Syntax:
//
//    sharpen [<sigma>]
//
// Parameters:
//
// sigma must be a positive floating point number and indicates how much the image will be
// sharpened. Default is 1.
func (ff SharpenFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 1 {
		return nil, imagefilter.ErrTooManyArgs
	}
	var sigma string
	if len(args) == 1 {
		sigma = args[0]
	}
	return Sharpen{Sigma: sigma}, nil
}

// Unmarshal decodes JSON data and returns a Sharpen instance.
func (ff SharpenFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Sharpen{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Sharpen) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var sigma float64
	sigmaRepl := repl.ReplaceAll(f.Sigma, "")
	if sigmaRepl == "" {
		sigma = 1
	} else {
		sigma, err = strconv.ParseFloat(sigmaRepl, 64)
		if err != nil {
			return img, fmt.Errorf("invalid sigma: %w", err)
		}
	}

	if sigma <= 0 {
		return img, errors.New("invalid sigma: cannot be less or equal 0")
	}

	return imaging.Sharpen(img, sigma), nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(SharpenFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*SharpenFactory)(nil)
	_ imagefilter.Filter        = (*Sharpen)(nil)
)
