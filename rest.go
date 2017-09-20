package rest

import "github.com/gopherjs/gopherjs/js"
import "runtime"

// Connect connects to an endpoint for channel creation/communication.
func Connect(addr string) {
	if js.Global == nil {
		endpt = newServer(addr)
	} else {
		endpt = newClient(addr)
	}
}

// Bytes creates a volatile byte slice REST channel.
func Bytes(pattern string, buf ...int) (func([]byte), func() []byte) {
	h := New(pattern)
	w, r := h.Bytes(buf...)
	runtime.SetFinalizer(w, func(_ func(_ []byte)) { h.Close() })
	runtime.SetFinalizer(r, func(_ func() []byte) { h.Close() })
	return w, r
}

// String creates a volatile string REST channel.
func String(pattern string, buf ...int) (func(string), func() string) {
	h := New(pattern)
	w, r := h.String(buf...)
	runtime.SetFinalizer(w, func(_ func(_ string)) { h.Close() })
	runtime.SetFinalizer(r, func(_ func() string) { h.Close() })
	return w, r
}

// Int creates a volatile int REST channel.
func Int(pattern string, buf ...int) (func(int), func() int) {
	h := New(pattern)
	w, r := h.Int(buf...)
	runtime.SetFinalizer(w, func(_ func(_ int)) { h.Close() })
	runtime.SetFinalizer(r, func(_ func() int) { h.Close() })
	return w, r
}

// Bool creates a volatile bool REST channel.
func Bool(pattern string, buf ...int) (func(bool), func() bool) {
	h := New(pattern)
	w, r := h.Bool(buf...)
	runtime.SetFinalizer(w, func(_ func(_ bool)) { h.Close() })
	runtime.SetFinalizer(r, func(_ func() bool) { h.Close() })
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
