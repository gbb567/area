package main

import (
	"client"
	"fmt"
	"strings"
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
	if c.Add() {
		fmt.Println("add")
		c.Nat()
		str := c.List()
		fmt.Println(str)
		strs := strings.Split(str, ";")
		if len(strs) > 1 {
			fmt.Println("link")
			c.Link(strs[0])
			c.Wait()
		} else {
			c.Wait()
		}
	} else {
		fmt.Println("no add")
	}
}
