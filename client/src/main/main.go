package main

import (
	"client"
	"fmt"
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
}
