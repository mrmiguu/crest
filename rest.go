package rest

import "github.com/gopherjs/gopherjs/js"

// Connect connects to an endpoint for channel creation/communication.
func Connect(addr string) {
	if js.Global == nil {
		endpt = newServer(addr)
	} else {
		endpt = newClient(addr)
	}
}

// Bytes creates a byte slice REST channel for one-time use.
func Bytes(pattern string, buf ...int) (func([]byte), func() []byte) {
	h := New(pattern)
	wcloser, rcloser := h.Bytes(buf...)
	w := func(b []byte) { defer h.Close(); wcloser(b) }
	r := func() []byte { defer h.Close(); return rcloser() }
	return w, r
}

// String creates a string REST channel for one-time use.
func String(pattern string, buf ...int) (func(string), func() string) {
	h := New(pattern)
	wcloser, rcloser := h.String(buf...)
	w := func(s string) { defer h.Close(); wcloser(s) }
	r := func() string { defer h.Close(); return rcloser() }
	return w, r
}

// Int creates a int REST channel for one-time use.
func Int(pattern string, buf ...int) (func(int), func() int) {
	h := New(pattern)
	wcloser, rcloser := h.Int(buf...)
	w := func(i int) { defer h.Close(); wcloser(i) }
	r := func() int { defer h.Close(); return rcloser() }
	return w, r
}

// Bool creates a bool REST channel for one-time use.
func Bool(pattern string, buf ...int) (func(bool), func() bool) {
	h := New(pattern)
	wcloser, rcloser := h.Bool(buf...)
	w := func(b bool) { defer h.Close(); wcloser(b) }
	r := func() bool { defer h.Close(); return rcloser() }
	return w, r
}

// New creates a handler for REST channel building.
func New(pattern string) *Handler {
	if endpt == nil {
		panic("must connect first")
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
