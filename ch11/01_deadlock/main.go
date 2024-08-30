package main

import (
	"fmt"
	"sync"
	"time"
)

func red(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Red: Lock1 and lock2 acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: Lock1 and lock2 released")
	}
}

func blue(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Blue: Lock1 and lock2 acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Lock1 and lock2 released")
	}
}

func main() {
	loackA := new(sync.Mutex)
	loackB := new(sync.Mutex)
	go red(loackA, loackB)
	go blue(loackA, loackB)
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
