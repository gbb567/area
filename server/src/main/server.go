package main

import (
	"fmt"
	"net"
)

var connectError = []byte("connection was existing")

type Server interface {
	Bind()
}

type udpConfig struct {
	network  string
	addr     string
	protocol string
	sbs      []byte
}
type ServerConfig struct {
	Protocol  string
	Port      string
	ICatalog  *ICatalog
	codec     ICodec
	handler   MapHandler
	udpConfig udpConfig
}

func NewServer(protocol string, port string) *Server {
	var server Server
	udpC := udpConfig{network: "udp4", addr: "127.0.0.1:9998", protocol: "udp", sbs: []byte("udp4\\127.0.0.1:9998")}
	server = ServerConfig{Protocol: protocol, Port: port, ICatalog: NewCatalog(), codec: NewCodecer(), handler: NewHandler(), udpConfig: udpC}
	return &server
}

func (sc ServerConfig) Bind() {
	ln, err := net.Listen(sc.Protocol, ":"+sc.Port)
	if err != nil {
		panic(err)
	}
	udpAddr, _ := net.ResolveUDPAddr(sc.udpConfig.network, sc.udpConfig.addr)
	udpConn, err2 := net.ListenUDP(sc.udpConfig.protocol, udpAddr)
	if err2 != nil {
		panic(err2)
	}
	go handleUdp(sc, udpConn)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print("server accept error")
			continue
		}
		go handleConnection(sc, &conn)
	}
}

func handleUdp(sc ServerConfig, conn *net.UDPConn) {
	buf := make([]byte, 5)
	for {
		len, udpAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len != 5 && string(buf) != "hello" {
			fmt.Println("error udp read")
			return
		}
		go natHandle(sc, conn, udpAddr)
	}
}

func handleConnection(sc ServerConfig, conn *net.Conn) {
	if !(*sc.ICatalog).Append(conn) {
		sc.codec.encode(connectError, *conn)
		fmt.Println("connect error,this connection was existing")
		(*conn).Close()
		return
	}
	for {
		bs := sc.codec.decode(*conn)
		if bs == nil {
			fmt.Println("close this connection")
			(*sc.ICatalog).Remove(conn)
			(*conn).Close()
			return
		}
		sc.handler.handle(bs, sc, *conn)
	}
}
