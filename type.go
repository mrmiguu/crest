package rest

type endpoint interface {
	Bytes(string) (func([]byte), func() []byte)
	String(string) (func(string), func() string)
	Int(string) (func(int), func() int)
}

type server struct {
	h map[string]*Handler
}

type client struct {
	h    map[string]*Handler
	addr string
}

// Handler holds pattern-relative typed channels.
type Handler struct {
	pattern    string
	getBytes   chan []byte
	postBytes  chan []byte
	getString  chan string
	postString chan string
	getInt     chan int
	postInt    chan int
}
