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
		h, exists := s.h.m[pattern]
		if !exists {
			http.Error(w, "pattern does not exist", http.StatusNotFound)
			return
		}
		switch t {
		case tbytes:
			h.postBytes.RLock()
			h.postBytes.sl[idx].c <- msg
			h.postBytes.RUnlock()
		case tstring:
			h.postString.RLock()
			h.postString.sl[idx].c <- string(msg)
			h.postString.RUnlock()
		case tint:
			h.postInt.RLock()
			h.postInt.sl[idx].c <- btoi(msg)
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
		h, exists := s.h.m[pattern]
		if !exists {
			http.Error(w, "pattern does not exist", http.StatusNotFound)
			return
		}

		switch t {
		case tbytes:
			h.getBytes.RLock()
			v := h.getBytes.sl[idx]
			h.getBytes.RUnlock()
			v.n <- 1
			b = <-v.c

		case tstring:
			h.getString.RLock()
			v := h.getString.sl[idx]
			h.getString.RUnlock()
			v.n <- 1
			b = []byte(<-v.c)

		case tint:
			h.getInt.RLock()
			v := h.getInt.sl[idx]
			h.getInt.RUnlock()
			v.n <- 1
			b = itob(<-v.c)
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
	h.getBytes.sl = append(h.getBytes.sl, &getbytes{make(chan int, xreads), make(chan []byte, n)})
	h.postBytes.sl = append(h.postBytes.sl, &getbytes{c: make(chan []byte, n)})

	get := h.getBytes.sl[idx]
	w := func(b []byte) {
		for z := len(get.n); z > 0; z-- {
			<-get.n
			get.c <- b
		}
	}

	r := func() []byte { return <-h.postBytes.sl[idx].c }

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
	h.getString.sl = append(h.getString.sl, &getstring{make(chan int, xreads), make(chan string, n)})
	h.postString.sl = append(h.postString.sl, &getstring{c: make(chan string, n)})

	get := h.getString.sl[idx]
	w := func(x string) {
		for z := len(get.n); z > 0; z-- {
			<-get.n
			get.c <- x
		}
	}

	r := func() string { return <-h.postString.sl[idx].c }

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
	h.getInt.sl = append(h.getInt.sl, &getint{make(chan int, xreads), make(chan int, n)})
	h.postInt.sl = append(h.postInt.sl, &getint{c: make(chan int, n)})

	get := h.getInt.sl[idx]
	w := func(i int) {
		for z := len(get.n); z > 0; z-- {
			<-get.n
			get.c <- i
		}
	}

	r := func() int { return <-h.postInt.sl[idx].c }

	return w, r
}
