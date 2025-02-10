package main

import (
	"fmt"
)

func UnBuff() {
	fmt.Println("UnBuff channel")
	channel := make(chan int)
	fmt.Println(len(channel))
	go func() {
		channel <- 1
	}()
	fmt.Println(len(channel))
}

func Buff() {
	fmt.Println("Buff channel")
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	fmt.Println(len(ch))
	fmt.Println(cap(ch))
}

func main() {
	UnBuff()
	Buff()
}
