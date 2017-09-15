package rest

// New creates a handler for REST channel building.
func New(pattern string) *Handler {
	return &Handler{pattern: pattern}
}

// Bytes creates a byte slice REST channel.
func (h *Handler) Bytes() (func([]byte), func() []byte) {
	return e.Bytes(h.pattern)
}

// String creates a string REST channel.
func (h *Handler) String() (func(string), func() string) {
	return e.String(h.pattern)
}

// Int creates an int REST channel.
func (h *Handler) Int() (func(int), func() int) {
	return e.Int(h.pattern)
}
