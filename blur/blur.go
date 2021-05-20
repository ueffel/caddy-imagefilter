package blur

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

// BlurFactory creates BlurFilter instances.
type BlurFactory struct{}

// Blur produces a Blured version of the image.
type Blur struct {
	Sigmna string `json:"sigma,omitempty"`
}

// Name retrurns the name if the filter, which is also the directive used in the image filter block.
func (ff BlurFactory) Name() string { return "blur" }

// New intitialises and returns a BlurFilter instance.
func (ff BlurFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) > 1 {
		return nil, errors.New("too many arguments")
	}
	var sigma string
	if len(args) == 1 {
		sigma = args[0]
	}
	return Blur{Sigmna: sigma}, nil
}

// Unmarshal decodes JSON data and returns a BlurFilter instance.
func (ff BlurFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := Blur{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func (f Blur) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	var err error
	var sigma float64
	sigmaRepl := repl.ReplaceAll(f.Sigmna, "")
	if sigmaRepl == "" {
		sigma = 1
	} else {
		sigma, err = strconv.ParseFloat(sigmaRepl, 64)
		if err != nil {
			return img, fmt.Errorf("invalid sigma: %v", err)
		}
	}

	if sigma <= 0 {
		return img, errors.New("invalid sigma: cannot be less or equal 0")
	}

	return imaging.Blur(img, sigma), nil
}

func init() {
	imagefilter.Register(BlurFactory{})
}

// Interface Guards
var (
	_ imagefilter.FilterFactory = (*BlurFactory)(nil)
	_ imagefilter.Filter        = (*Blur)(nil)
)
