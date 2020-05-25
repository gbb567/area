package main

import (
	"fmt"
	"net"
)

type IHandler interface {
	handle(bs []byte, sc ServerConfig, conn net.Conn)
}

type Handler struct {
}

type MapHandler struct {
	Handler
	handlerMap map[string]IHandler
}

func _newHandlerMap() map[string]IHandler {
	m := make(map[string]IHandler)
	m["list"] = CatalogHandler{}
	m["link"] = LinkHandler{}
	m["_nat"] = NatCheckHandler{}
	return m
}

func NewHandler() MapHandler {
	return MapHandler{handlerMap: _newHandlerMap()}
}

func (handler MapHandler) handle(bs []byte, sc ServerConfig, conn net.Conn) {
	key := string(bs[0:4])
	if handler.handlerMap[key] != nil {
		handler.handlerMap[key].handle(bs[4:], sc, conn)
	} else {
		fmt.Println(key)
	}
}
