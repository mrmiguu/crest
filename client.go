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
	if _, exists := c.h[pattern]; exists {
		panic("pattern already exists")
	}
	h := &Handler{hptr: &c.h, pattern: pattern}
	c.h[pattern] = h
	return h
}

func (c *client) write(pattern string, t byte, idx int, msg []byte) {
	i := make([]byte, 8)
	binary.BigEndian.PutUint64(i, uint64(idx))
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, i, msg}, v)
	_, err := http.Post(c.addr+"/post", "text/plain", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
}

func (c *client) read(pattern string, t byte, idx int) []byte {
	i := make([]byte, 8)
	binary.BigEndian.PutUint64(i, uint64(idx))
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, i}, v)
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
	idx := len(c.h[pattern].getBytes)
	c.h[pattern].getBytes = append(c.h[pattern].getBytes, nil)
	c.h[pattern].postBytes = append(c.h[pattern].postBytes, nil)
	w := func(b []byte) { c.write(pattern, tbytes, idx, b) }
	r := func() []byte { return c.read(pattern, tbytes, idx) }
	return w, r
}

func (c *client) String(pattern string) (func(string), func() string) {
	idx := len(c.h[pattern].getString)
	c.h[pattern].getString = append(c.h[pattern].getString, nil)
	c.h[pattern].postString = append(c.h[pattern].postString, nil)
	w := func(s string) { c.write(pattern, tstring, idx, []byte(s)) }
	r := func() string { return string(c.read(pattern, tstring, idx)) }
	return w, r
}

func (c *client) Int(pattern string) (func(int), func() int) {
	idx := len(c.h[pattern].getInt)
	c.h[pattern].getInt = append(c.h[pattern].getInt, nil)
	c.h[pattern].postInt = append(c.h[pattern].postInt, nil)
	w := func(i int) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(i))
		c.write(pattern, tint, idx, b)
	}
	r := func() int { return int(binary.BigEndian.Uint64(c.read(pattern, tint, idx))) }
	return w, r
}
