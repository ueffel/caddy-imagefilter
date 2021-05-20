package invert

import (
	"encoding/json"
	"errors"
	"image"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
)

type InvertFactory struct{}

type Invert struct{}

func (ff InvertFactory) Name() string { return "invert" }

func (ff InvertFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 0 {
		return nil, errors.New("too many arguments")
	}

	return Invert{}, nil
}

func (ff InvertFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Invert{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f Invert) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	return imaging.Invert(img), nil
}

func init() {
	imagefilter.Register(InvertFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*InvertFactory)(nil)
	_ imagefilter.Filter        = (*Invert)(nil)
)
