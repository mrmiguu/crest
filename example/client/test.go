package main

import (
	"time"

	"github.com/mrmiguu/rest"
)

func main() {
	rest.Connect("http://127.0.0.1/")
	h := rest.New("login")
	w, r := h.Int()
	println("w(69)...")
	time.Sleep(3 * time.Second)
	w(69)
	println("r()...")
	time.Sleep(3 * time.Second)
	println(r())
	h.Close()
	h2 := rest.New("login")
	ws, rs := h2.String()
	println("rs()...")
	time.Sleep(3 * time.Second)
	println(rs())
	println(`ws("You got my '69' earlier")...`)
	time.Sleep(3 * time.Second)
	ws("You got my '69' earlier")
	select {}
}
