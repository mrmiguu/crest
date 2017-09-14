package crest

import (
	"regexp"
	"sync"
)

var (
	address string

	isClientExpr = regexp.MustCompile(`^[Hh][Tt][Tt][Pp][Ss]{0,1}:`)
	isClient     bool

	handlers = struct {
		sync.RWMutex
		m map[string]*Handler
	}{m: map[string]*Handler{}}
)
