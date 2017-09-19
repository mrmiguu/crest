package main

import (
	"encoding/json"
	"fmt"

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
	rest.Connect("127.0.0.1:80")

	h := rest.New("test")
	w, r := h.Bytes()
	t := &test{
		This: "Th!s.",
		Is:   []bool{false, true},
	}
	t.Working.Period = 420.69
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	println("w(b)...")
	w(b)
	var t2 test
	println("r()...")
	err = json.Unmarshal(r(), &t2)
	if err != nil {
		panic(err)
	}
	fmt.Println(t2)

	select {}
}
