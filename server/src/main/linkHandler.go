package main

import (
	"fmt"
	"net"
)

type LinkHandler struct {
	Handler
}

func (handler LinkHandler) handle(bs []byte, sc ServerConfig, conn net.Conn) {
	fmt.Println(string(bs))
}
