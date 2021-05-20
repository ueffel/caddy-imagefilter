package fit

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

type FitFactory struct{}

type Fit struct {
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

func (ff FitFactory) Name() string { return "fit" }

func (ff FitFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}
	if len(args) > 2 {
		return nil, errors.New("too many arguments")
	}

	return Fit{Width: args[0], Height: args[1]}, nil
}

func (ff FitFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Fit{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f Fit) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
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

	if height <= 0 || width <= 0 {
		return img, fmt.Errorf("invalid width height combination %d %d", width, height)
	}

	return imaging.Fit(img, width, height, imaging.Linear), nil
}

func init() {
	imagefilter.Register(FitFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*FitFactory)(nil)
	_ imagefilter.Filter        = (*Fit)(nil)
)
