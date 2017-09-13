package crest

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/mrmiguu/jsutil"
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
			client.Lock()
			client.b = true
			client.Unlock()
		}()
		err := http.ListenAndServe(url, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func get(w http.ResponseWriter, r *http.Request) {
	// this connection is part of an ephemeral session—
	// it will contribute to part of the connections that are
	// broadcasted to with respect to a channel write done on
	// the server's end

	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error; returning")
		return
	}
	println("/get", b)
	select {}

	parts := strings.Split(string(b), Sep)

	pattern, index := parts[0], parts[1]
	i, err := strconv.Atoi(index)
	if err != nil {
		println("error; returning")
		return
	}

	handlers.RLock()
	handlers.m[pattern].rbytes.RLock()
	b = <-handlers.m[pattern].rbytes.sl[i]
	handlers.m[pattern].rbytes.RUnlock()
	handlers.RUnlock()

	w.Write(b)
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error; returning")
		return
	}
	println("/post", b)
	select {}

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
		sl []chan []byte
	}
	rbytes struct {
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
	index := strconv.Itoa(len(h.wbytes.sl))
	h.wbytes.sl = append(h.wbytes.sl, w)
	h.rbytes.sl = append(h.rbytes.sl, r)
	h.wbytes.Unlock()
	h.rbytes.Unlock()

	// read
	go func() {
		client.RLock()
		if !client.b {
			client.RUnlock()
			return
		}
		client.RUnlock()

		for {
			s := h.pattern + Sep + index
			resp, err := http.Post(address+"/get", "text/plain", strings.NewReader(s))
			if err != nil {
				jsutil.Alert(err.Error())
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				jsutil.Alert(err.Error())
			}
			jsutil.Alert("/get ! " + string(b))
			r <- b
		}
	}()

	// write
	go func() {
		client.RLock()
		if !client.b {
			client.RUnlock()
			return
		}
		client.RUnlock()

		for b := range w {
			s := h.pattern + Sep + index + Sep + string(b)
			_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
			if err != nil {
				jsutil.Alert(err.Error())
			}
			jsutil.Alert("/post ! " + string(s))
		}
	}()

	return w, r
}
