package main

import "fmt"

func primeMultipleFilter(numbers <-chan int, quit chan<- int) {
	var right chan int
	p := <-numbers
	fmt.Println(p)

	for n := range numbers { // Reads the next number from the input channel
		if n%p != 0 {
			if right == nil {
				right = make(chan int)
				go primeMultipleFilter(right, quit) // If the current goroutine has no right, it starts a new goroutine and connects to it with a channel.
			}
			right <- n
		}
	}
	if right == nil {
		close(quit) // If the current goroutine has no right, it closes the quit channel.
	} else {
		close(right) // If the current goroutine has a right, it closes the right channel.
	}

}

func main() {
	numbers := make(chan int)
	quit := make(chan int)
	go primeMultipleFilter(numbers, quit)

	for i := 2; i < 100000; i++ {
		numbers <- i
	}

	close(numbers)
	<-quit // Waits for the quit channel to close
}
