package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
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
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
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

func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}

		sort.Slice(uniqueWords, func(i, j int) bool {
			return len(uniqueWords[i]) > len(uniqueWords[j])
		})

		longWords <- strings.Join(uniqueWords[:10], ", ") // Once the input channel is closed, sends a string with the 10 longest words on the output channel
	}()
	return longWords
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

// Broadcast function takes a quit channel, an input channel, and an integer n, and returns a slice of n output channels.
func Broadcast[K any](quit chan int, input <-chan K, n int) []chan K {
	outputs := CreateAll[K](n) //Creates n output channels of type K
	go func() {
		defer CloseAll(outputs...) // closes all output channels when the input channel is closed

		var msg K
		moreData := true

		for moreData {
			select {
			case msg, moreData = <-input: // receives messages from the input channel
				if moreData {
					for _, output := range outputs {
						output <- msg // forwards the received message to all output channels
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return outputs
}

// CreateAll function takes an integer n and returns a slice of n output channels.
func CreateAll[K any](n int) []chan K {
	outputs := make([]chan K, n)
	for i := range outputs {
		outputs[i] = make(chan K)
	}
	return outputs
}

// CloseAll function takes a variadic number of channels and closes all of them.
func CloseAll[K any](channels ...chan K) {
	for _, c := range channels {
		close(c)
	}
}

func frequentWords(quit <-chan int, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)
	go func() {
		defer close(mostFrequentWords)
		freqMap := make(map[string]int) // Creates a map to store the frequency occurrence of each unique word
		freqList := make([]string, 0)   // Creates a slice to store the unique words
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData {
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					freqMap[word]++
				}
			case <-quit:
				return
			}
		}
		sort.Slice(freqList, func(i, j int) bool {
			return freqMap[freqList[i]] > freqMap[freqList[j]]
		})
		mostFrequentWords <- strings.Join(freqList[:10], ", ") // Sends the 10 most frequent words to the output channel
	}()
	return mostFrequentWords
}

// Take function takes a quit channel, an integer n, and an input channel and returns an output channel.
func Take[K any](quit chan int, n int, input <-chan K) <-chan K {
	output := make(chan K)
	go func() {
		defer close(output)
		moreData := true
		var msg K
		for n > 0 && moreData { // Continuesforwardingmessagesas long as there is more data and countdown n is greater than 0
			select {
			case msg, moreData = <-input: // Reads the next message from the input
				if moreData {
					output <- msg
					n--
				}
			case <-quit:
				return
			}
		}

		if n == 0 {
			close(quit)
		}
	}()
	return output
}

func main() {
	quitWords := make(chan int)
	quit := make(chan int)
	defer close(quit)

	urls := generateURLs(quit)

	pages := make([]<-chan string, downloaders)
	for i := 0; i < downloaders; i++ {
		pages[i] = downloadPages(quit, urls) // fanout
	}

	// pipeling
	// results := longestWords(quit, extractWords(quit, FanIn(quit, pages...))) // fanin
	// for word := range results {
	// 	fmt.Println(word)
	// }

	// broadcast
	// words := extractWords(quit, FanIn(quit, pages...))
	// wordsMulti := Broadcast(quit, words, 2)               //Creates a goroutine that will broadcast the contents of the words channel to two output channels
	// longestResults := longestWords(quit, wordsMulti[0])   // Creates the goroutine to find the longest words from the input channel
	// frequentResults := frequentWords(quit, wordsMulti[1]) // Creates the goroutine to find the most frequently used words from the input channel
	// fmt.Println("Longest Words:", <-longestResults)
	// fmt.Println("Most frequent Words:", <-frequentResults)

	// Take
	words := Take(quitWords, 10000, extractWords(quitWords, FanIn(quit, pages...)))
	wordsMulti := Broadcast(quit, words, 2)               //Creates a goroutine that will broadcast the contents of the words channel to two output channels
	longestResults := longestWords(quit, wordsMulti[0])   // Creates the goroutine to find the longest words from the input channel
	frequentResults := frequentWords(quit, wordsMulti[1]) // Creates the goroutine to find the most frequently used words from the input channel
	fmt.Println("Longest Words:", <-longestResults)
	fmt.Println("Most frequent Words:", <-frequentResults)
}
