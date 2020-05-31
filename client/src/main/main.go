package main

import (
	"client"
	"fmt"
	"strings"
	"time"
)

func main() {
	//component.ComponentListBox()
	defer func() {
		fmt.Println("program over")
		if err := recover(); err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}
	}()
	c := client.NewClient()
	if c.IsLive() {
		c.Nat()
		str := c.List()
		fmt.Println(str)
		strs := strings.Split(str, ";")
		if len(strs) > 1 {
			fmt.Println("link")
			c.Link(strs[0])
			c.Send(strs[0], []byte("阿鲁巴"))
			time.Sleep(2 * time.Second)
			c.Send(strs[0], []byte("阿鲁巴"))
		}
		time.Sleep(20 * time.Second)
	} else {
		fmt.Println("client error")
	}
}
