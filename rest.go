package rest

// Connect connects to an endpoint for channel creation/communication.
func Connect(addr string) {
	if isServer.MatchString(addr) {
		endpt = newServer(addr)
	} else {
		endpt = newClient(addr)
	}
}

// New creates a handler for REST channel building.
func New(pattern string) *Handler {
	return endpt.New(pattern)
}

// Close closes the handler and releases all of its REST channels.
func (h *Handler) Close() error {
	delete(*h.hptr, h.pattern)
	return nil
}

// Bytes creates a byte slice REST channel.
func (h *Handler) Bytes() (func([]byte), func() []byte) {
	return endpt.Bytes(h.pattern)
}

// String creates a string REST channel.
func (h *Handler) String() (func(string), func() string) {
	return endpt.String(h.pattern)
}

// Int creates an int REST channel.
func (h *Handler) Int() (func(int), func() int) {
	return endpt.Int(h.pattern)
}
