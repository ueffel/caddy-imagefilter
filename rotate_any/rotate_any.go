package rotate

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/disintegration/imaging"
	imagefilter "github.com/ueffel/caddy-imagefilter"
	"gopkg.in/go-playground/colors.v1"
)

// RotateAnyFactory creates RotateAny instances.
type RotateAnyFactory struct{}

// RotateAny rotates an image by a specific angle counter-clockwise. Uncovered areas after the
// rotation are filled with the specified color.
type RotateAny struct {
	Angle string `json:"angle,omitempty"`
	Color string `json:"color,omitempty"`
}

// Name returns the name of the filter, which is also the directive used in the image filter block.
func (ff RotateAnyFactory) Name() string { return "rotate_any" }

// New initialises and returns a configured Fit instance.
//
// Syntax:
//
//    rotate_any <angle> <color>
//
// Parameters:
//
// angle is the angle as floating point number in degrees by which the image is rotated
// counter-clockwise.
//
// color is the color which is used to fill uncovered areas after the rotation.
// Supported formats are:
//    "#FFAADD" (in quotes because otherwise it will be a comment in a caddyfile)
//    rgb(255,170,221)
//    rgba(255,170,221,0.5)
//    transparent, black, white, blue or about 140 more
//
// (see for many more supported color words https://www.w3schools.com/colors/colors_names.asp)
func (ff RotateAnyFactory) New(args ...string) (imagefilter.Filter, error) {
	if len(args) < 2 {
		return nil, imagefilter.ErrTooFewArgs
	}
	if len(args) > 2 {
		return nil, imagefilter.ErrTooManyArgs
	}
	return RotateAny{Angle: args[0], Color: args[1]}, nil
}

// Unmarshal decodes JSON data and returns a RotateAny instance.
func (ff RotateAnyFactory) Unmarshal(data []byte) (imagefilter.Filter, error) {
	filter := RotateAny{}
	err := json.Unmarshal(data, &filter)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

// Apply applies the image filter to an image and returns the new image.
func (f RotateAny) Apply(repl *caddy.Replacer, img image.Image) (image.Image, error) {
	angleRepl := repl.ReplaceAll(f.Angle, "")
	angle, err := strconv.ParseFloat(angleRepl, 64)
	if err != nil {
		return img, fmt.Errorf("invalid angle: %w", err)
	}
	colorRepl := repl.ReplaceAll(f.Color, "")
	bgColor := getColorFromName(colorRepl)
	if bgColor == nil {
		extractedColor, err := colors.Parse(colorRepl)
		if err != nil {
			return img, fmt.Errorf("invalid color: %w", err)
		}

		converted := extractedColor.ToRGBA()
		bgColor = color.NRGBA{R: converted.R, G: converted.G, B: converted.B, A: uint8(converted.A * 0xff)}
	}
	return imaging.Rotate(img, angle, bgColor), nil
}

// getColorFromName returns the RGB-Color for a color name. See
// https://www.w3schools.com/colors/colors_names.asp for supported names.
func getColorFromName(colorName string) color.Color {
	switch strings.ToLower(colorName) {
	case "aliceblue":
		return color.RGBA{R: 0xF0, G: 0xF8, B: 0xFF, A: 0xFF}
	case "antiquewhite":
		return color.RGBA{R: 0xFA, G: 0xEB, B: 0xD7, A: 0xFF}
	case "aqua":
		return color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
	case "aquamarine":
		return color.RGBA{R: 0x7F, G: 0xFF, B: 0xD4, A: 0xFF}
	case "azure":
		return color.RGBA{R: 0xF0, G: 0xFF, B: 0xFF, A: 0xFF}
	case "beige":
		return color.RGBA{R: 0xF5, G: 0xF5, B: 0xDC, A: 0xFF}
	case "bisque":
		return color.RGBA{R: 0xFF, G: 0xE4, B: 0xC4, A: 0xFF}
	case "black":
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	case "blanchedalmond":
		return color.RGBA{R: 0xFF, G: 0xEB, B: 0xCD, A: 0xFF}
	case "blue":
		return color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	case "blueviolet":
		return color.RGBA{R: 0x8A, G: 0x2B, B: 0xE2, A: 0xFF}
	case "brown":
		return color.RGBA{R: 0xA5, G: 0x2A, B: 0x2A, A: 0xFF}
	case "burlywood":
		return color.RGBA{R: 0xDE, G: 0xB8, B: 0x87, A: 0xFF}
	case "cadetblue":
		return color.RGBA{R: 0x5F, G: 0x9E, B: 0xA0, A: 0xFF}
	case "chartreuse":
		return color.RGBA{R: 0x7F, G: 0xFF, B: 0x00, A: 0xFF}
	case "chocolate":
		return color.RGBA{R: 0xD2, G: 0x69, B: 0x1E, A: 0xFF}
	case "coral":
		return color.RGBA{R: 0xFF, G: 0x7F, B: 0x50, A: 0xFF}
	case "cornflowerblue":
		return color.RGBA{R: 0x64, G: 0x95, B: 0xED, A: 0xFF}
	case "cornsilk":
		return color.RGBA{R: 0xFF, G: 0xF8, B: 0xDC, A: 0xFF}
	case "crimson":
		return color.RGBA{R: 0xDC, G: 0x14, B: 0x3C, A: 0xFF}
	case "cyan":
		return color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
	case "darkblue":
		return color.RGBA{R: 0x00, G: 0x00, B: 0x8B, A: 0xFF}
	case "darkcyan":
		return color.RGBA{R: 0x00, G: 0x8B, B: 0x8B, A: 0xFF}
	case "darkgoldenrod":
		return color.RGBA{R: 0xB8, G: 0x86, B: 0x0B, A: 0xFF}
	case "darkgray":
		return color.RGBA{R: 0xA9, G: 0xA9, B: 0xA9, A: 0xFF}
	case "darkgrey":
		return color.RGBA{R: 0xA9, G: 0xA9, B: 0xA9, A: 0xFF}
	case "darkgreen":
		return color.RGBA{R: 0x00, G: 0x64, B: 0x00, A: 0xFF}
	case "darkkhaki":
		return color.RGBA{R: 0xBD, G: 0xB7, B: 0x6B, A: 0xFF}
	case "darkmagenta":
		return color.RGBA{R: 0x8B, G: 0x00, B: 0x8B, A: 0xFF}
	case "darkolivegreen":
		return color.RGBA{R: 0x55, G: 0x6B, B: 0x2F, A: 0xFF}
	case "darkorange":
		return color.RGBA{R: 0xFF, G: 0x8C, B: 0x00, A: 0xFF}
	case "darkorchid":
		return color.RGBA{R: 0x99, G: 0x32, B: 0xCC, A: 0xFF}
	case "darkred":
		return color.RGBA{R: 0x8B, G: 0x00, B: 0x00, A: 0xFF}
	case "darksalmon":
		return color.RGBA{R: 0xE9, G: 0x96, B: 0x7A, A: 0xFF}
	case "darkseagreen":
		return color.RGBA{R: 0x8F, G: 0xBC, B: 0x8F, A: 0xFF}
	case "darkslateblue":
		return color.RGBA{R: 0x48, G: 0x3D, B: 0x8B, A: 0xFF}
	case "darkslategray":
		return color.RGBA{R: 0x2F, G: 0x4F, B: 0x4F, A: 0xFF}
	case "darkslategrey":
		return color.RGBA{R: 0x2F, G: 0x4F, B: 0x4F, A: 0xFF}
	case "darkturquoise":
		return color.RGBA{R: 0x00, G: 0xCE, B: 0xD1, A: 0xFF}
	case "darkviolet":
		return color.RGBA{R: 0x94, G: 0x00, B: 0xD3, A: 0xFF}
	case "deeppink":
		return color.RGBA{R: 0xFF, G: 0x14, B: 0x93, A: 0xFF}
	case "deepskyblue":
		return color.RGBA{R: 0x00, G: 0xBF, B: 0xFF, A: 0xFF}
	case "dimgray":
		return color.RGBA{R: 0x69, G: 0x69, B: 0x69, A: 0xFF}
	case "dimgrey":
		return color.RGBA{R: 0x69, G: 0x69, B: 0x69, A: 0xFF}
	case "dodgerblue":
		return color.RGBA{R: 0x1E, G: 0x90, B: 0xFF, A: 0xFF}
	case "firebrick":
		return color.RGBA{R: 0xB2, G: 0x22, B: 0x22, A: 0xFF}
	case "floralwhite":
		return color.RGBA{R: 0xFF, G: 0xFA, B: 0xF0, A: 0xFF}
	case "forestgreen":
		return color.RGBA{R: 0x22, G: 0x8B, B: 0x22, A: 0xFF}
	case "fuchsia":
		return color.RGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}
	case "gainsboro":
		return color.RGBA{R: 0xDC, G: 0xDC, B: 0xDC, A: 0xFF}
	case "ghostwhite":
		return color.RGBA{R: 0xF8, G: 0xF8, B: 0xFF, A: 0xFF}
	case "gold":
		return color.RGBA{R: 0xFF, G: 0xD7, B: 0x00, A: 0xFF}
	case "goldenrod":
		return color.RGBA{R: 0xDA, G: 0xA5, B: 0x20, A: 0xFF}
	case "gray":
		return color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	case "grey":
		return color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	case "green":
		return color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF}
	case "greenyellow":
		return color.RGBA{R: 0xAD, G: 0xFF, B: 0x2F, A: 0xFF}
	case "honeydew":
		return color.RGBA{R: 0xF0, G: 0xFF, B: 0xF0, A: 0xFF}
	case "hotpink":
		return color.RGBA{R: 0xFF, G: 0x69, B: 0xB4, A: 0xFF}
	case "indianred":
		return color.RGBA{R: 0xCD, G: 0x5C, B: 0x5C, A: 0xFF}
	case "indigo":
		return color.RGBA{R: 0x4B, G: 0x00, B: 0x82, A: 0xFF}
	case "ivory":
		return color.RGBA{R: 0xFF, G: 0xFF, B: 0xF0, A: 0xFF}
	case "khaki":
		return color.RGBA{R: 0xF0, G: 0xE6, B: 0x8C, A: 0xFF}
	case "lavender":
		return color.RGBA{R: 0xE6, G: 0xE6, B: 0xFA, A: 0xFF}
	case "lavenderblush":
		return color.RGBA{R: 0xFF, G: 0xF0, B: 0xF5, A: 0xFF}
	case "lawngreen":
		return color.RGBA{R: 0x7C, G: 0xFC, B: 0x00, A: 0xFF}
	case "lemonchiffon":
		return color.RGBA{R: 0xFF, G: 0xFA, B: 0xCD, A: 0xFF}
	case "lightblue":
		return color.RGBA{R: 0xAD, G: 0xD8, B: 0xE6, A: 0xFF}
	case "lightcoral":
		return color.RGBA{R: 0xF0, G: 0x80, B: 0x80, A: 0xFF}
	case "lightcyan":
		return color.RGBA{R: 0xE0, G: 0xFF, B: 0xFF, A: 0xFF}
	case "lightgoldenrodyellow":
		return color.RGBA{R: 0xFA, G: 0xFA, B: 0xD2, A: 0xFF}
	case "lightgray":
		return color.RGBA{R: 0xD3, G: 0xD3, B: 0xD3, A: 0xFF}
	case "lightgrey":
		return color.RGBA{R: 0xD3, G: 0xD3, B: 0xD3, A: 0xFF}
	case "lightgreen":
		return color.RGBA{R: 0x90, G: 0xEE, B: 0x90, A: 0xFF}
	case "lightpink":
		return color.RGBA{R: 0xFF, G: 0xB6, B: 0xC1, A: 0xFF}
	case "lightsalmon":
		return color.RGBA{R: 0xFF, G: 0xA0, B: 0x7A, A: 0xFF}
	case "lightseagreen":
		return color.RGBA{R: 0x20, G: 0xB2, B: 0xAA, A: 0xFF}
	case "lightskyblue":
		return color.RGBA{R: 0x87, G: 0xCE, B: 0xFA, A: 0xFF}
	case "lightslategray":
		return color.RGBA{R: 0x77, G: 0x88, B: 0x99, A: 0xFF}
	case "lightslategrey":
		return color.RGBA{R: 0x77, G: 0x88, B: 0x99, A: 0xFF}
	case "lightsteelblue":
		return color.RGBA{R: 0xB0, G: 0xC4, B: 0xDE, A: 0xFF}
	case "lightyellow":
		return color.RGBA{R: 0xFF, G: 0xFF, B: 0xE0, A: 0xFF}
	case "lime":
		return color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	case "limegreen":
		return color.RGBA{R: 0x32, G: 0xCD, B: 0x32, A: 0xFF}
	case "linen":
		return color.RGBA{R: 0xFA, G: 0xF0, B: 0xE6, A: 0xFF}
	case "magenta":
		return color.RGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}
	case "maroon":
		return color.RGBA{R: 0x80, G: 0x00, B: 0x00, A: 0xFF}
	case "mediumaquamarine":
		return color.RGBA{R: 0x66, G: 0xCD, B: 0xAA, A: 0xFF}
	case "mediumblue":
		return color.RGBA{R: 0x00, G: 0x00, B: 0xCD, A: 0xFF}
	case "mediumorchid":
		return color.RGBA{R: 0xBA, G: 0x55, B: 0xD3, A: 0xFF}
	case "mediumpurple":
		return color.RGBA{R: 0x93, G: 0x70, B: 0xDB, A: 0xFF}
	case "mediumseagreen":
		return color.RGBA{R: 0x3C, G: 0xB3, B: 0x71, A: 0xFF}
	case "mediumslateblue":
		return color.RGBA{R: 0x7B, G: 0x68, B: 0xEE, A: 0xFF}
	case "mediumspringgreen":
		return color.RGBA{R: 0x00, G: 0xFA, B: 0x9A, A: 0xFF}
	case "mediumturquoise":
		return color.RGBA{R: 0x48, G: 0xD1, B: 0xCC, A: 0xFF}
	case "mediumvioletred":
		return color.RGBA{R: 0xC7, G: 0x15, B: 0x85, A: 0xFF}
	case "midnightblue":
		return color.RGBA{R: 0x19, G: 0x19, B: 0x70, A: 0xFF}
	case "mintcream":
		return color.RGBA{R: 0xF5, G: 0xFF, B: 0xFA, A: 0xFF}
	case "mistyrose":
		return color.RGBA{R: 0xFF, G: 0xE4, B: 0xE1, A: 0xFF}
	case "moccasin":
		return color.RGBA{R: 0xFF, G: 0xE4, B: 0xB5, A: 0xFF}
	case "navajowhite":
		return color.RGBA{R: 0xFF, G: 0xDE, B: 0xAD, A: 0xFF}
	case "navy":
		return color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xFF}
	case "oldlace":
		return color.RGBA{R: 0xFD, G: 0xF5, B: 0xE6, A: 0xFF}
	case "olive":
		return color.RGBA{R: 0x80, G: 0x80, B: 0x00, A: 0xFF}
	case "olivedrab":
		return color.RGBA{R: 0x6B, G: 0x8E, B: 0x23, A: 0xFF}
	case "orange":
		return color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}
	case "orangered":
		return color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}
	case "orchid":
		return color.RGBA{R: 0xDA, G: 0x70, B: 0xD6, A: 0xFF}
	case "palegoldenrod":
		return color.RGBA{R: 0xEE, G: 0xE8, B: 0xAA, A: 0xFF}
	case "palegreen":
		return color.RGBA{R: 0x98, G: 0xFB, B: 0x98, A: 0xFF}
	case "paleturquoise":
		return color.RGBA{R: 0xAF, G: 0xEE, B: 0xEE, A: 0xFF}
	case "palevioletred":
		return color.RGBA{R: 0xDB, G: 0x70, B: 0x93, A: 0xFF}
	case "papayawhip":
		return color.RGBA{R: 0xFF, G: 0xEF, B: 0xD5, A: 0xFF}
	case "peachpuff":
		return color.RGBA{R: 0xFF, G: 0xDA, B: 0xB9, A: 0xFF}
	case "peru":
		return color.RGBA{R: 0xCD, G: 0x85, B: 0x3F, A: 0xFF}
	case "pink":
		return color.RGBA{R: 0xFF, G: 0xC0, B: 0xCB, A: 0xFF}
	case "plum":
		return color.RGBA{R: 0xDD, G: 0xA0, B: 0xDD, A: 0xFF}
	case "powderblue":
		return color.RGBA{R: 0xB0, G: 0xE0, B: 0xE6, A: 0xFF}
	case "purple":
		return color.RGBA{R: 0x80, G: 0x00, B: 0x80, A: 0xFF}
	case "rebeccapurple":
		return color.RGBA{R: 0x66, G: 0x33, B: 0x99, A: 0xFF}
	case "red":
		return color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	case "rosybrown":
		return color.RGBA{R: 0xBC, G: 0x8F, B: 0x8F, A: 0xFF}
	case "royalblue":
		return color.RGBA{R: 0x41, G: 0x69, B: 0xE1, A: 0xFF}
	case "saddlebrown":
		return color.RGBA{R: 0x8B, G: 0x45, B: 0x13, A: 0xFF}
	case "salmon":
		return color.RGBA{R: 0xFA, G: 0x80, B: 0x72, A: 0xFF}
	case "sandybrown":
		return color.RGBA{R: 0xF4, G: 0xA4, B: 0x60, A: 0xFF}
	case "seagreen":
		return color.RGBA{R: 0x2E, G: 0x8B, B: 0x57, A: 0xFF}
	case "seashell":
		return color.RGBA{R: 0xFF, G: 0xF5, B: 0xEE, A: 0xFF}
	case "sienna":
		return color.RGBA{R: 0xA0, G: 0x52, B: 0x2D, A: 0xFF}
	case "silver":
		return color.RGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	case "skyblue":
		return color.RGBA{R: 0x87, G: 0xCE, B: 0xEB, A: 0xFF}
	case "slateblue":
		return color.RGBA{R: 0x6A, G: 0x5A, B: 0xCD, A: 0xFF}
	case "slategray":
		return color.RGBA{R: 0x70, G: 0x80, B: 0x90, A: 0xFF}
	case "slategrey":
		return color.RGBA{R: 0x70, G: 0x80, B: 0x90, A: 0xFF}
	case "snow":
		return color.RGBA{R: 0xFF, G: 0xFA, B: 0xFA, A: 0xFF}
	case "springgreen":
		return color.RGBA{R: 0x00, G: 0xFF, B: 0x7F, A: 0xFF}
	case "steelblue":
		return color.RGBA{R: 0x46, G: 0x82, B: 0xB4, A: 0xFF}
	case "tan":
		return color.RGBA{R: 0xD2, G: 0xB4, B: 0x8C, A: 0xFF}
	case "teal":
		return color.RGBA{R: 0x00, G: 0x80, B: 0x80, A: 0xFF}
	case "thistle":
		return color.RGBA{R: 0xD8, G: 0xBF, B: 0xD8, A: 0xFF}
	case "tomato":
		return color.RGBA{R: 0xFF, G: 0x63, B: 0x47, A: 0xFF}
	case "transparent":
		return color.Transparent
	case "turquoise":
		return color.RGBA{R: 0x40, G: 0xE0, B: 0xD0, A: 0xFF}
	case "violet":
		return color.RGBA{R: 0xEE, G: 0x82, B: 0xEE, A: 0xFF}
	case "wheat":
		return color.RGBA{R: 0xF5, G: 0xDE, B: 0xB3, A: 0xFF}
	case "white":
		return color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	case "whitesmoke":
		return color.RGBA{R: 0xF5, G: 0xF5, B: 0xF5, A: 0xFF}
	case "yellow":
		return color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}
	case "yellowgreen":
		return color.RGBA{R: 0x9A, G: 0xCD, B: 0x32, A: 0xFF}
	default:
		return nil
	}
}

// init registers the image filter.
func init() {
	imagefilter.Register(RotateAnyFactory{})
}

// Interface guards.
var (
	_ imagefilter.FilterFactory = (*RotateAnyFactory)(nil)
	_ imagefilter.Filter        = (*RotateAny)(nil)
)
