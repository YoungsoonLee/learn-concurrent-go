package main

import (
	"fmt"
	"sync"
	"time"
)

func receiver(message chan int, w *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = <-message
		fmt.Println("Received:", msg)
	}
	w.Done()
}

func main() {
	msgChannel := make(chan int, 3)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go receiver(msgChannel, &wg)

	for i := 1; i <= 6; i++ {
		size := len(msgChannel)
		fmt.Printf("%s Sending: %d. Buffer Size: %d\n", time.Now().Format("15:04:05"), i, size)
		msgChannel <- i
	}

	msgChannel <- -1
	wg.Wait()
}
