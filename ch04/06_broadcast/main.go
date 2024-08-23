package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	playersInGame := 4

	for playerID := 0; playerID < 4; playerID++ {
		go playerHandler(cond, &playersInGame, playerID)
		time.Sleep(1 * time.Second)
	}
}

func playerHandler(cond *sync.Cond, playersInGame *int, playerID int) {
	cond.L.Lock()

	fmt.Println("Player", playerID, "joined the game")
	*playersInGame--
	if *playersInGame == 0 {
		cond.Broadcast()
	}

	for *playersInGame > 0 {
		fmt.Println("Player", playerID, "waiting for others to join")
		cond.Wait()
	}

	cond.L.Unlock()

	fmt.Println("All players connected. Ready player", playerID)

	// Game starts here

}
