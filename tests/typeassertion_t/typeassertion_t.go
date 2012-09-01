package main

import (
	"fmt"
)

func main() {
	var i interface{}
	
	i=int32(14)
	
	i=make(chan int)
	
	switch i.(type) {
		case chan int: fmt.Println("It is a chan int")
		case int: fmt.Println("It is an int")
		case chan string: fmt.Println("It is a chan string")
		default: fmt.Println("Unknown!")
	}
}
