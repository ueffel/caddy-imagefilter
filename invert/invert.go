package invert

import (
	"encoding/json"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// InvertFactory creates Invert instances.
type InvertFactory struct{}

// Invert produces an inverted (negated) version of the image.
type Invert struct{}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff InvertFactory) Name() string { return "invert" }

// New initialises and returns a configured Grayscale instance.
//
// Syntax:
//
//    invert
//
// no parameters.
func (ff InvertFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 0 {
		return nil, imagefilter.ErrTooManyArgs
	}

	return Invert{}, nil
}

// Unmarshal decodes JSON data and returns a Invert instance.
func (ff InvertFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Invert{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Invert) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Invert(img), nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(InvertFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*InvertFactory)(nil)
	_ imagefilter.Filter        = (*Invert)(nil)
)
