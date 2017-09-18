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
	h := &Handler{hptr: &c.h.m, pattern: pattern}
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

func (c *client) Bytes(pattern string) (func([]byte), func() []byte) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.getBytes.Lock()
	h.postBytes.Lock()
	defer c.h.RUnlock()
	defer h.getBytes.Unlock()
	defer h.postBytes.Unlock()

	idx := len(h.getBytes.sl)
	h.getBytes.sl = append(h.getBytes.sl, nil)
	h.postBytes.sl = append(h.postBytes.sl, nil)
	w := func(b []byte) { c.write(pattern, tbytes, idx, b) }
	r := func() []byte { return c.read(pattern, tbytes, idx) }
	return w, r
}

func (c *client) String(pattern string) (func(string), func() string) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.getString.Lock()
	h.postString.Lock()
	defer c.h.RUnlock()
	defer h.getString.Unlock()
	defer h.postString.Unlock()

	idx := len(h.getString.sl)
	h.getString.sl = append(h.getString.sl, nil)
	h.postString.sl = append(h.postString.sl, nil)
	w := func(s string) { c.write(pattern, tstring, idx, []byte(s)) }
	r := func() string { return string(c.read(pattern, tstring, idx)) }
	return w, r
}

func (c *client) Int(pattern string) (func(int), func() int) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.getInt.Lock()
	h.postInt.Lock()
	defer c.h.RUnlock()
	defer h.getInt.Unlock()
	defer h.postInt.Unlock()

	idx := len(h.getInt.sl)
	h.getInt.sl = append(h.getInt.sl, nil)
	h.postInt.sl = append(h.postInt.sl, nil)
	w := func(i int) { c.write(pattern, tint, idx, itob(i)) }
	r := func() int { return btoi(c.read(pattern, tint, idx)) }
	return w, r
}
