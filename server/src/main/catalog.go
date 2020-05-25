package main

import (
	"bytes"
	"container/list"
	"net"
	"sync"
)

var spilter = []byte(";")

type ICatalog interface {
	List() *list.List
	Append(conn *net.Conn) bool
	Remove(conn *net.Conn)
}

type Catalog struct {
	list       *list.List
	dict       map[string]*list.Element
	appendLock sync.Mutex
	removeLock sync.Mutex
}

type CatalogHandler struct {
	Handler
}

func (handler CatalogHandler) handle(rbs []byte, sc ServerConfig, conn net.Conn) {
	ls := (*sc.ICatalog).List()
	bs := bytes.NewBuffer([]byte{})
	var cn *net.Conn
	header := ls.Front()
	if header != nil {
		cn = header.Value.(*net.Conn)
		bs.Write([]byte((*cn).RemoteAddr().String()))
	}
	for header = header.Next(); header != nil; header = header.Next() {
		cn = header.Value.(*net.Conn)
		bs.Write(spilter)
		bs.Write([]byte((*cn).RemoteAddr().String()))
	}
	sc.codec.encode(bs.Bytes(), conn)
}

func NewCatalog() *ICatalog {
	var ic ICatalog
	ic = Catalog{list: list.New(), dict: make(map[string]*list.Element)}
	return &ic
}

func (cl Catalog) List() *list.List {
	return cl.list
}

func (cl Catalog) Append(conn *net.Conn) bool {
	cl.appendLock.Lock()
	defer cl.appendLock.Unlock()
	key := getKey(conn)
	if cl.dict[key] != nil {
		return false
	}
	cl.dict[key] = (*cl.list).PushBack(conn)
	return true
}
func (cl Catalog) Remove(conn *net.Conn) {
	cl.removeLock.Lock()
	defer cl.removeLock.Unlock()
	key := getKey(conn)
	if cl.dict[key] != nil {
		(*cl.list).Remove(cl.dict[key])
		delete(cl.dict, key)
	}
}
func getKey(conn *net.Conn) string {
	str := (*conn).RemoteAddr().String()
	//return str[:strings.Index(str, ":")]
	return str
}
