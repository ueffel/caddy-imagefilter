package crop

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

type CropFactory struct{}

type Crop struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Anchor string `json:"anchor,omitempty"`
}

func (ff CropFactory) Name() string { return "crop" }

func (ff CropFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}
	if len(args) > 3 {
		return nil, errors.New("too many arguments")
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

func (ff CropFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Crop{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f Crop) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var width int
	widthRepl := repl.ReplaceAll(f.Width, "")
	width, err = strconv.Atoi(widthRepl)
	if err != nil {
		return img, fmt.Errorf("invalid width %s %v", widthRepl, err)
	}
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %d", width)
	}

	var height int
	heightRepl := repl.ReplaceAll(f.Height, "")
	height, err = strconv.Atoi(heightRepl)
	if err != nil {
		return img, fmt.Errorf("invalid height %s %v", heightRepl, err)
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

func init() {
	imagefilter.Register(CropFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*CropFactory)(nil)
	_ imagefilter.Filter        = (*Crop)(nil)
)
