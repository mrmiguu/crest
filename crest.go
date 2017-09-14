package crest

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Connect prepares communication with the specified URL.
func Connect(url string) {
	if strings.LastIndex(url, "/")+1 == len(url) {
		url = url[:len(url)-1]
	}
	address = url + Endpoint

	isClient = isClientExpr.MatchString(url)
	if isClient {
		return
	}

	http.HandleFunc(Endpoint+"/get", get)
	http.HandleFunc(Endpoint+"/post", post)
	go log.Fatal(http.ListenAndServe(url, nil))
}

// New creates a new pattern-specific handler for creating REST channels.
func New(pattern string) *Handler {
	h := &Handler{pattern: pattern}
	handlers.Lock()
	handlers.m[pattern] = h
	handlers.Unlock()
	return h
}

// Handler holds pattern-specific read/write channels.
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

// Bytes creates a byte slice channel split into write-only and read-only parts.
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

	// write
	go func() {
		if !isClient {
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
				s := h.pattern + V + index + V + string(b)
				for {
					println("http.Post:POST...")
					_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
					if err == nil {
						break
					}
					print("RESTARTING ! ")
				}
				println("/post ! " + string(s))
			}
		}
	}()

	// read
	go func() {
		if !isClient {
			return
		}
		for {
			s := h.pattern + V + index
			println("http.Post:GET...")
			resp, err := http.Post(address+"/get", "text/plain", strings.NewReader(s))
			if err != nil {
				print("RESTARTING ! ")
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
