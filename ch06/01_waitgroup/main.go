package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, &wg)
	}
	wg.Wait()
	fmt.Println("All done!")
}
