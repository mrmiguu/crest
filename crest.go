package crest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	Endpoint = "/83c66fb4ee98eddef7cba94d787e4dc135e70a95"
)

var (
	address string
)

func Connect(url string) {
	address = url + Endpoint

	// this won't be used if not a server
	mux := http.NewServeMux()

	mux.HandleFunc(Endpoint+"/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "")
	})

	mux.HandleFunc(Endpoint+"/post", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "")
	})

	go func() {
		fmt.Println("Server?")
		defer func() {
			e := recover()
			if e == nil {
				return
			}
			fmt.Println("Client!")
		}()
		err := http.ListenAndServe(url, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

type Handler struct {
	rstring struct {
		sync.RWMutex
		sl []chan string
	}
}

func New(pattern string) *Handler {
	var h Handler
	return &h
}

func (h *Handler) String(buf ...int) (chan<- string, <-chan string) {
	n := 0
	if len(buf) > 0 {
		n = buf[0]
	}
	w := make(chan string, n)
	r := make(chan string, n)

	// client read
	go func() {
		for {
			resp, err := http.Get(address)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			r <- string(b)
		}
	}()

	return w, r
}
