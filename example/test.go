package main

import "github.com/mrmiguu/rest"
import "time"

func main() {
	rest.Connect("127.0.0.1:80")
	h := rest.New("login")
	w, r := h.Int(1)
	time.Sleep(1 * time.Second)
	println("w(420)...")
	w(420)
	time.Sleep(1 * time.Second)
	println("r()...")
	println(r())
	h.Close()
	h2 := rest.New("login")
	ws, rs := h2.String()
	time.Sleep(1 * time.Second)
	println(`ws("You just got my '420'")...`)
	ws("You just got my '420'")
	time.Sleep(1 * time.Second)
	println("rs()...")
	println(rs())
	select {}
}
