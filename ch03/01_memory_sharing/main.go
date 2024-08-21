package main

import (
	"fmt"
	"time"
)

func main() {
	count := 5
	go countdown(&count) // sharing memory

	for count > 0 {
		time.Sleep((500 & time.Millisecond))
		fmt.Println(count)
	}
}

func countdown(count *int) {
	for *count > 0 {
		time.Sleep(1 * time.Second)
		*count--
	}
}
