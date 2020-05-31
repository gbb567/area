package main

import (
	"fmt"
	"net"
)

var b1 = []byte("1")
var b2 = []byte("2")

type NatHandler struct {
	Handler
}

func init() {
	ch := NatHandler{}
	ch.index = 2
	var ih IHandler = ch
	registerMap(&ih)
}

func (h NatHandler) doHandle(sr Server, addr *net.UDPAddr, no int64, bs []byte) {
	if bs == nil {
		return
	}
	key := string(bs)
	if key == "check" {
		//check Restricted Cone
		con, err2 := net.DialUDP("udp", nil, addr)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		con.Write(encode(no, b2))
		//check Full Cone,need two public ip
	}
}
