package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("number from chan: ", <-numbers)
		}
		close(quit)
	}()
}

func main() {
	numbers := make(chan int)
	quit := make(chan int)

	printNumbers(numbers, quit)

	next := 0
	for i := 1; ; i++ {
		//next++
		next += i
		fmt.Println("Sending number: ", next)

		select {
		case numbers <- next:
		case <-quit:
			fmt.Println("Quitting")
			return
		}
	}
}
