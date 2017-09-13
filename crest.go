package crest

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	Endpoint = "" //"/83c66fb4ee98eddef7cba94d787e4dc135e70a95"
	Sep      = "▼"
)

var (
	address string

	client struct {
		sync.RWMutex
		b bool
	}

	handlers = struct {
		sync.RWMutex
		m map[string]*Handler
	}{m: map[string]*Handler{}}
)

func Connect(url string) {
	address = url + Endpoint

	mux := http.NewServeMux()
	mux.HandleFunc(Endpoint+"/get", get)
	mux.HandleFunc(Endpoint+"/post", post)

	go func() {
		defer func() {
			recover()

			//
			// THIS IS TRUE BEFORE IT'S LOCK-N-READ BY h.Bytes(); it's false even though it shouldn't be
			// THIS IS TRUE BEFORE IT'S LOCK-N-READ BY h.Bytes(); it's false even though it shouldn't be
			// THIS IS TRUE BEFORE IT'S LOCK-N-READ BY h.Bytes(); it's false even though it shouldn't be
			//
			client.Lock()
			client.b = true
			println(client.b)
			client.Unlock()
		}()
		err := http.ListenAndServe(url, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func get(w http.ResponseWriter, r *http.Request) {
	println("/get")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error; returning")
		return
	}

	parts := strings.Split(string(b), Sep)

	pattern, index := parts[0], parts[1]
	i, err := strconv.Atoi(index)
	if err != nil {
		println("error; returning")
		return
	}

	cb := make(chan []byte)

	handlers.RLock()
	handlers.m[pattern].wbytes.RLock()
	handlers.m[pattern].wbytes.sl[i].cb.Lock()
	handlers.m[pattern].wbytes.sl[i].cb.sl = append(handlers.m[pattern].wbytes.sl[i].cb.sl, cb)
	handlers.m[pattern].wbytes.sl[i].cb.Unlock()
	handlers.m[pattern].wbytes.RUnlock()
	handlers.RUnlock()

	w.Write(<-cb)
}

func post(w http.ResponseWriter, r *http.Request) {
	println("/post")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error; returning")
		return
	}

	parts := strings.Split(string(b), Sep)

	pattern, index, bytes := parts[0], parts[1], parts[2]
	i, err := strconv.Atoi(index)
	if err != nil {
		println("error; returning")
		return
	}

	handlers.RLock()
	handlers.m[pattern].rbytes.RLock()
	handlers.m[pattern].rbytes.sl[i] <- []byte(bytes)
	handlers.m[pattern].rbytes.RUnlock()
	handlers.RUnlock()
}

func New(pattern string) *Handler {
	h := &Handler{pattern: pattern}
	return h
}

type Handler struct {
	pattern string

	wbytes struct {
		sync.RWMutex
		sl []wbyte
	}
	rbytes struct {
		sync.RWMutex
		sl []chan []byte
	}
}

type wbyte struct {
	c  chan []byte
	cb struct {
		sync.RWMutex
		sl []chan []byte
	}
}

func (h *Handler) Bytes(buf ...int) (chan<- []byte, <-chan []byte) {
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}

	w := make(chan []byte, n)
	r := make(chan []byte, n)

	h.wbytes.Lock()
	h.rbytes.Lock()
	i := len(h.wbytes.sl)
	index := strconv.Itoa(i)
	h.wbytes.sl = append(h.wbytes.sl, wbyte{c: w})
	h.rbytes.sl = append(h.rbytes.sl, r)
	h.wbytes.Unlock()
	h.rbytes.Unlock()

	client.RLock()
	notServer := client.b
	client.RUnlock()
	println("notServer=", notServer)

	// write
	go func() {
		if !notServer {
			println("for b := range w...")
			for b := range w {
				h.wbytes.RLock()
				h.wbytes.sl[i].cb.Lock()
				for _, cb := range h.wbytes.sl[i].cb.sl {
					print("cb <- b... ")
					cb <- b
					println("!")
				}
				h.wbytes.sl[i].cb.sl = []chan []byte{}
				h.wbytes.sl[i].cb.Unlock()
				h.wbytes.RUnlock()
			}
		} else {
			println("for b := range w...")
			for b := range w {
				println("b! read from writer ch")
				s := h.pattern + Sep + index + Sep + string(b)
				for {
					println("http.Post:POST...")
					_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
					if err == nil {
						break
					}
				}
				println("/post ! " + string(s))
			}
		}
	}()

	// read
	go func() {
		if !notServer {
			return
		}
		for {
			s := h.pattern + Sep + index
			println("http.Post:GET...")
			resp, err := http.Post(address+"/get", "text/plain", strings.NewReader(s))
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				println(err.Error())
			}
			println("/get ! " + string(b))
			r <- b
		}
	}()

	return w, r
}
