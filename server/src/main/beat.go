package main

import (
	"net"
)

type BeatHandler struct {
	Handler
}

func init() {
	ch := BeatHandler{}
	ch.index = 0
	var ih IHandler = ch
	registerMap(&ih)
}

func (h BeatHandler) doHandle(sr Server, addr *net.UDPAddr, bs []byte) {
	sr.dict.Update(addr)
}
