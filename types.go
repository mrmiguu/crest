package crest

import "sync"

// Handler holds pattern-specific read/write channels.
type Handler struct {
	pattern string

	wbytes struct {
		sync.RWMutex
		sl []callbacks
	}
	rbytes struct {
		sync.RWMutex
		sl []chan interface{}
	}
}

type callbacks struct {
	sync.RWMutex
	sl []chan interface{}
}
