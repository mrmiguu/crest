package crest

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	println("/get")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("error; returning")
		return
	}

	parts := strings.Split(string(b), V)

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

	parts := strings.Split(string(b), V)

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
