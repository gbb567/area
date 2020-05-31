package main

import (
	"fmt"
	"net"
)

type IHandler interface {
	getIndex() int
	doHandle(sr Server, addr *net.UDPAddr, no int64, bs []byte)
}

type Handler struct {
	index int
}

type UdpHandler struct {
	hmap map[int]*IHandler
}

var udpHandler = UdpHandler{hmap: make(map[int]*IHandler)}

func (h Handler) getIndex() int {
	return h.index
}

func registerMap(ihandler *IHandler) {
	_, ok := udpHandler.hmap[(*ihandler).getIndex()]
	if !ok {
		udpHandler.hmap[(*ihandler).getIndex()] = ihandler
	} else {
		fmt.Println("已经存在的handler")
	}
}

//instruct len 4
func handleUdp(sr Server) {
	for {
		i, no, bs, udpAddr := sr.Read()
		if i == unKnow {
			continue
		}
		h, ok := udpHandler.hmap[i]
		if !ok {
			fmt.Println("error udp read")
		} else {
			go (*h).doHandle(sr, udpAddr, no, bs)
		}
	}
}
