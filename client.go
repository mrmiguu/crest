package rest

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"
	"strings"
)

func newClient(addr string) *client {
	if i := strings.LastIndex(addr, "/"); i+1 == len(addr) {
		addr = addr[:i]
	}
	return &client{map[string]*Handler{}, addr}
}

func (s *server) New(pattern string) *Handler {
	h := &Handler{pattern: pattern}
	s.h[pattern] = h
	return h
}

func (c *client) write(t byte, b []byte) {
	b = append([]byte{t}, b...)
	http.Post(c.addr+"/post", "text/plain", bytes.NewReader(b))
}

func (c *client) read(t byte) []byte {
	b := []byte{t}
	resp, err := http.Post(c.addr+"/get", "text/plain", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return b
}

// Bytes creates a byte slice REST channel.
func (c *client) Bytes(pattern string) (func([]byte), func() []byte) {
	w := func(b []byte) {
		c.write(tbytes, b)
	}
	r := func() []byte {
		return c.read(tbytes)
	}
	return w, r
}

// String creates a string REST channel.
func (c *client) String(pattern string) (func(string), func() string) {
	w := func(s string) {
		c.write(tstring, []byte(s))
	}
	r := func() string {
		return string(c.read(tstring))
	}
	return w, r
}

// Int creates an int REST channel.
func (c *client) Int(pattern string) (func(int), func() int) {
	w := func(i int) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(i))
		c.write(tint, b)
	}
	r := func() int {
		return int(binary.BigEndian.Uint64(c.read(tint)))
	}
	return w, r
}
