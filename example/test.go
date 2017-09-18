package main

import (
	"encoding/json"
	"fmt"
	"strings"

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
	rest.Connect("127.0.0.1:80")
	h := rest.New("test")
	w, r := h.Bytes()
	var t test
	rd := r()
	fmt.Println(string(rd))
	err := json.Unmarshal(rd, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	t2 := test{
		this: strings.ToUpper(t.this),
		is:   []bool{t.is[1], t.is[0]},
	}
	t2.working.period = t.working.period / 2
	b, err := json.Marshal(t2)
	if err != nil {
		panic(err)
	}
	w(b)

	select {}
}
