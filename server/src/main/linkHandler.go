package main

import (
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

func (h LinkHandler) doHandle(sr Server, addr *net.UDPAddr, bs []byte) {
	if bs == nil {
		return
	}
	key := string(bs[0:4])
	if key == "link" {
		sr.dict.notify(sr.udpConn, string(bs[4:]), addr.String())
	}
}
