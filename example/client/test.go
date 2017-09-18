package main

import (
	"time"

	"github.com/mrmiguu/rest"
)

func main() {
	rest.Connect("http://127.0.0.1/")
	h := rest.New("login")
	w, r := h.Int(1)
	time.Sleep(1 * time.Second)
	println("r()...")
	println(r())
	time.Sleep(1 * time.Second)
	println("w(69)...")
	w(69)
	h.Close()
	h2 := rest.New("login")
	ws, rs := h2.String()
	time.Sleep(1 * time.Second)
	println("rs()...")
	println(rs())
	time.Sleep(1 * time.Second)
	println(`ws("You got my '69' earlier")...`)
	ws("You got my '69' earlier")
	select {}
}
