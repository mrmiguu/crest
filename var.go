package rest

import (
	"regexp"
)

var (
	isServer = regexp.MustCompile(`:[0-9]+\b`)
	endpt    endpoint
	v        = []byte(`▼`)
)

func init() {
}
