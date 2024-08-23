package main

import (
	"fmt"
	"time"
)

func receiver(message <-chan int) { // Receive-only channel
	for {
		msg := <-message
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

func sender(message chan<- int) { // Send-only channel
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		message <- i
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	go sender(msgChannel)
	time.Sleep(5 * time.Second)
}
