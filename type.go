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
	addr    string
	h       map[string]*Handler
	bytesc  int
	stringc int
	intc    int
}

// Handler holds pattern-relative typed channels.
type Handler struct {
	pattern    string
	getBytes   []chan []byte
	postBytes  []chan []byte
	getString  []chan string
	postString []chan string
	getInt     []chan int
	postInt    []chan int
}
