package main

import (
	"fmt"
	"time"
)

const (
	ovenTime           = 5
	everyThingElseTime = 2
)

func prepareTray(trayNumber int) string {
	fmt.Println("Preparing empty tray", trayNumber)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("tray number %d", trayNumber)
}

func mixture(tray string) string {
	fmt.Println("Pouring cupcake Mixture in", tray)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("cupcake in %s", tray)
}

func bake(mixture string) string {
	fmt.Println("Baking", mixture)
	time.Sleep(ovenTime * time.Second)
	return fmt.Sprintf("baked %s", mixture)
}

func addToppings(baked string) string {
	fmt.Println("Adding toppings on", baked)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("topped %s", baked)
}

func box(finishedCupCake string) string {
	fmt.Println("Boxing", finishedCupCake)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("boxed %s", finishedCupCake)
}

func sequential() {
	//startTime := time.Now()
	for i := 0; i < 10; i++ {
		result := box(addToppings(bake(mixture(prepareTray(i)))))
		fmt.Println("Accepting:", result)
	}
	//fmt.Println("Sequential execution took:", time.Since(startTime))
}

func addOnPipe[X, Y any](q <-chan int, f func(X) Y, in <-chan X) <-chan Y {
	out := make(chan Y) // Creates an output channel of type Y
	go func() {         // Starts the goroutine
		defer close(out)
		for {
			select {
			case <-q: // When the quit channel is closed,exits the loop and terminates the goroutine
				return
			case v := <-in: // Receives a message on the input channel if one is available
				out <- f(v) // Calls the function f and outputs the function’s return value on the output channel

			}
		}
	}()
	return out
}

func pipeline() {
	input := make(chan int) // Creates an input channel of type int
	quit := make(chan int)  // Creates a quit channel

	output := addOnPipe(quit, box,
		addOnPipe(quit, addToppings,
			addOnPipe(quit, bake,
				addOnPipe(quit, mixture,
					addOnPipe(quit, prepareTray, input)))))

	go func() {
		for i := 0; i < 10; i++ {
			input <- i
		}
	}()

	for i := 0; i < 10; i++ {
		fmt.Println("Accepting:", <-output)
	}

}

func main() {
	//sequential()
	pipeline()
}
