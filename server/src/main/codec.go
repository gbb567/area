package main

import (
	"bytes"
	"fmt"
)

var handleLen = 4
var noLen = 8
var headLen = 12
var empty = []byte("")

func resp(ins int, bs []byte) []byte {
	b := new(bytes.Buffer)
	b.Write(i2bs(ins))
	b.Write(bs)
	return b.Bytes()
}

func encode(no int64, bs []byte) []byte {
	b := new(bytes.Buffer)
	b.Write(i642bs(no))
	if bs != nil {
		b.Write(bs)
	}
	return b.Bytes()
}

func decode(len int, bs []byte) (int, int64, []byte) {
	if len < headLen {
		fmt.Println("dont decode")
		return 0, 0, empty
	}
	return bs2i(bs), bs2i64(bs[handleLen : handleLen+noLen]), bs[headLen:len]
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
