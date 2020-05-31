package client

import (
	"fmt"
	"net"
	"time"
)

var (
	beat      = 0
	catalog   = 1
	nat       = 2
	SLinkIns  = 3
	CLinkIns  = 4
	CLink2Ins = 5
	Client_T  = 120
)
var insLen = 4
var dur = 5 * time.Second
var hello = []byte("hello")

func (c *Client) beat() {
	tiker := time.NewTicker(time.Second)
	for {
		select {
		case <-tiker.C:
			if !c.IsLive() {
				fmt.Println("beat exit")
				return
			}
			c.Check(2 * dur)
			c.Write(beat, "", nil)
		}
	}
}

func (c *Client) List() string {
	c.WriteForSync(catalog, "list", nil)
	return string(c.Read())
}

func (c *Client) add() bool {
	c.write0(catalog, "add", nil, false)
	_, bs := c.read0(false)
	return string(bs) == "ok"
}

func (c *Client) Remove() {
	c.Write(catalog, "remove", nil)
}

func (c *Client) Link(ipAndPort string) {
	_link(c, ipAndPort, true)
}

func (c *Client) Nat() bool {
	fmt.Println("nat3")
	c.WriteForSync(nat, "check", nil)
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

func _link(c *Client, ipAndPort string, flag bool) {
	addr, err1 := net.ResolveUDPAddr("udp4", ipAndPort)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//todo dist nat check
	c.ch.Write(Client_T, hello, addr)
	if flag {
		c.Write(SLinkIns, "link", []byte(ipAndPort))
	}
}

func (c *Client) Wait() {
	for {
		if !c.IsLive() {
			fmt.Println("client exit")
			break
		}
		no, bs := c.read0(true)
		if len(bs) == 0 && no == 0 {
			c.UpdateTime()
		} else {
			c.DoResp(bs)
		}
	}
}
