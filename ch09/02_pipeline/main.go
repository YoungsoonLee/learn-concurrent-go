package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

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

func main() {
	quit := make(chan int)
	defer close(quit)

	results := extractWords(quit, downloadPages(quit, generateURLs(quit))) /// Pipeline

	for result := range results {
		fmt.Println(result)
	}

}
