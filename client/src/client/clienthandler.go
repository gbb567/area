package client

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type ClientHandler struct {
	udpConn *net.UDPConn
}

type RH interface {
	r(*Client, []byte)
}
type _SLinkR struct{}
type _CLinkR struct{}
type _CLink2R struct{}
type _ClientT struct{}

var rm map[int]RH

func init() {
	rm = make(map[int]RH)
	rm[SLinkIns] = _SLinkR{}
	rm[CLinkIns] = _CLinkR{}
	rm[CLink2Ins] = _CLink2R{}
	rm[Client_T] = _ClientT{}
}
func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (c *Client) DoResp(bs []byte) {
	if len(bs) < insLen {
		return
	}
	ins, b := decodeResp(bs)
	rh, ok := rm[ins]
	if ok {
		rh.r(c, b)
	} else {
		fmt.Println(string(bs))
	}
}

func (ch ClientHandler) Write(ins int, bs []byte, addr *net.UDPAddr) {
	b := new(bytes.Buffer)
	b.Write(i642bs(time.Now().Unix()))
	b.Write(i2bs(ins))
	b.Write(bs)
	ch.udpConn.WriteToUDP(b.Bytes(), addr)
}

func (r _SLinkR) r(c *Client, bs []byte) {
	addr, err := net.ResolveUDPAddr("udp4", string(bs))
	if err != nil {
		return
	}
	c.ch.Write(CLinkIns, c.addrBs, addr)
}

func (r _CLinkR) r(c *Client, bs []byte) {
	addr, err := net.ResolveUDPAddr("udp4", string(bs))
	if err != nil {
		return
	}
	fmt.Println("hello")
	c.ch.Write(CLink2Ins, c.addrBs, addr)
}

func (r _CLink2R) r(c *Client, bs []byte) {
	addr, err := net.ResolveUDPAddr("udp4", string(bs))
	if err != nil {
		return
	}
	fmt.Println("hello")
	c.ch.Write(Client_T, hello, addr)
}

func (r _ClientT) r(c *Client, bs []byte) {
	fmt.Println(string(bs))
	//c.ch.Write(CLinkIns, hello, addr)
}
