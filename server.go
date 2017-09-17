package rest

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"
)

func newServer(addr string) endpoint {
	s := &server{map[string]*Handler{}}
	go s.run(addr)
	return s
}

func (s *server) New(pattern string) *Handler {
	if _, exists := s.h[pattern]; exists {
		panic("pattern already exists")
	}
	h := &Handler{h: &s.h, pattern: pattern}
	s.h[pattern] = h
	return h
}

func (s *server) run(addr string) {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		parts := bytes.Split(b, v)
		pattern, t, idx, msg := string(parts[0]), parts[1][0], int(binary.BigEndian.Uint64(parts[2])), parts[3]
		switch t {
		case tbytes:
			s.h[pattern].postBytes[idx] <- msg
		case tstring:
			s.h[pattern].postString[idx] <- string(msg)
		case tint:
			s.h[pattern].postInt[idx] <- int(binary.BigEndian.Uint64(msg))
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		parts := bytes.Split(b, v)
		pattern, t, idx := string(parts[0]), parts[1][0], int(binary.BigEndian.Uint64(parts[2]))
		switch t {
		case tbytes:
			b = <-s.h[pattern].getBytes[idx]
		case tstring:
			b = []byte(<-s.h[pattern].getString[idx])
		case tint:
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(<-s.h[pattern].getInt[idx]))
		}
		w.Write(b)
	})

	http.ListenAndServe(addr, nil)
}

func (s *server) Bytes(pattern string) (func([]byte), func() []byte) {
	idx := len(s.h[pattern].getBytes)
	s.h[pattern].getBytes = append(s.h[pattern].getBytes, make(chan []byte))
	s.h[pattern].postBytes = append(s.h[pattern].postBytes, make(chan []byte))
	w := func(b []byte) { s.h[pattern].getBytes[idx] <- b }
	r := func() []byte { return <-s.h[pattern].postBytes[idx] }
	return w, r
}

func (s *server) String(pattern string) (func(string), func() string) {
	idx := len(s.h[pattern].getString)
	s.h[pattern].getString = append(s.h[pattern].getString, make(chan string))
	s.h[pattern].postString = append(s.h[pattern].postString, make(chan string))
	w := func(x string) { s.h[pattern].getString[idx] <- x }
	r := func() string { return <-s.h[pattern].postString[idx] }
	return w, r
}

func (s *server) Int(pattern string) (func(int), func() int) {
	idx := len(s.h[pattern].getInt)
	s.h[pattern].getInt = append(s.h[pattern].getInt, make(chan int))
	s.h[pattern].postInt = append(s.h[pattern].postInt, make(chan int))
	w := func(i int) { s.h[pattern].getInt[idx] <- i }
	r := func() int { return <-s.h[pattern].postInt[idx] }
	return w, r
}
