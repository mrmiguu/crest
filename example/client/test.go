package main

import (
	"encoding/json"
	"fmt"

	"github.com/mrmiguu/rest"
)

type test struct {
	this    string
	is      []bool
	working struct {
		period float64
	}
}

func main() {
	rest.Connect("http://127.0.0.1/")
	h := rest.New("test")
	w, r := h.Bytes()
	t := &test{
		this: "Th!s.",
		is:   []bool{false, true},
	}
	t.working.period = 420.69
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	w(b)
	var t2 test
	err = json.Unmarshal(r(), &t2)
	if err != nil {
		panic(err)
	}
	fmt.Println(t2)

	select {}
}
