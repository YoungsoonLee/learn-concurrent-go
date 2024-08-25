package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// Assigning a nil value to the channel variable after the receiver detects
// that the channel has been closed has the effect of disabling that case statement.
// This allows the receiving goroutine to read from the remaining open channels.
func generateAmounts(n int) <-chan int {
	amounts := make(chan int)
	go func() {
		defer close(amounts)
		for i := 0; i < n; i++ {
			amounts <- rand.Intn(1000) + 1
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return amounts
}

func main() {
	sales := generateAmounts(50)
	expenses := generateAmounts(40)

	endOfDayAmount := 0
	for sales != nil || expenses != nil {
		select {
		case sale, moreData := <-sales:
			if moreData {
				fmt.Println("Sale of: $", sale)
				endOfDayAmount += sale
			} else {
				sales = nil // If the channel has been closed, marks the channel as nil, disabling this select case !!!
			}
		case expense, moreData := <-expenses:
			if moreData {
				fmt.Println("Expense of: $", expense)
				endOfDayAmount -= expense
			} else {
				expenses = nil // If the channel has been closed, marks the channel as nil, disabling this select case !!!
			}
		}
	}
	fmt.Println("End of day amount:", endOfDayAmount)
}
