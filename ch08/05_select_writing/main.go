package main

import (
	"fmt"
	"math"

	"golang.org/x/exp/rand"
)

func primesOnly(inputs <-chan int) <-chan int {
	results := make(chan int)
	go func() {
		for c := range inputs {
			isPrime := c != 1
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				if c%i == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				results <- c
			}
		}
	}()
	return results
}

func main() {
	numbersChannel := make(chan int)
	primes := primesOnly(numbersChannel)

	for i := 0; i < 100; i++ {
		select {
		case numbersChannel <- rand.Intn(100000000) + 1:
		case prime := <-primes:
			fmt.Println("Prime number found:", prime)
			i++
		}
	}
}
