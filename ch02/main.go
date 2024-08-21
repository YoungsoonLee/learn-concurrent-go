package main

import (
	"fmt"
	"runtime"
	"time"
)

func doWork(id int) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func sayHello() {
	fmt.Println("Hello, World!")
}

func main() {
	// for i := 0; i < 5; i++ {
	// 	go doWork(i)
	// }
	// time.Sleep(2 * time.Second)

	// fmt.Println("Number of CPUs:", runtime.NumCPU())
	// fmt.Println("Number of GOMAXPROCS:", runtime.GOMAXPROCS(0))

	go sayHello()
	runtime.Gosched()
	fmt.Println("Finished")
}
