package crest

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func rbytesf(h *Handler, t string) func() interface{} {
	i := len(h.wbytes.sl)
	index := strconv.Itoa(i)

	return func() interface{} {
		if !isClient {
			h.rbytes.RLock()
			defer h.rbytes.RUnlock()
			switch t {
			case tint:
				return (<-h.rbytes.sl[i]).(int)
			case tstring:
				return (<-h.rbytes.sl[i]).(string)
			case tbytes:
				return (<-h.rbytes.sl[i]).([]byte)
			default:
				println("unknown type")
			}
		}

		// println("r()...")
		s := h.pattern + V + t + V + index
		var resp *http.Response
		var err error
		for {
			println("/get...")
			resp, err = http.Post(address+"/get", "text/plain", strings.NewReader(s))
			if err == nil {
				break
			}
			println("RESTARTING...")
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
		}
		println("/get !")
		// println("r() !")

		switch t {
		case tint:
			i, err := strconv.Atoi(string(b))
			if err != nil {
				println("could not convert outgoing to int")
			}
			return i
		case tstring:
			return string(b)
		case tbytes:
			return b
		default:
			panic("unknown type")
		}
	}
}
