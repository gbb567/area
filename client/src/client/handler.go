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
	con, err2 := net.DialUDP("udp", nil, addr)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	con.Write(hello)
	buf := make([]byte, 1)
	fmt.Println(111)
	con.Read(buf)
	fmt.Print("nat")
	fmt.Println(string(buf))
}
