package grayscale

import (
	"encoding/json"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// GrayscaleFactory creates Grayscale instances.
type GrayscaleFactory struct{}

// Grayscale produces a grayscaled version of the image.
type Grayscale struct{}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff GrayscaleFactory) Name() string { return "grayscale" }

// New initialises and returns a configured Grayscale instance.
//
// Syntax:
//
//	grayscale
//
// no parameters.
func (ff GrayscaleFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 0 {
		return nil, imagefilter.ErrTooManyArgs
	}

	return Grayscale{}, nil
}

// Unmarshal decodes JSON data and returns a Grayscale instance.
func (ff GrayscaleFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Grayscale{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Grayscale) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Grayscale(img), nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(GrayscaleFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*GrayscaleFactory)(nil)
	_ imagefilter.Filter        = (*Grayscale)(nil)
)
