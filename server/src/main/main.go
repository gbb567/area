package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("program over")
		if err := recover(); err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}
	}()
	(*NewServer()).Bind()
}
