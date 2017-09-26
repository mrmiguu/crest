package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

func newClient() endpoint {
	return &client{h: safeh{m: map[string]*Handler{}}}
}

// TODO: add thread safety
// TODO: add thread safety
func (c *client) Connect(addr string) {
	if i := strings.LastIndex(addr, "/"); i+1 == len(addr) {
		addr = addr[:i]
	}
	c.addr = addr
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
	var resp *http.Response
	for ok := true; ok; ok = (err != nil || resp.StatusCode > 299) {
		resp, err = http.Post(c.addr+Write, "text/plain", bytes.NewReader(b))
	}
}

func (c *client) read(pattern string, t byte, idx int) []byte {
	i := itob(idx)
	b := bytes.Join([][]byte{[]byte(pattern), []byte{t}, i}, v)
	var err error
	var resp *http.Response
	for ok := true; ok; ok = (err != nil || resp.StatusCode > 299) {
		resp, err = http.Post(c.addr+Read, "text/plain", bytes.NewReader(b))
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
	h.postBytes.sl = append(h.postBytes.sl, &getbytes{c: make(chan []byte, n)})
	h.getBytes.sl = append(h.getBytes.sl, &getbytes{c: make(chan []byte, n)})
	w := func(b []byte) {
		go func() {
			c.write(pattern, Tbytes, idx, b)
			<-h.postBytes.sl[idx].c
		}()
		h.postBytes.sl[idx].c <- nil
	}
	r := func() []byte { return c.read(pattern, Tbytes, idx) }
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
	h.postString.sl = append(h.postString.sl, &getstring{c: make(chan string, n)})
	h.getString.sl = append(h.getString.sl, &getstring{c: make(chan string, n)})
	w := func(s string) {
		go func() {
			c.write(pattern, Tstring, idx, []byte(s))
			<-h.postString.sl[idx].c
		}()
		h.postString.sl[idx].c <- ""
	}
	r := func() string { return string(c.read(pattern, Tstring, idx)) }
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
	h.postInt.sl = append(h.postInt.sl, &getint{c: make(chan int, n)})
	h.getInt.sl = append(h.getInt.sl, &getint{c: make(chan int, n)})
	w := func(i int) {
		go func() {
			c.write(pattern, Tint, idx, itob(i))
			<-h.postInt.sl[idx].c
		}()
		h.postInt.sl[idx].c <- 0
	}
	r := func() int { return btoi(c.read(pattern, Tint, idx)) }
	return w, r
}

func (c *client) Bool(pattern string, n int) (func(bool), func() bool) {
	c.h.RLock()
	h := c.h.m[pattern]
	h.postBool.Lock()
	h.getBool.Lock()
	defer c.h.RUnlock()
	defer h.postBool.Unlock()
	defer h.getBool.Unlock()

	idx := len(h.postBool.sl)
	h.postBool.sl = append(h.postBool.sl, &getbool{c: make(chan bool, n)})
	h.getBool.sl = append(h.getBool.sl, &getbool{c: make(chan bool, n)})
	w := func(b bool) {
		go func() {
			c.write(pattern, Tbool, idx, bool2bytes(b))
			<-h.postBool.sl[idx].c
		}()
		h.postBool.sl[idx].c <- true
	}
	r := func() bool { return bytes2bool(c.read(pattern, Tbool, idx)) }
	return w, r
}
