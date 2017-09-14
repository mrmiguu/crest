package crest

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Connect prepares communication with the specified URL.
func Connect(url string) {
	if strings.LastIndex(url, "/")+1 == len(url) {
		url = url[:len(url)-1]
	}
	address = url + Endpoint

	isClient = isClientExpr.MatchString(url)
	if isClient {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc(Endpoint+"/get", get)
	mux.HandleFunc(Endpoint+"/post", post)
	s := &http.Server{
		Addr:        url,
		Handler:     mux,
		ReadTimeout: 1 * time.Minute,
	}
	go func() { log.Fatal(s.ListenAndServe()) }()
}

// New creates a new pattern-specific handler for creating REST channels.
func New(pattern string) *Handler {
	h := &Handler{pattern: pattern}
	handlers.Lock()
	handlers.m[pattern] = h
	handlers.Unlock()
	return h
}

func (h *Handler) generic(t string, buf ...int) (func(interface{}), func() interface{}) {
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}

	h.wbytes.Lock()
	h.rbytes.Lock()
	defer h.wbytes.Unlock()
	defer h.rbytes.Unlock()

	w := wbytesf(h, t)
	r := rbytesf(h, t)

	h.wbytes.sl = append(h.wbytes.sl, callbacks{})
	h.rbytes.sl = append(h.rbytes.sl, make(chan interface{}, n))

	return w, r
}

// Int creates a writer and reader int REST channel.
func (h *Handler) Int(buf ...int) (func(int), func() int) {
	w, r := h.generic(tint, buf...)
	return func(i int) { w(i) }, func() int { return r().(int) }
}

// String creates a writer and reader int REST channel.
func (h *Handler) String(buf ...int) (func(string), func() string) {
	w, r := h.generic(tstring, buf...)
	return func(s string) { w(s) }, func() string { return r().(string) }
}

// Bytes creates a writer and reader byte slice REST channel.
func (h *Handler) Bytes(buf ...int) (func([]byte), func() []byte) {
	w, r := h.generic(tbytes, buf...)
	return func(b []byte) { w(b) }, func() []byte { return r().([]byte) }
}
