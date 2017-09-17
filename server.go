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
	h := &Handler{pattern: pattern}
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
		pattern, t, msg := string(parts[0]), parts[1][0], parts[2]
		switch t {
		case tbytes:
			s.h[pattern].postBytes <- msg
		case tstring:
			s.h[pattern].postString <- string(msg)
		case tint:
			s.h[pattern].postInt <- int(binary.BigEndian.Uint64(msg))
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		parts := bytes.Split(b, v)
		pattern, t := string(parts[0]), parts[1][0]
		switch t {
		case tbytes:
			b = <-s.h[pattern].getBytes
		case tstring:
			b = []byte(<-s.h[pattern].getString)
		case tint:
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(<-s.h[pattern].getInt))
		}
		w.Write(b)
	})

	http.ListenAndServe(addr, nil)
}

// Bytes creates a byte slice REST channel.
func (s *server) Bytes(pattern string) (func([]byte), func() []byte) {
	s.h[pattern].getBytes = make(chan []byte)
	s.h[pattern].postBytes = make(chan []byte)
	w := func(b []byte) { s.h[pattern].getBytes <- b }
	r := func() []byte { return <-s.h[pattern].postBytes }
	return w, r
}

// String creates a string REST channel.
func (s *server) String(pattern string) (func(string), func() string) {
	s.h[pattern].getString = make(chan string)
	s.h[pattern].postString = make(chan string)
	w := func(x string) { s.h[pattern].getString <- x }
	r := func() string { return <-s.h[pattern].postString }
	return w, r
}

// Int creates an int REST channel.
func (s *server) Int(pattern string) (func(int), func() int) {
	s.h[pattern].getInt = make(chan int)
	s.h[pattern].postInt = make(chan int)
	w := func(i int) { s.h[pattern].getInt <- i }
	r := func() int { return <-s.h[pattern].postInt }
	return w, r
}
