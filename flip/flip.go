package flip

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// FlipFactory creates Flip instances.
type FlipFactory struct{}

// Flips flips (mirrors) a image vertically or horizontally.
type Flip struct {
	Direction string `json:"direction,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff FlipFactory) Name() string { return "flip" }

// New initialises and returns a configured Flip instance.
//
// Syntax:
//
//    flip <h|v>
//
// Parameters:
//
// h|v determines if the image flipped horizontally or vertically.
func (ff FlipFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 1 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 1 {
		return nil, imagefilter.ErrTooManyArgs
	}

	return Flip{Direction: args[0]}, nil
}

// Unmarshal decodes JSON data and returns a Flip instance.
func (ff FlipFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Flip{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Flip) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	direction := repl.ReplaceAll(f.Direction, "")

	switch direction {
	case "h":
		return imaging.FlipH(img), nil
	case "v":
		return imaging.FlipV(img), nil
	default:
		return nil, fmt.Errorf("unknown flip direction %s", direction)
	}
}

// init registers the image filter.
func init() {
	imagefilter.Register(FlipFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*FlipFactory)(nil)
	_ imagefilter.Filter        = (*Flip)(nil)
)
