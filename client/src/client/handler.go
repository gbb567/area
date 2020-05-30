package client

import (
	"fmt"
	"net"
)

var (
	beat    = 0
	catalog = 1
	nat     = 2
	link    = 3
)

var hello = []byte("hello")

func (c Client) beat() {
	c.Write(beat, "", nil)
}

func (c Client) List() string {
	c.Write(catalog, "list", nil)
	return string(c.Read())
}

func (c Client) Add() bool {
	c.Write(catalog, "add", nil)
	return string(c.Read()) == "ok"
}

func (c Client) Remove() {
	c.Write(catalog, "remove", nil)
}

func (c Client) Link(ipAndPort string) {
	_link(c, ipAndPort, true)
}

func (c Client) Nat() bool {
	c.Write(nat, "check", nil)
	fmt.Println("nat3")
	//check Restricted Cone
	if string(c.Read()) == "2" {
		fmt.Println("nat2")
		//check Full Cone,need two public ips
		// if printNat(con, addr) == "1"{
		// 	fmt.Println("nat1")
		// }
	}
	return true
}

func _link(c Client, ipAndPort string, flag bool) {
	addr, err1 := net.ResolveUDPAddr("udp4", ipAndPort)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//todo dist nat check
	c.udpConn.WriteToUDP(hello, addr)
	if flag {
		c.Write(link, "link", []byte(ipAndPort))
	}
}

func (c Client) Wait() {
	for {
		bs := c.Read()
		key := string(bs)
		if key != "" {
			if key == "beat" {
				c.beat()
			} else if key == "hello" {
				fmt.Println(key)
			} else {
				// link
				_link(c, key, false)
			}
		}
	}
}
