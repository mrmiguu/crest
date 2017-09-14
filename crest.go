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
	go func() { log.Fatal(http.ListenAndServe(url, nil)) }()
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
	f  func([]byte)
	cb struct {
		sync.RWMutex
		sl []chan []byte
	}
}

// Bytes creates a byte slice channel split into write-only and read-only parts.
func (h *Handler) Bytes(buf ...int) (func([]byte), func() []byte) {
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}

	h.wbytes.Lock()
	h.rbytes.Lock()
	defer h.wbytes.Unlock()
	defer h.rbytes.Unlock()

	i := len(h.wbytes.sl)
	index := strconv.Itoa(i)

	w := func(b []byte) {
		if !isClient {
			println("w(?)...")
			defer println("w(?) !")
			h.wbytes.RLock()
			h.wbytes.sl[i].cb.Lock()
			if len(h.wbytes.sl[i].cb.sl) < 1 {
				h.wbytes.sl[i].cb.Unlock()
				h.wbytes.RUnlock()
				println("Sleeping...")
				<-reboot
				println("REBOOTING !")
				h.wbytes.RLock()
				h.wbytes.sl[i].cb.Lock()
			}
			for _, cb := range h.wbytes.sl[i].cb.sl {
				println("cb <- ...")
				cb <- b
				println("cb <- !")
			}
			h.wbytes.sl[i].cb.sl = []chan []byte{}
			h.wbytes.sl[i].cb.Unlock()
			h.wbytes.RUnlock()
		} else {
			println("w(?)...")
			defer println("w(?) !")
			s := h.pattern + V + index + V + string(b)
			for {
				println("/post...")
				_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
				if err == nil {
					println("/post !")
					break
				}
				println("RESTARTING...")
			}
		}
	}

	r := func() []byte {
		if !isClient {
			h.rbytes.RLock()
			defer h.rbytes.RUnlock()
			return <-h.rbytes.sl[i]
		}

		println("r()...")
		s := h.pattern + V + index
		var resp *http.Response
		var err error
		for {
			println("/get...")
			resp, err = http.Post(address+"/get", "text/plain", strings.NewReader(s))
			if err == nil {
				break
			}
			println("RESTARTING...")
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
		}
		println("/get !")
		println("r() !")
		return b
	}

	h.wbytes.sl = append(h.wbytes.sl, wbyte{f: w})
	h.rbytes.sl = append(h.rbytes.sl, make(chan []byte, n))

	// write
	// go func() {
	// 	if !isClient {
	// 		println("w <- ...")
	// 		for b := range w {
	// 			defer println("w <- !")
	// 			h.wbytes.RLock()
	// 			h.wbytes.sl[i].cb.Lock()
	// 			if len(h.wbytes.sl[i].cb.sl) < 1 {
	// 				h.wbytes.sl[i].cb.Unlock()
	// 				h.wbytes.RUnlock()
	// 				println("Sleeping...")
	// 				<-reboot
	// 				println("REBOOTING !")
	// 				h.wbytes.RLock()
	// 				h.wbytes.sl[i].cb.Lock()
	// 			}
	// 			for _, cb := range h.wbytes.sl[i].cb.sl {
	// 				println("cb <- ...")
	// 				cb <- b
	// 				println("cb <- !")
	// 			}
	// 			h.wbytes.sl[i].cb.sl = []chan []byte{}
	// 			h.wbytes.sl[i].cb.Unlock()
	// 			h.wbytes.RUnlock()
	// 		}
	// 	} else {
	// 		println("w <- ...")
	// 		for b := range w {
	// 			defer println("w <- !")
	// 			s := h.pattern + V + index + V + string(b)
	// 			for {
	// 				println("/post...")
	// 				_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
	// 				if err == nil {
	// 					println("/post !")
	// 					break
	// 				}
	// 				println("RESTARTING...")
	// 			}
	// 		}
	// 	}
	// }()

	// read
	// go func() {
	// 	if !isClient {
	// 		return
	// 	}
	// 	for {
	// 		s := h.pattern + V + index
	// 		println("/get...")
	// 		resp, err := http.Post(address+"/get", "text/plain", strings.NewReader(s))
	// 		if err != nil {
	// 			println("RESTARTING...")
	// 			continue
	// 		}
	// 		defer resp.Body.Close()
	// 		b, err := ioutil.ReadAll(resp.Body)
	// 		if err != nil {
	// 			println(err.Error())
	// 		}
	// 		println("/get !")
	// 		println("<-r...")
	// 		r <- b
	// 		println("<-r !")
	// 	}
	// }()

	return w, r
}
