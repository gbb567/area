package client

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var empty = []byte("")

type IClient interface {
	Read() []byte
	Write(i int, ins string, bs []byte)
}

type Client struct {
	serverAddr *net.UDPAddr
	udpConn    *net.UDPConn
}

func NewClient() Client {
	sAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:9998")
	if err != nil {
		panic(err)
	}
	c := Client{serverAddr: sAddr}
	bind(&c)
	return c
}

func (c Client) Write(i int, ins string, bs []byte) {
	b := new(bytes.Buffer)
	b.Write(i2bs(i))
	b.Write([]byte(ins))
	if bs != nil {
		b.Write(bs)
	}
	c.udpConn.WriteToUDP(b.Bytes(), c.serverAddr)
}

func (c Client) Read() []byte {
	t := time.Now()
	c.udpConn.SetDeadline(t.Add(time.Duration(5 * time.Second)))
	b := make([]byte, 1024)
	len, _, err := c.udpConn.ReadFromUDP(b)
	if err != nil {
		return empty
	}
	return b[0:len]
}

func bind(client *Client) {
	rand.Seed(time.Now().UnixNano())
	udpAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", rand.Intn(100)+10010))
	con, err2 := net.ListenUDP("udp", udpAddr)
	if err2 != nil {
		fmt.Println(err2)
	}
	client.udpConn = con
}

func i2bs(i int) []byte {
	b := make([]byte, 4)
	b[3] = byte(i)
	b[2] = byte(i >> 8)
	b[1] = byte(i >> 16)
	b[0] = byte(i >> 24)
	return b
}
