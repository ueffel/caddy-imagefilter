package grayscale

import (
	"encoding/json"
	"errors"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

type GrayscaleFactory struct{}

type Grayscale struct{}

func (ff GrayscaleFactory) Name() string { return "grayscale" }

func (ff GrayscaleFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 0 {
		return nil, errors.New("too many arguments")
	}

	return Grayscale{}, nil
}

func (ff GrayscaleFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Grayscale{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f Grayscale) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Grayscale(img), nil
}

func init() {
	imagefilter.Register(GrayscaleFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*GrayscaleFactory)(nil)
	_ imagefilter.Filter        = (*Grayscale)(nil)
)
