package main

import (
	"fmt"
	"sync"
	"time"
)

func stingy(money *int, m *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		m.Lock()
		*money += 10
		m.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, m *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		m.Lock()
		*money -= 10
		m.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	m := sync.Mutex{}

	go stingy(&money, &m)
	go spendy(&money, &m)

	time.Sleep(2 * time.Second)

	m.Lock()
	fmt.Println("Money left:", money)
	m.Unlock()
}
