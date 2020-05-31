package client

import (
	"bytes"
	"container/list"
	"net"
	"sync"
	"time"
)

var ok = []byte("ok")
var fail = []byte("fail")

var spilter = []byte(";")
var appendLock sync.Mutex
var removeLock sync.Mutex

type UDPAddrWrap struct {
	addr *net.UDPAddr
	time time.Time
}

type CatalogHandler struct{}
type Catalog struct {
	list *list.List
	dict map[string]*list.Element
}

func NewCatalog() *Catalog {
	return &Catalog{list: list.New(), dict: make(map[string]*list.Element)}
}

func listToByte(ls *list.List) []byte {
	if ls.Len() == 0 {
		return []byte{}
	}
	bs := bytes.NewBuffer([]byte{})
	var addr *net.UDPAddr
	header := ls.Front()
	if header != nil {
		addr = header.Value.(*UDPAddrWrap).addr
		bs.Write([]byte(addr.String()))
	}
	for header = header.Next(); header != nil; header = header.Next() {
		addr = header.Value.(*UDPAddrWrap).addr
		bs.Write(spilter)
		bs.Write([]byte(addr.String()))
	}
	return bs.Bytes()
}

func (cl *Catalog) List() *list.List {
	return cl.list
}

func (cl *Catalog) Update(addr *net.UDPAddr) {
	key := getKey(addr)
	ele := cl.dict[key]
	if ele != nil {
		ele.Value.(*UDPAddrWrap).time = time.Now()
	}
}

func (cl *Catalog) Remove(addr *net.UDPAddr) {
	removeLock.Lock()
	defer removeLock.Unlock()
	key := getKey(addr)
	if cl.dict[key] != nil {
		(*cl.list).Remove(cl.dict[key])
		delete(cl.dict, key)
	}
}

func (cl *Catalog) append(addr *net.UDPAddr) bool {
	appendLock.Lock()
	defer appendLock.Unlock()
	key := getKey(addr)
	if cl.dict[key] != nil {
		return false
	}
	cl.dict[key] = (*cl.list).PushBack(wrapAddr(addr))
	return true
}

func (cl *Catalog) Get(key string) *net.UDPAddr {
	ele := cl.dict[key]
	if ele == nil {
		return nil
	}
	return ele.Value.(*UDPAddrWrap).addr
}

func getKey(addr *net.UDPAddr) string {
	return (*addr).String()
}

func wrapAddr(addr *net.UDPAddr) *UDPAddrWrap {
	wrap := new(UDPAddrWrap)
	wrap.addr = addr
	wrap.time = time.Now()
	return wrap
}
