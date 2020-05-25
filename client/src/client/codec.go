package client

import (
	"bytes"
	"fmt"
	"net"
)

var STATUS = []byte{0, 0}

type ICodec interface {
	encode(bs []byte, conn net.Conn)
	decode(conn net.Conn) []byte
}

type Codecer struct {
	len     int
	bodyLen int
	status  int
}

func NewCodecer() ICodec {
	return Codecer{len: 6, bodyLen: 4, status: 2}
}

func (codecer Codecer) encode(bs []byte, conn net.Conn) {
	conn.Write(makeBSLen(bs, codecer))
	conn.Write(STATUS)
	conn.Write(bs)
}

func (codecer Codecer) decode(conn net.Conn) []byte {
	b := new(bytes.Buffer)
	len := 1024
	data := make([]byte, len)
	n := 0
	_len := 0
	bodyLen := codecer.len
	flag := true
	var err error
	for {
		n, err = conn.Read(data)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		b.Write(data[0:n])
		_len += n
		if flag && _len >= codecer.len {
			bodyLen += makeLen(b.Next(codecer.len), codecer.bodyLen)
			flag = false
		}
		if bodyLen == _len {
			break
		}
	}
	return b.Bytes()
}

func makeLen(bs []byte, len int) int {
	var l int = 0
	for i := 0; i < len; i++ {
		l = (l << 8) | int(bs[i])
	}
	return l
}

func makeBSLen(bs []byte, codecer Codecer) []byte {
	return makeLLen(len(bs), codecer)
}

func makeLLen(l int, codecer Codecer) []byte {
	var lbs []byte
	lbs = make([]byte, codecer.bodyLen)
	for i := codecer.bodyLen - 1; i > -1; i-- {
		lbs[i] = byte(l)
		l = l >> 8
	}
	return lbs
}
