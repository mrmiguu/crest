package rest

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"
	"strings"
)

func newClient(addr string) endpoint {
	if i := strings.LastIndex(addr, "/"); i+1 == len(addr) {
		addr = addr[:i]
	}
	return &client{h: map[string]*Handler{}, addr: addr}
}

func (c *client) New(pattern string) *Handler {
	h := &Handler{pattern: pattern}
	if _, exists := c.h[pattern]; exists {
		panic("pattern already exists")
	}
	c.h[pattern] = h
	return h
}

func (c *client) write(pattern string, t byte, idx []byte, msg []byte) {
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, idx, msg}, v)
	_, err := http.Post(c.addr+"/post", "text/plain", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
}

func (c *client) read(pattern string, t byte, idx []byte) []byte {
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, idx}, v)
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

func (c *client) Bytes(pattern string) (func([]byte), func() []byte) {
	idx := make([]byte, 8)
	binary.BigEndian.PutUint64(idx, uint64(len(c.h[pattern].getBytes)))
	c.h[pattern].getBytes = append(c.h[pattern].getBytes, nil)
	w := func(b []byte) { c.write(pattern, tbytes, idx, b) }
	r := func() []byte { return c.read(pattern, tbytes, idx) }
	return w, r
}

func (c *client) String(pattern string) (func(string), func() string) {
	idx := make([]byte, 8)
	binary.BigEndian.PutUint64(idx, uint64(len(c.h[pattern].getString)))
	c.h[pattern].getString = append(c.h[pattern].getString, nil)
	w := func(s string) { c.write(pattern, tstring, idx, []byte(s)) }
	r := func() string { return string(c.read(pattern, tstring, idx)) }
	return w, r
}

func (c *client) Int(pattern string) (func(int), func() int) {
	idx := make([]byte, 8)
	binary.BigEndian.PutUint64(idx, uint64(len(c.h[pattern].getInt)))
	c.h[pattern].getInt = append(c.h[pattern].getInt, nil)
	w := func(i int) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(i))
		c.write(pattern, tint, idx, b)
	}
	r := func() int { return int(binary.BigEndian.Uint64(c.read(pattern, tint, idx))) }
	return w, r
}
