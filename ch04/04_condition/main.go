package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 1000000; i++ {
		cond.L.Lock()
		*money += 10
		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 1000000; i++ {
		cond.L.Lock()

		for *money < 50 {
			cond.Wait()
		}

		*money -= 50

		if *money < 0 {
			fmt.Println("Money is now negative")
			os.Exit(1)
		}

		cond.L.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	go stingy(&money, cond)
	go spendy(&money, cond)

	time.Sleep(2 * time.Second)

	m.Lock()
	fmt.Println("Money left:", money)
	m.Unlock()
}
