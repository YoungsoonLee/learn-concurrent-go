package main

import (
	"fmt"
	"sync"
)

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work done")

	cond.L.Lock()
	cond.Signal() // !!!
	cond.L.Unlock()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()

	for i := 0; i < 50000; i++ {
		go doWork(cond)
		fmt.Println("Waiting for work to be done")
		cond.Wait()
		fmt.Println("Work completed")
	}

	cond.L.Unlock()
}
