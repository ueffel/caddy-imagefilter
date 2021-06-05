package rotate

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

// RotateFactory creates Rotate instances.
type RotateFactory struct{}

// Rotate rotates a image 90, 180 or 270 degrees counter-clockwise.
type Rotate struct {
	Angle string `json:"angle,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff RotateFactory) Name() string { return "rotate" }

// New initialises and returns a configured Rotate instance.
//
// Syntax:
//
//    rotate <angle>
//
// Parameters:
//
// angle is one of the following: 0, 90, 180, 270 (0 is valid, but nothing will be done to the
// image).
func (ff RotateFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 1 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 1 {
		return nil, imagefilter.ErrTooManyArgs
	}
	return Rotate{Angle: args[0]}, nil
}

// Unmarshal decodes JSON data and returns a Rotate instance.
func (ff RotateFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Rotate{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Rotate) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	angleRepl := repl.ReplaceAll(f.Angle, "")
	angle, err := strconv.Atoi(angleRepl)
	if err != nil {
		return img, fmt.Errorf("invalid angle: %w", err)
	}

	switch angle {
	case 0:
		return img, nil
	case 90:
		return imaging.Rotate90(img), nil
	case 180:
		return imaging.Rotate180(img), nil
	case 270:
		return imaging.Rotate270(img), nil
	default:
		return nil, errors.New("invalid angle (only 0, 90, 180, 270 allowed)")
	}
}

// init registers the image filter.
func init() {
	imagefilter.Register(RotateFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*RotateFactory)(nil)
	_ imagefilter.Filter        = (*Rotate)(nil)
)
