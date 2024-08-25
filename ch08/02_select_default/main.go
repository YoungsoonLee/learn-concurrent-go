package main

import (
	"fmt"
	"time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
	msg := make(chan string)
	go func() {
		time.Sleep(seconds)
		msg <- "Hello"
	}()
	return msg
}

func main() {
	messages := sendMsgAfter(2 * time.Second)
	for {
		select {
		case msg := <-messages:
			fmt.Println("Message received:", msg)
			return
		default:
			fmt.Println("No messages waiting")
			time.Sleep(1 * time.Second)
		}
	}
}
