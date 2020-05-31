package client

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var ch = make(chan []byte)
var unKnow int64 = -99

var (
	start        = 0
	connected    = 1
	unconnected  = 2
	disconnected = 3
)

type IClient interface {
	Read() []byte
	Write(i int, ins string, bs []byte)
	WriteForSync(i int, ins string, bs []byte)
	Send(string, []byte)
}

type Client struct {
	serverAddr *net.UDPAddr
	udpConn    *net.UDPConn
	addr       *net.UDPAddr
	addrBs     []byte
	ch         *ClientHandler
	time       time.Time
	dict       *Catalog
	status     int
}

var noMap = make(map[int64]bool)

func NewClient() *Client {
	sAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:9998")
	if err != nil {
		panic(err)
	}
	c := Client{serverAddr: sAddr, ch: NewClientHandler(), dict: NewCatalog()}
	c.status = start
	bind(&c)
	if c.add() {
		c.status = connected
		go c.Wait()
		go c.beat()
	} else {
		c.status = unconnected
	}
	return &c
}

func (c Client) Write(i int, ins string, bs []byte) {
	c.write0(i, ins, bs, true)
}

func (c Client) WriteForSync(i int, ins string, bs []byte) {
	k, b := encode(i, ins, bs)
	noMap[k] = true
	c.udpConn.WriteToUDP(b, c.serverAddr)
}

func (c *Client) write0(i int, ins string, bs []byte, flag bool) {
	if flag && c.status != connected {
		fmt.Println("client status is not connected")
		return
	}
	_, b := encode(i, ins, bs)
	c.udpConn.WriteToUDP(b, c.serverAddr)
}

func (c Client) Read() []byte {
	return <-ch
}

func (c *Client) read0(flag bool) (int64, []byte) {
	if flag && c.status != connected {
		fmt.Println("client status is not connected")
		return unKnow, empty
	}
	t := time.Now()
	c.udpConn.SetDeadline(t.Add(time.Duration(5 * time.Second)))
	b := make([]byte, 1024)
	len, addr, err := c.udpConn.ReadFromUDP(b)
	if err != nil {
		return unKnow, empty
	}
	c.dict.append(addr)
	i, bs := decode(len, b)
	_, ok := noMap[i]
	if flag && ok {
		delete(noMap, i)
		ch <- bs
		return unKnow, empty
	}
	return i, bs
}

func (c *Client) UpdateTime() {
	c.time = time.Now()
}

func (c *Client) Check(dur time.Duration) {
	if time.Now().Sub(c.time) > dur {
		c.status = disconnected
	}
}

func (c *Client) IsLive() bool {
	return c.status == connected
}

func (c *Client) Send(ipAndPort string, bs []byte) {
	addr := c.dict.Get(ipAndPort)
	if addr != nil {
		c.ch.Write(Client_T, bs, addr)
	}
}

func bind(client *Client) {
	rand.Seed(time.Now().UnixNano())
	udpAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", rand.Intn(100)+10010))
	client.addr = udpAddr
	client.addrBs = []byte(udpAddr.String())
	con, err2 := net.ListenUDP("udp", udpAddr)
	if err2 != nil {
		fmt.Println(err2)
	}
	client.time = time.Now()
	client.udpConn = con
	client.ch.udpConn = con
}
