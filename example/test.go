package main

import "github.com/mrmiguu/rest"

func main() {
	rest.Connect("127.0.0.1:80")
	h := rest.New("login")
	w, r := h.Int()
	println(r())
	w(420)
	h2 := rest.New("logout")
	ws, rs := h2.String()
	ws("You just got my '420'")
	println(rs())
	select {}
}

// package main

// import (
// 	"encoding/binary"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// )

// func main() {
// 	go run()
// 	w, r := String()
// 	msg := r()
// 	println(msg)
// 	w(strings.ToUpper(msg))
// }

// const (
// 	tbytes byte = iota
// 	tstring
// 	tint
// )

// func run() {
// 	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		b, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			panic(err)
// 		}
// 		t := b[0]
// 		b = b[1:]
// 		switch t {
// 		case tbytes:
// 			postBytes <- b
// 		case tstring:
// 			postString <- string(b)
// 		case tint:
// 			postInt <- int(binary.BigEndian.Uint64(b))
// 		}
// 	})

// 	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		b, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			panic(err)
// 		}
// 		t := b[0]
// 		switch t {
// 		case tbytes:
// 			b = <-getBytes
// 		case tstring:
// 			b = []byte(<-getString)
// 		case tint:
// 			b = make([]byte, 8)
// 			binary.BigEndian.PutUint64(b, uint64(<-getInt))
// 		}
// 		w.Write(b)
// 	})

// 	http.ListenAndServe("127.0.0.1:80", nil)
// }

// var (
// 	getBytes   = make(chan []byte)
// 	postBytes  = make(chan []byte)
// 	getString  = make(chan string)
// 	postString = make(chan string)
// 	getInt     = make(chan int)
// 	postInt    = make(chan int)
// )

// func Bytes() (func([]byte), func() []byte) {
// 	w := func(b []byte) { getBytes <- b }
// 	r := func() []byte { return <-postBytes }
// 	return w, r
// }

// func String() (func(string), func() string) {
// 	w := func(s string) { getString <- s }
// 	r := func() string { return <-postString }
// 	return w, r
// }

// func Int() (func(int), func() int) {
// 	w := func(i int) { getInt <- i }
// 	r := func() int { return <-postInt }
// 	return w, r
// }
