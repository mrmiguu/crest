package rest

import "github.com/gopherjs/gopherjs/js"

// Connect connects to an endpoint for channel creation/communication.
func Connect(addr string) {
	if endpt != nil {
		panic("already connected")
	}
	if js.Global == nil {
		endpt = newServer(addr)
	} else {
		endpt = newClient(addr)
	}
	global = New("")
}

// Bytes creates a byte slice REST channel.
func Bytes(buf ...int) (func([]byte), func() []byte) {
	return global.Bytes(buf...)
}

// String creates a string REST channel.
func String(buf ...int) (func(string), func() string) {
	return global.String(buf...)
}

// Int creates an int REST channel.
func Int(buf ...int) (func(int), func() int) {
	return global.Int(buf...)
}

// Bool creates a bool REST channel.
func Bool(buf ...int) (func(bool), func() bool) {
	return global.Bool(buf...)
}

// New creates a handler for REST channel building.
func New(pattern string) *Handler {
	if endpt == nil {
		Connect("/")
	}
	return endpt.New(pattern)
}

// Close closes the handler and releases all of its REST channels.
func (h *Handler) Close() error {
	h.hptr.Lock()
	defer h.hptr.Unlock()
	delete(h.hptr.m, h.pattern)
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
	return endpt.Bytes(h.pattern, n)
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
	return endpt.String(h.pattern, n)
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
	return endpt.Int(h.pattern, n)
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
	return endpt.Bool(h.pattern, n)
}
