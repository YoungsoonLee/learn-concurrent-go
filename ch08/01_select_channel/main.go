package main

import "time"

func writeEvery(msg string, seconds time.Duration) <-chan string {
	c := make(chan string)
	go func() {
		for {
			time.Sleep(seconds)
			c <- msg
		}
	}()
	return c
}

func main() {
	messageFromA := writeEvery("A", 1*time.Second)
	messageFromB := writeEvery("B", 3*time.Second)

	for {
		select {
		case msg := <-messageFromA:
			println(msg)
		case msg := <-messageFromB:
			println(msg)
		}
	}
}
