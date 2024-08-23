package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const AllLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	m := sync.Mutex{}
	var frequency = make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, &m, frequency)
	}

	time.Sleep(60 * time.Second)
	m.Lock()
	for i, c := range AllLetters {
		fmt.Printf("%c: %d\n", c, frequency[i])
	}
	m.Unlock()
}

func countLetters(url string, m *sync.Mutex, frequency []int) {
	//m.Lock()
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to get %s", url))
	}

	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, c)
		if cIndex > 0 {
			m.Lock()
			frequency[cIndex]++
			m.Unlock()
		}
	}

	fmt.Println("Completed: ", url, time.Now().Format("2006-01-02 15:04:05"))
	//m.Unlock()
}
