package main

import (
	"fmt"
	"net"
)

type LinkHandler struct {
	Handler
}

func init() {
	ch := LinkHandler{}
	ch.index = 3
	var ih IHandler = ch
	registerMap(&ih)
}

func (h LinkHandler) doHandle(sr Server, addr *net.UDPAddr, no int64, bs []byte) {
	if bs == nil {
		return
	}
	key := string(bs[0:4])
	if key == "link" {
		ipAndPort := string(bs[4:])
		dist := addr.String()
		if ipAndPort == dist {
			fmt.Println("link source target same")
		} else {
			sr.dict.notify(no, h.index, sr.udpConn, string(bs[4:]), addr.String())
		}
	}
}
