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

func (h BeatHandler) doHandle(sr Server, addr *net.UDPAddr, no int64, bs []byte) {
	sr.dict.Update(addr)
	sr.Write(0, empty, addr)
}
