package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func newServer(addr string) endpoint {
	s := &server{safeh{m: map[string]*Handler{}}}
	go s.run(addr)
	return s
}

func (s *server) New(pattern string) *Handler {
	s.h.Lock()
	defer s.h.Unlock()
	if _, exists := s.h.m[pattern]; exists {
		panic("pattern already exists")
	}
	h := &Handler{hptr: &s.h, pattern: pattern}
	s.h.m[pattern] = h
	return h
}

func (s *server) run(addr string) {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		must(err)
		parts := bytes.Split(b, v)
		pattern, t, idx, msg := string(parts[0]), parts[1][0], btoi(parts[2]), parts[3]
		s.h.RLock()
		defer s.h.RUnlock()
		h := s.h.m[pattern]
		switch t {
		case tbytes:
			h.postBytes.RLock()
			h.postBytes.sl[idx] <- msg
			h.postBytes.RUnlock()
		case tstring:
			h.postString.RLock()
			h.postString.sl[idx] <- string(msg)
			h.postString.RUnlock()
		case tint:
			h.postInt.RLock()
			h.postInt.sl[idx] <- btoi(msg)
			h.postInt.RUnlock()
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		must(err)
		parts := bytes.Split(b, v)
		pattern, t, idx := string(parts[0]), parts[1][0], btoi(parts[2])
		s.h.RLock()
		defer s.h.RUnlock()
		h := s.h.m[pattern]
		switch t {
		case tbytes:
			h.getBytes.RLock()
			b = <-h.getBytes.sl[idx]
			h.getBytes.RUnlock()
		case tstring:
			h.getString.RLock()
			b = []byte(<-h.getString.sl[idx])
			h.getString.RUnlock()
		case tint:
			h.getInt.RLock()
			b = itob(<-h.getInt.sl[idx])
			h.getInt.RUnlock()
		}
		w.Write(b)
	})

	http.ListenAndServe(addr, nil)
}

func (s *server) Bytes(pattern string, n int) (func([]byte), func() []byte) {
	s.h.RLock()
	h := s.h.m[pattern]
	h.getBytes.Lock()
	h.postBytes.Lock()
	defer s.h.RUnlock()
	defer h.getBytes.Unlock()
	defer h.postBytes.Unlock()

	idx := len(h.getBytes.sl)
	h.getBytes.sl = append(h.getBytes.sl, make(chan []byte, n))
	h.postBytes.sl = append(h.postBytes.sl, make(chan []byte, n))
	w := func(b []byte) { h.getBytes.sl[idx] <- b }
	r := func() []byte { return <-h.postBytes.sl[idx] }
	return w, r
}

func (s *server) String(pattern string, n int) (func(string), func() string) {
	s.h.RLock()
	h := s.h.m[pattern]
	h.getString.Lock()
	h.postString.Lock()
	defer s.h.RUnlock()
	defer h.getString.Unlock()
	defer h.postString.Unlock()

	idx := len(h.getString.sl)
	h.getString.sl = append(h.getString.sl, make(chan string, n))
	h.postString.sl = append(h.postString.sl, make(chan string, n))
	w := func(x string) { h.getString.sl[idx] <- x }
	r := func() string { return <-h.postString.sl[idx] }
	return w, r
}

func (s *server) Int(pattern string, n int) (func(int), func() int) {
	s.h.RLock()
	h := s.h.m[pattern]
	h.getInt.Lock()
	h.postInt.Lock()
	defer s.h.RUnlock()
	defer h.getInt.Unlock()
	defer h.postInt.Unlock()

	idx := len(h.getInt.sl)
	h.getInt.sl = append(h.getInt.sl, make(chan int, n))
	h.postInt.sl = append(h.postInt.sl, make(chan int, n))
	w := func(i int) { h.getInt.sl[idx] <- i }
	r := func() int { return <-h.postInt.sl[idx] }
	return w, r
}
