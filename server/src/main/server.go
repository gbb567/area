package main

import (
	"container/list"
	"fmt"
	"net"
	"time"
)

var beat = []byte("beat")
var dur = 2 * time.Second
var unKnow = -99

type udpConfig struct {
	network  string
	addr     string
	protocol string
}
type Server struct {
	udpConfig udpConfig
	dict      ICatalog
	udpConn   *net.UDPConn
}

func NewServer() *Server {
	var server Server
	udpC := udpConfig{network: "udp4", addr: "127.0.0.1:9998", protocol: "udp"}
	server = Server{udpConfig: udpC, dict: NewCatalog()}
	return &server
}

func (sr *Server) Bind() {
	udpAddr, _ := net.ResolveUDPAddr(sr.udpConfig.network, sr.udpConfig.addr)
	udpConn, err2 := net.ListenUDP(sr.udpConfig.protocol, udpAddr)
	if err2 != nil {
		panic(err2)
	}
	(*sr).udpConn = udpConn
	go makeTimer(udpConn, *sr)
	handleUdp(*sr)
}

func (sr *Server) Read() (int, int64, []byte, *net.UDPAddr) {
	bs := make([]byte, 128)
	len, addr, err := sr.udpConn.ReadFromUDP(bs)
	if err != nil {
		fmt.Println(err)
		return unKnow, 0, empty, nil
	}
	i, no, b := decode(len, bs)
	return i, no, b, addr
}

func (sr *Server) Write(no int64, bs []byte, addr *net.UDPAddr) {
	_, err := sr.udpConn.WriteToUDP(encode(no, bs), addr)
	if err != nil {
		fmt.Println(err)
	}
}

func makeTimer(udpConn *net.UDPConn, sr Server) {
	tiker := time.NewTicker(dur)
	for {
		select {
		case <-tiker.C:
			var j *list.Element
			now := time.Now()
			for i := sr.dict.List().Front(); i != nil; {
				j = i.Next()
				wrap := i.Value.(*UDPAddrWrap)
				if now.Sub(i.Value.(*UDPAddrWrap).time) > dur {
					sr.dict.Remove(wrap.addr)
				}
				i = j
			}
		}
	}
}
