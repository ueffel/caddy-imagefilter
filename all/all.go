// Package all include all implemented image filters as a bundle in a single
// import.
package all

import (
	_ "github.com/ueffel/caddy-imagefilter/v2/blur"
	_ "github.com/ueffel/caddy-imagefilter/v2/crop"
	_ "github.com/ueffel/caddy-imagefilter/v2/fit"
	_ "github.com/ueffel/caddy-imagefilter/v2/flip"
	_ "github.com/ueffel/caddy-imagefilter/v2/grayscale"
	_ "github.com/ueffel/caddy-imagefilter/v2/invert"
	_ "github.com/ueffel/caddy-imagefilter/v2/resize"
	_ "github.com/ueffel/caddy-imagefilter/v2/rotate"
	_ "github.com/ueffel/caddy-imagefilter/v2/rotate_any"
	_ "github.com/ueffel/caddy-imagefilter/v2/sharpen"
	_ "github.com/ueffel/caddy-imagefilter/v2/smartcrop"
)
