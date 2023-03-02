// Package defaults provides commonly used image filters as a bundle in a single
// import.
package defaults

import (
	_ "github.com/ueffel/caddy-imagefilter/v2/crop"
	_ "github.com/ueffel/caddy-imagefilter/v2/fit"
	_ "github.com/ueffel/caddy-imagefilter/v2/flip"
	_ "github.com/ueffel/caddy-imagefilter/v2/resize"
	_ "github.com/ueffel/caddy-imagefilter/v2/rotate"
	_ "github.com/ueffel/caddy-imagefilter/v2/sharpen"
)
