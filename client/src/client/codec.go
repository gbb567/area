package client

import (
	"bytes"
	"fmt"
	"time"
)

var handleLen = 4
var noLen = 8
var headLen = 12
var empty = []byte("")

func encode(i int, ins string, bs []byte) (int64, []byte) {
	b := new(bytes.Buffer)
	k := time.Now().Unix()
	b.Write(i2bs(i))
	b.Write(i642bs(k))
	if ins != "" {
		b.Write([]byte(ins))
		if bs != nil {
			b.Write(bs)
		}
	}
	return k, b.Bytes()
}

func decode(len int, bs []byte) (int64, []byte) {
	if len < noLen {
		fmt.Println("dont decode")
		return unKnow, empty
	}
	return bs2i64(bs), bs[noLen:len]
}

func decodeResp(bs []byte) (int, []byte) {
	return bs2i(bs), bs[insLen:]
}

func i2bs(v int) []byte {
	b := make([]byte, handleLen)
	for i := 0; i < handleLen; i++ {
		b[i] = byte(v >> (i * 8))
	}
	return b
}

func i642bs(v int64) []byte {
	b := make([]byte, noLen)
	for i := 0; i < noLen; i++ {
		b[i] = byte(v >> (i * 8))
	}
	return b
}

func bs2i(bs []byte) int {
	var r = 0
	for i := 0; i < handleLen; i++ {
		r |= int(bs[i]) << (i * 8)
	}
	return r
}

func bs2i64(bs []byte) int64 {
	var r int64 = 0
	for i := 0; i < noLen; i++ {
		r |= int64(bs[i]) << (i * 8)
	}
	return r
}
