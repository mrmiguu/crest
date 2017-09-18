package rest

import "sync"

type endpoint interface {
	New(string) *Handler
	Bytes(string, int) (func([]byte), func() []byte)
	String(string, int) (func(string), func() string)
	Int(string, int) (func(int), func() int)
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
	hptr                *safeh
	pattern             string
	getBytes, postBytes struct {
		sync.RWMutex
		sl []chan []byte
	}
	getString, postString struct {
		sync.RWMutex
		sl []chan string
	}
	getInt, postInt struct {
		sync.RWMutex
		sl []chan int
	}
}
