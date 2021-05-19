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

type RotateFilterFactory struct{}

type RotateFilter struct {
	Angle string `json:"angle,omitempty"`
}

func (ff RotateFilterFactory) Name() string { return "rotate" }

func (ff RotateFilterFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}
	if len(args) > 1 {
		return nil, errors.New("too many arguments")
	}
	return RotateFilter{Angle: args[0]}, nil
}

func (ff RotateFilterFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := RotateFilter{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f RotateFilter) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	angleRepl := repl.ReplaceAll(f.Angle, "")
	angle, err := strconv.Atoi(angleRepl)
	if err != nil {
		return img, fmt.Errorf("invalid angle: %v", err)
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

func init() {
	imagefilter.Register(RotateFilterFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*RotateFilterFactory)(nil)
	_ imagefilter.Filter        = (*RotateFilter)(nil)
)
