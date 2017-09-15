package rest

import (
	"encoding/binary"
	"io/ioutil"
	"net/http"
)

func newServer(addr string) *server {
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
		t := b[0]
		b = b[1:]
		switch t {
		case tbytes:
			s.postBytes <- b
		case tstring:
			s.postString <- string(b)
		case tint:
			s.postInt <- int(binary.BigEndian.Uint64(b))
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		t := b[0]
		switch t {
		case tbytes:
			b = <-s.getBytes
		case tstring:
			b = []byte(<-s.getString)
		case tint:
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(<-s.getInt))
		}
		w.Write(b)
	})

	http.ListenAndServe(addr, nil)
}

// Bytes creates a byte slice REST channel.
func (s *server) Bytes() (func([]byte), func() []byte) {
	s.getBytes = make(chan []byte)
	s.postBytes = make(chan []byte)
	w := func(b []byte) { s.getBytes <- b }
	r := func() []byte { return <-s.postBytes }
	return w, r
}

// String creates a string REST channel.
func (s *server) String() (func(string), func() string) {
	s.getString = make(chan string)
	s.postString = make(chan string)
	w := func(x string) { s.getString <- x }
	r := func() string { return <-s.postString }
	return w, r
}

// Int creates an int REST channel.
func (s *server) Int() (func(int), func() int) {
	s.getInt = make(chan int)
	s.postInt = make(chan int)
	w := func(i int) { s.getInt <- i }
	r := func() int { return <-s.postInt }
	return w, r
}
