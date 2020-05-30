package main

import (
	"fmt"
	"net"
)

type IHandler interface {
	getIndex() int
	doHandle(sr Server, addr *net.UDPAddr, bs []byte)
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
	buf := make([]byte, 128)
	for {
		len, udpAddr, err := sr.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len < 4 {
			fmt.Println("read error")
			return
		}
		h, ok := udpHandler.hmap[makeLen(buf[0:4])]
		if !ok {
			fmt.Println("error udp read")
			return
		}
		if len > 4 {
			go (*h).doHandle(sr, udpAddr, buf[4:len])
		} else {
			go (*h).doHandle(sr, udpAddr, nil)
		}
	}
}

func makeLen(bs []byte) int {
	return (int(bs[0]) << 24) | (int(bs[1]) << 16) | (int(bs[2]) << 8) | int(bs[3])
}
