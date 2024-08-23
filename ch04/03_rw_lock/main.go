package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func matchRecorder(matchEvents *[]string, m *sync.RWMutex) {
	for i := 0; ; i++ {
		m.Lock()
		*matchEvents = append(*matchEvents, "Match event "+strconv.Itoa(i)) // write
		m.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, m *sync.RWMutex, st time.Time) {
	for i := 0; i < 100; i++ {
		m.RLock()
		allEvents := copyAllEvents(mEvents) // read
		m.RUnlock()

		timeTaken := time.Since(st)
		fmt.Println(len(allEvents), "events in", timeTaken)
	}
}

func copyAllEvents(events *[]string) []string {
	allEvents := make([]string, len(*events))
	copy(allEvents, *events)
	return allEvents
}

func main() {
	m := sync.RWMutex{}
	var matchEvents = make([]string, 0, 10000)
	for j := 0; j < 10000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}

	go matchRecorder(&matchEvents, &m)

	st := time.Now()
	for j := 0; j < 5000; j++ {
		go clientHandler(&matchEvents, &m, st)
	}
	time.Sleep(100 * time.Second)
}
