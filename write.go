package crest

import (
	"net/http"
	"strconv"
	"strings"
)

func wbytesf(h *Handler, t string) func(interface{}) {
	i := len(h.wbytes.sl)
	index := strconv.Itoa(i)

	return func(x interface{}) {
		if !isClient {
			// println("w(?)...")
			// defer println("w(?) !")
			h.wbytes.RLock()
			h.wbytes.sl[i].Lock()
			if len(h.wbytes.sl[i].sl) < 1 {
				h.wbytes.sl[i].Unlock()
				h.wbytes.RUnlock()
				println("Sleeping...")
				<-reboot
				println("REBOOTING !")
				h.wbytes.RLock()
				h.wbytes.sl[i].Lock()
			}
			for _, cb := range h.wbytes.sl[i].sl {
				// println("cb <- ...")
				cb <- x
				// println("cb <- !")
			}
			h.wbytes.sl[i].sl = make([]chan interface{}, 0) // change 0 to 1??
			h.wbytes.sl[i].Unlock()
			h.wbytes.RUnlock()
		} else {
			// println("w(?)...")
			// defer println("w(?) !")
			var msg string
			switch t {
			case tint:
				msg = strconv.Itoa(x.(int))
			case tstring:
				msg = x.(string)
			case tbytes:
				msg = string(x.([]byte))
			default:
				println("unknown type")
			}
			s := h.pattern + V + t + V + index + V + msg
			for {
				println("/post...")
				_, err := http.Post(address+"/post", "text/plain", strings.NewReader(s))
				if err == nil {
					println("/post !")
					break
				}
				println("RESTARTING...")
			}
		}
	}
}
