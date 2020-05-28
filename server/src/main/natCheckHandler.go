package main

import (
	"fmt"
	"net"
)

var m = make(map[string]map[int]*net.UDPAddr)
var b1 = []byte("1")
var b2 = []byte("2")
var b3 = []byte("3")

type NatCheckHandler struct {
	Handler
}

func (handler NatCheckHandler) handle(bs []byte, sc ServerConfig, conn net.Conn) {
	sc.codec.encode(sc.udpConfig.sbs, conn)
}

func natHandle(sc ServerConfig, conn *net.UDPConn, addr *net.UDPAddr) {
	//check Symmetric
	//check Port Restricted Cone
	_, err := conn.WriteToUDP(b3, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	//check Restricted Cone
	con, err2 := net.DialUDP("udp", nil, addr)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	con.Write(b2)
	ip4 := (*addr).String()
	m[ip4] = make(map[int]*net.UDPAddr)
	m[ip4][(*addr).Port] = addr
}
