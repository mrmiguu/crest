package crest

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	println("/get...")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("/get !", err.Error())
		return
	}

	parts := strings.Split(string(b), V)

	pattern, t, index := parts[0], parts[1], parts[2]
	i, err := strconv.Atoi(index)
	if err != nil {
		println("/get !", err.Error())
		return
	}

	cb := make(chan interface{})
	handlers.RLock()
	handlers.m[pattern].wbytes.RLock()
	handlers.m[pattern].wbytes.sl[i].Lock()
	if len(handlers.m[pattern].wbytes.sl[i].sl) < 1 {
		reboot <- true
	}
	handlers.m[pattern].wbytes.sl[i].sl = append(handlers.m[pattern].wbytes.sl[i].sl, cb)
	handlers.m[pattern].wbytes.sl[i].Unlock()
	handlers.m[pattern].wbytes.RUnlock()
	handlers.RUnlock()

	switch t {
	case tint:
		w.Write([]byte(strconv.Itoa((<-cb).(int))))
	case tstring:
		w.Write([]byte((<-cb).(string)))
	case tbytes:
		w.Write((<-cb).([]byte))
	default:
		println("unknown type")
	}

	println("/get !")
}

func post(w http.ResponseWriter, r *http.Request) {
	println("/post...")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("/post !", err.Error())
		return
	}

	parts := strings.Split(string(b), V)

	pattern, t, index, msg := parts[0], parts[1], parts[2], parts[3]
	i, err := strconv.Atoi(index)
	if err != nil {
		println("/post !", err.Error())
		return
	}

	handlers.RLock()
	handlers.m[pattern].rbytes.RLock()
	switch t {
	case tint:
		i, err := strconv.Atoi(msg)
		if err != nil {
			println("could not convert incoming to int")
		}
		handlers.m[pattern].rbytes.sl[i] <- i
	case tstring:
		handlers.m[pattern].rbytes.sl[i] <- msg
	case tbytes:
		handlers.m[pattern].rbytes.sl[i] <- []byte(msg)
	default:
		println("unknown type")
	}
	handlers.m[pattern].rbytes.RUnlock()
	handlers.RUnlock()

	println("/post !")
}
