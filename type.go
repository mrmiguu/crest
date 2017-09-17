package rest

type endpoint interface {
	New(string) *Handler
	Bytes(string) (func([]byte), func() []byte)
	String(string) (func(string), func() string)
	Int(string) (func(int), func() int)
}

type server struct {
	h map[string]*Handler
}

type client struct {
	addr string
	h    map[string]*Handler
}

// Handler holds pattern-relative typed channels.
type Handler struct {
	hptr                  *map[string]*Handler
	pattern               string
	getBytes, postBytes   []chan []byte
	getString, postString []chan string
	getInt, postInt       []chan int
}
