package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mrmiguu/rest"
)

type test struct {
	This    string
	Is      []bool
	Working struct {
		Period float64
	}
}

func main() {
	rest.Connect("http://127.0.0.1/")

	time.Sleep(3 * time.Second)
	h := rest.New("test")
	w, r := h.Bytes()
	var t test
	println("r()...")
	// r()
	err := json.Unmarshal(r(), &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	t2 := test{
		This: strings.ToUpper(t.This),
		Is:   []bool{t.Is[1], t.Is[0]},
	}
	t2.Working.Period = t.Working.Period / 2
	b, err := json.Marshal(t2)
	if err != nil {
		panic(err)
	}
	println("w(b)... [3 sec]")
	time.Sleep(3 * time.Second)
	w(b)

	select {}
}
