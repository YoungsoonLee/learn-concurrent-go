package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const downloaders = 20

func generateURLs(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			select {
			case urls <- fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i):
			case <-quit:
				return
			}
		}
	}()
	return urls
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != http.StatusOK {
						panic("Server returning error code: " + resp.Status)
					}
					defer resp.Body.Close()
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
				}
			case <-quit: // When a message arrives on the quit channel, terminates the goroutine
				return
			}
		}
	}()
	return pages
}

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-z]+`)
		moreData, page := true, ""
		for moreData {
			select {
			case page, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(page, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}

// FanIn function takes a quit channel and a variadic number of input channels and returns a single output channel.
func FanIn[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))

	output := make(chan K)

	for _, c := range allChannels {
		go func(channel <-chan K) { // starts a goroutine for every input channel
			defer wg.Done()
			for val := range channel {
				select {
				case output <- val: // forwards each received message to the shared output channel
				case <-quit: // if quit channel is closed, terminates the goroutine
					return
				}
			}
		}(c) // passes one input channel to the goroutine
	}

	go func() {
		wg.Wait() // wait for all the goroutines to finish and then closes the output channel
		close(output)
	}()

	return output
}

func main() {
	quit := make(chan int)
	defer close(quit)

	urls := generateURLs(quit)

	pages := make([]<-chan string, downloaders)
	for i := 0; i < downloaders; i++ {
		pages[i] = downloadPages(quit, urls) // fanout
	}

	results := extractWords(quit, FanIn(quit, pages...)) // fanin
	for word := range results {
		fmt.Println(word)
	}
}
