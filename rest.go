package rest

import (
	"github.com/gopherjs/gopherjs/js"
)

// Connect connects to an endpoint for channel creation/communication.
func Connect(addr string) {
	if online {
		panic("already connected")
	}
	checkInstance()
	endpt.Connect(addr)
	online = true
}

// TODO: add thread safety
// TODO: add thread safety
func checkInstance() {
	if endpt != nil {
		return
	}
	if js.Global == nil {
		endpt = newServer()
	} else {
		endpt = newClient()
	}
	global = New("")
}

// Bytes creates a byte slice REST channel.
func Bytes(buf ...int) (func([]byte), func() []byte) {
	checkInstance()
	return global.Bytes(buf...)
}

// String creates a string REST channel.
func String(buf ...int) (func(string), func() string) {
	checkInstance()
	return global.String(buf...)
}

// Int creates an int REST channel.
func Int(buf ...int) (func(int), func() int) {
	checkInstance()
	return global.Int(buf...)
}

// Bool creates a bool REST channel.
func Bool(buf ...int) (func(bool), func() bool) {
	checkInstance()
	return global.Bool(buf...)
}

// New creates a handler for REST channel building.
func New(pattern string) *Handler {
	checkInstance()
	return endpt.New(pattern)
}

// Close closes the handler and releases all of its REST channels.
func (h *Handler) Close() error {
	h.hptr.Lock()
	defer h.hptr.Unlock()
	delete(h.hptr.m, h.Pattern)
	return nil
}

// Bytes creates a byte slice REST channel.
func (h *Handler) Bytes(buf ...int) (func([]byte), func() []byte) {
	if len(buf) > 1 {
		panic("too many arguments")
	}
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}
	return endpt.Bytes(h.Pattern, n)
}

// String creates a string REST channel.
func (h *Handler) String(buf ...int) (func(string), func() string) {
	if len(buf) > 1 {
		panic("too many arguments")
	}
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}
	return endpt.String(h.Pattern, n)
}

// Int creates an int REST channel.
func (h *Handler) Int(buf ...int) (func(int), func() int) {
	if len(buf) > 1 {
		panic("too many arguments")
	}
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}
	return endpt.Int(h.Pattern, n)
}

// Bool creates a bool REST channel.
func (h *Handler) Bool(buf ...int) (func(bool), func() bool) {
	if len(buf) > 1 {
		panic("too many arguments")
	}
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}
	return endpt.Bool(h.Pattern, n)
}
