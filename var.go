package rest

import (
	"regexp"
)

var (
	isServer = regexp.MustCompile(`:[0-9]+/`)
	endpt    endpoint
	v        = []byte(`▼`)
)

func init() {
}
