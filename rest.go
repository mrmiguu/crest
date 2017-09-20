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
func Bytes(pattern string, buf ...int) func(...[]byte) []byte {
	h := New(pattern)
	w, r := h.Bytes(buf...)
	return func(x ...[]byte) (y []byte) {
		defer h.Close()
		if len(x) > 1 {
			panic("too many arguments")
		} else if len(x) > 0 {
			w(x[0])
		} else {
			y = r()
		}
		return
	}
}

// String creates a string REST channel for one-time use.
func String(pattern string, buf ...int) func(...string) string {
	h := New(pattern)
	w, r := h.String(buf...)
	return func(x ...string) (y string) {
		defer h.Close()
		if len(x) > 1 {
			panic("too many arguments")
		} else if len(x) > 0 {
			w(x[0])
		} else {
			y = r()
		}
		return
	}
}

// Int creates a int REST channel for one-time use.
func Int(pattern string, buf ...int) func(...int) int {
	h := New(pattern)
	w, r := h.Int(buf...)
	return func(x ...int) (y int) {
		defer h.Close()
		if len(x) > 1 {
			panic("too many arguments")
		} else if len(x) > 0 {
			w(x[0])
		} else {
			y = r()
		}
		return
	}
}

// Bool creates a bool REST channel for one-time use.
func Bool(pattern string, buf ...int) func(...bool) bool {
	h := New(pattern)
	w, r := h.Bool(buf...)
	return func(x ...bool) (y bool) {
		defer h.Close()
		if len(x) > 1 {
			panic("too many arguments")
		} else if len(x) > 0 {
			w(x[0])
		} else {
			y = r()
		}
		return
	}
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
