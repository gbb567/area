package client

import (
	"fmt"
	"net"
	"strings"
)

var hello = []byte("hello")
var _list = []byte("list")
var _nat = []byte("_nat")

func (c Client) List() {
	c.Write(_list)
	fmt.Println(string(c.Read()))
}

func (c Client) Link() {
	c.Write([]byte("link123:32:123:21"))
}

func (c Client) Nat() {
	c.Write(_nat)
	v := string(c.Read())
	vs := strings.Split(v, "\\")
	addr, err1 := net.ResolveUDPAddr(vs[0], vs[1])
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	udpAddr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:10010")
	con, err2 := net.ListenUDP("udp", udpAddr)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	con.WriteToUDP(hello, addr)
	buf := make([]byte, 1)
	con.ReadFromUDP(buf)
	fmt.Print("nat")
	fmt.Println(string(buf))
}
