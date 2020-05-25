package main

import (
	"client"
	"fmt"
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
	c.List()
	c.Nat()
	time.Sleep(10 * time.Second)
}
