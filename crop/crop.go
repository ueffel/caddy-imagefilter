package crop

import (
	"encoding/json"
	"fmt"
	"image"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

// CropFactory creates Crop instances.
type CropFactory struct{}

// Crop produces a cropped image as rectangular region of a specific size.
type Crop struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Anchor string `json:"anchor,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff CropFactory) Name() string { return "crop" }

// New initialises and returns a configured Crop instance.
//
// Syntax:
//
//    crop <width> <height> [<anchor>]
//
// Parameters:
//
// width must be a positive integer and determines the width of the cropped image.
//
// height must be a positive integer and determines the height of the cropped image.
//
// anchor determines the anchor point of the rectangular region that is cut out. Possible values
// are: center, topleft, top, topright, left, right, bottomleft, bottom, bottomright.
// Default is center.
func (ff CropFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 3 {
		return nil, imagefilter.ErrTooManyArgs
	}

	var anchor string
	if len(args) < 3 {
		anchor = "center"
	} else {
		anchor = args[2]
	}

	return Crop{
		Width:  args[0],
		Height: args[1],
		Anchor: anchor,
	}, nil
}

// Unmarshal decodes JSON data and returns a Crop instance.
func (ff CropFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Crop{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f Crop) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var width int
	widthRepl := repl.ReplaceAll(f.Width, "")
	width, err = strconv.Atoi(widthRepl)
	if err != nil {
		return img, fmt.Errorf("invalid width %s %w", widthRepl, err)
	}
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %d", width)
	}

	var height int
	heightRepl := repl.ReplaceAll(f.Height, "")
	height, err = strconv.Atoi(heightRepl)
	if err != nil {
		return img, fmt.Errorf("invalid height %s %w", heightRepl, err)
	}
	if height <= 0 {
		return img, fmt.Errorf("invalid height %d", height)
	}

	var anchor imaging.Anchor
	anchorRepl := repl.ReplaceAll(f.Anchor, "")
	switch anchorRepl {
	case "center":
		anchor = imaging.Center
	case "topleft":
		anchor = imaging.TopLeft
	case "top":
		anchor = imaging.Top
	case "topright":
		anchor = imaging.TopRight
	case "left":
		anchor = imaging.Left
	case "right":
		anchor = imaging.Right
	case "bottomleft":
		anchor = imaging.BottomLeft
	case "bottom":
		anchor = imaging.Bottom
	case "bottomright":
		anchor = imaging.BottomRight
	default:
		return nil, fmt.Errorf("invalid anchor '%s'", anchorRepl)
	}

	return imaging.CropAnchor(img, width, height, anchor), nil
}

// init registers the image filter.
func init() {
	imagefilter.Register(CropFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*CropFactory)(nil)
	_ imagefilter.Filter        = (*Crop)(nil)
)
