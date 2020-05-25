package client

import (
	"fmt"
	"net"
)

type IClient interface {
	Read() []byte
	Write([]byte)
}

type Client struct {
	codec ICodec
	conn  net.Conn
}

func NewClient() Client {
	c := Client{codec: NewCodecer()}
	connect(&c)
	return c
}

func (c Client) Write(bs []byte) {
	c.codec.encode(bs, c.conn)
}

func (client Client) Read() []byte {
	return client.codec.decode(client.conn)
}

func connect(client *Client) {
	conn, err := net.Dial("tcp", "127.0.0.1:4040")
	if err != nil {
		fmt.Print("err")
		panic(err)
	}
	client.conn = conn
}
