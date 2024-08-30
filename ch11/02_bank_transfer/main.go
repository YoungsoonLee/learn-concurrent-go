package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

// Arbitrator is a struct that arbitrates access to accounts
type Arbitrator struct {
	accountsInUse map[string]bool
	cond          *sync.Cond
}

// NewArbitrator creates a new Arbitrator
func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: make(map[string]bool),
		cond:          sync.NewCond(new(sync.Mutex)),
	}
}

// Lock locks the account
func (a *Arbitrator) Lock(ids ...string) {
	a.cond.L.Lock()
	for allAvailable := false; !allAvailable; {
		allAvailable = true
		for _, id := range ids {
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
			}
		}
	}
	for _, id := range ids {
		a.accountsInUse[id] = true
	}
	a.cond.L.Unlock()
}

// Unlock unlocks the account
func (a *Arbitrator) Unlock(ids ...string) {
	a.cond.L.Lock()
	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	a.cond.Broadcast()
	a.cond.L.Unlock()
}

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

// Transfer transfers money from one account to another
func (ba *BankAccount) Transfer(to *BankAccount, amount int, tellerID int, arb *Arbitrator) {
	fmt.Printf("%d Locking %s and %s\n", tellerID, ba.id, to.id)
	arb.Lock(ba.id, to.id)
	ba.balance -= amount
	to.balance += amount
	arb.Unlock(ba.id, to.id)
	fmt.Printf("%d Unlocked %s and %s\n", tellerID, ba.id, to.id)
}

func main() {
	accounts := []BankAccount{
		*NewBankAccount("Sam"),
		*NewBankAccount("Paul"),
		*NewBankAccount("Amy"),
		*NewBankAccount("John"),
	}

	total := len(accounts)
	arb := NewArbitrator()

	for i := 0; i < total; i++ {
		go func(tellerID int) {
			for i := 1; i < 1000; i++ {
				from, to := rand.Intn(total), rand.Intn(total)

				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(&accounts[to], 10, tellerID, arb)
			}
			fmt.Println("Teller", tellerID, "done")
		}(i)
	}
	time.Sleep(60 * time.Second)
}
