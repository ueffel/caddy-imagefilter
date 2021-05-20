package flip

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

type FlipFactory struct{}

type Flip struct {
	Direction string `json:"direction,omitempty"`
}

func (ff FlipFactory) Name() string { return "flip" }

func (ff FlipFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}
	if len(args) > 1 {
		return nil, errors.New("too many arguments")
	}

	return Flip{Direction: args[0]}, nil
}

func (ff FlipFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Flip{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

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

func init() {
	imagefilter.Register(FlipFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*FlipFactory)(nil)
	_ imagefilter.Filter        = (*Flip)(nil)
)
