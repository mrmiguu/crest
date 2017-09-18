package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

func newClient(addr string) endpoint {
	if i := strings.LastIndex(addr, "/"); i+1 == len(addr) {
		addr = addr[:i]
	}
	return &client{h: safeh{m: map[string]*Handler{}}, addr: addr}
}

func (c *client) New(pattern string) *Handler {
	c.h.Lock()
	defer c.h.Unlock()
	if _, exists := c.h.m[pattern]; exists {
		panic("pattern already exists")
	}
	h := &Handler{hptr: &c.h, pattern: pattern}
	c.h.m[pattern] = h
	return h
}

func (c *client) write(pattern string, t byte, idx int, msg []byte) {
	i := itob(idx)
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, i, msg}, v)
	var err error
	for ok := true; ok; ok = (err != nil) {
		_, err = http.Post(c.addr+"/post", "text/plain", bytes.NewReader(b))
	}
}

func (c *client) read(pattern string, t byte, idx int) []byte {
	i := itob(idx)
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, i}, v)
	var err error
	var resp *http.Response
	for ok := true; ok; ok = (err != nil) {
		resp, err = http.Post(c.addr+"/get", "text/plain", bytes.NewReader(b))
	}
	b, err = ioutil.ReadAll(resp.Body)
	must(err)
	return b
}

func (c *client) Bytes(pattern string, n int) (func([]byte), func() []byte) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.postBytes.Lock()
	h.getBytes.Lock()
	defer c.h.RUnlock()
	defer h.postBytes.Unlock()
	defer h.getBytes.Unlock()

	idx := len(h.postBytes.sl)
	h.postBytes.sl = append(h.postBytes.sl, make(chan []byte, n))
	h.getBytes.sl = append(h.getBytes.sl, make(chan []byte, n))
	w := func(b []byte) {
		go func() { c.write(pattern, tbytes, idx, <-h.postBytes.sl[idx]) }()
		h.postBytes.sl[idx] <- b
	}
	r := func() []byte { return c.read(pattern, tbytes, idx) }
	return w, r
}

func (c *client) String(pattern string, n int) (func(string), func() string) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.postString.Lock()
	h.getString.Lock()
	defer c.h.RUnlock()
	defer h.postString.Unlock()
	defer h.getString.Unlock()

	idx := len(h.postString.sl)
	h.postString.sl = append(h.postString.sl, make(chan string, n))
	h.getString.sl = append(h.getString.sl, make(chan string, n))
	w := func(s string) {
		go func() { c.write(pattern, tstring, idx, []byte(<-h.postString.sl[idx])) }()
		h.postString.sl[idx] <- s
	}
	r := func() string { return string(c.read(pattern, tstring, idx)) }
	return w, r
}

func (c *client) Int(pattern string, n int) (func(int), func() int) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.postInt.Lock()
	h.getInt.Lock()
	defer c.h.RUnlock()
	defer h.postInt.Unlock()
	defer h.getInt.Unlock()

	idx := len(h.postInt.sl)
	h.postInt.sl = append(h.postInt.sl, make(chan int, n))
	h.getInt.sl = append(h.getInt.sl, make(chan int, n))
	w := func(i int) {
		go func() { c.write(pattern, tint, idx, itob(<-h.postInt.sl[idx])) }()
		h.postInt.sl[idx] <- i
	}
	r := func() int { return btoi(c.read(pattern, tint, idx)) }
	return w, r
}
