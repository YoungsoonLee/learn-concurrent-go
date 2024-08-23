package main

import "fmt"

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	fmt.Println("Sending hello...")
	msgChannel <- "hello"
	fmt.Println("Sending world...")
	msgChannel <- "world"
	fmt.Println("Sending goodbye...")
	msgChannel <- "goodbye"
	close(msgChannel)
}

func receiver(msgChannel chan string) {
	for {
		msg := <-msgChannel
		fmt.Println("Received:", msg)
	}
}
