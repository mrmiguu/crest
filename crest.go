package crest

import (
	"fmt"
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
	mux.HandleFunc(Endpoint+"/w", wserver)
	mux.HandleFunc(Endpoint+"/r", rserver)

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

func wserver(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error; returning")
		return
	}

	parts := strings.Split(string(b), Sep)
	fmt.Println("/w", parts)

	pattern, index := parts[0], parts[1]
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println("error; returning")
		return
	}

	handlers.RLock()
	handlers.m[pattern].rbytes.RLock()
	b = <-handlers.m[pattern].rbytes.sl[i]
	handlers.m[pattern].rbytes.RUnlock()
	handlers.RUnlock()

	w.Write(b)
}

func rserver(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error; returning")
		return
	}
	parts := strings.Split(string(b), Sep)
	fmt.Println("/r", parts)

	pattern, index, bytes := parts[0], parts[1], parts[2]
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println("error; returning")
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
			_, err := http.Post(address+"/r", "text/html", strings.NewReader(s))
			if err != nil {
				jsutil.Alert(err.Error())
			}
			jsutil.Alert("/w ! " + string(s))
		}
	}()

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
			resp, err := http.Post(address+"/w", "text/html", strings.NewReader(s))
			if err != nil {
				jsutil.Alert(err.Error())
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				jsutil.Alert(err.Error())
			}
			jsutil.Alert("/r ! " + string(b))
			r <- b
		}
	}()

	return w, r
}
