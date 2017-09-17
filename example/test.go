package main

import "github.com/mrmiguu/rest"
import "time"

func main() {
	rest.Connect("127.0.0.1:80")
	h := rest.New("login")
	w, r := h.Int()
	println("r()...")
	time.Sleep(3 * time.Second)
	println(r())
	println("w(420)...")
	time.Sleep(3 * time.Second)
	w(420)
	h.Close()
	h2 := rest.New("login")
	ws, rs := h2.String()
	println(`ws("You just got my '420'")...`)
	time.Sleep(3 * time.Second)
	ws("You just got my '420'")
	println("rs()...")
	time.Sleep(3 * time.Second)
	println(rs())
	select {}
}
