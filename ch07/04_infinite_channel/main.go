package main

import (
	"fmt"
	"time"
)

func receiver(message <-chan int) { // Receive-only channel
	for {
		msg, more := <-message // Check if the channel is closed
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, more)
		time.Sleep(1 * time.Second)
		if !more { // check if the channel is closed
			return
		}
	}
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}
