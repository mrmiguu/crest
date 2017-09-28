package rest

import "sync"

type endpoint interface {
	Connect(string)
	New(string) *Handler
	Bytes(string, int) (func([]byte), func() []byte)
	String(string, int) (func(string), func() string)
	Int(string, int) (func(int), func() int)
	Bool(string, int) (func(bool), func() bool)
}

type server struct {
	h safeh
}

type client struct {
	addr string
	h    safeh
}

type safeh struct {
	sync.RWMutex
	m map[string]*Handler
}

// Handler holds pattern-relative typed channels.
type Handler struct {
	Pattern             string
	hptr                *safeh
	getBytes, postBytes struct {
		sync.RWMutex
		sl []*getbytes
	}
	getString, postString struct {
		sync.RWMutex
		sl []*getstring
	}
	getInt, postInt struct {
		sync.RWMutex
		sl []*getint
	}
	getBool, postBool struct {
		sync.RWMutex
		sl []*getbool
	}
}

type getbytes struct {
	n chan int
	c chan []byte
}

type getstring struct {
	n chan int
	c chan string
}

type getint struct {
	n chan int
	c chan int
}

type getbool struct {
	n chan int
	c chan bool
}
