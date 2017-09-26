package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func newServer() endpoint {
	s := &server{safeh{m: map[string]*Handler{}}}
	go s.run()
	return s
}

// TODO: add thread safety
// TODO: add thread safety
func (s *server) Connect(addr string) {
	go http.ListenAndServe(addr, nil)
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

func (s *server) run() {
	http.HandleFunc(Write, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parts := bytes.Split(b, v)
		if len(parts) != 4 {
			http.Error(w, "bad pattern-type-index-message format", http.StatusBadRequest)
			return
		}
		pattern, t, idx, msg := string(parts[0]), parts[1][0], btoi(parts[2]), parts[3]

		s.h.RLock()
		defer s.h.RUnlock()
		h, exists := s.h.m[pattern]
		if !exists {
			http.Error(w, "pattern does not exist", http.StatusNotFound)
			return
		}

		switch t {
		case Tbytes:
			h.postBytes.RLock()
			defer h.postBytes.RUnlock()
			if idx > len(h.postBytes.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.postBytes.sl[idx].c <- msg

		case Tstring:
			h.postString.RLock()
			defer h.postString.RUnlock()
			if idx > len(h.postString.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.postString.sl[idx].c <- string(msg)

		case Tint:
			h.postInt.RLock()
			defer h.postInt.RUnlock()
			if idx > len(h.postInt.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.postInt.sl[idx].c <- btoi(msg)

		case Tbool:
			h.postBool.RLock()
			defer h.postBool.RUnlock()
			if idx > len(h.postBool.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.postBool.sl[idx].c <- bytes2bool(msg)

		default:
			http.Error(w, "type does not exist", http.StatusNotFound)
		}
	})

	http.HandleFunc(Read, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parts := bytes.Split(b, v)
		if len(parts) != 3 {
			http.Error(w, "bad pattern-type-index format", http.StatusBadRequest)
			return
		}
		pattern, t, idx := string(parts[0]), parts[1][0], btoi(parts[2])

		s.h.RLock()
		defer s.h.RUnlock()
		h, exists := s.h.m[pattern]
		if !exists {
			http.Error(w, "pattern does not exist", http.StatusNotFound)
			return
		}

		switch t {
		case Tbytes:
			h.getBytes.RLock()
			defer h.getBytes.RUnlock()
			if idx > len(h.getBytes.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.getBytes.sl[idx].n <- 1
			b = <-h.getBytes.sl[idx].c

		case Tstring:
			h.getString.RLock()
			defer h.getString.RUnlock()
			if idx > len(h.getString.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.getString.sl[idx].n <- 1
			b = []byte(<-h.getString.sl[idx].c)

		case Tint:
			h.getInt.RLock()
			defer h.getInt.RUnlock()
			if idx > len(h.getInt.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.getInt.sl[idx].n <- 1
			b = itob(<-h.getInt.sl[idx].c)

		case Tbool:
			h.getBool.RLock()
			defer h.getBool.RUnlock()
			if idx > len(h.getBool.sl)-1 {
				http.Error(w, "index does not exist", http.StatusNotFound)
				return
			}
			h.getBool.sl[idx].n <- 1
			b = bool2bytes(<-h.getBool.sl[idx].c)

		default:
			http.Error(w, "type does not exist", http.StatusNotFound)
		}

		w.Write(b)
	})
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
		for {
			<-get.n
			get.c <- b
			if len(get.n) < 1 {
				return
			}
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
		for {
			<-get.n
			get.c <- x
			if len(get.n) < 1 {
				return
			}
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
		for {
			<-get.n
			get.c <- i
			if len(get.n) < 1 {
				return
			}
		}
	}

	r := func() int { return <-h.postInt.sl[idx].c }

	return w, r
}

func (s *server) Bool(pattern string, n int) (func(bool), func() bool) {
	s.h.RLock()
	h := s.h.m[pattern]
	h.getBool.Lock()
	h.postBool.Lock()
	defer s.h.RUnlock()
	defer h.getBool.Unlock()
	defer h.postBool.Unlock()

	idx := len(h.getBool.sl)
	h.getBool.sl = append(h.getBool.sl, &getbool{make(chan int, xreads), make(chan bool, n)})
	h.postBool.sl = append(h.postBool.sl, &getbool{c: make(chan bool, n)})

	get := h.getBool.sl[idx]
	w := func(b bool) {
		for {
			<-get.n
			get.c <- b
			if len(get.n) < 1 {
				return
			}
		}
	}

	r := func() bool { return <-h.postBool.sl[idx].c }

	return w, r
}
