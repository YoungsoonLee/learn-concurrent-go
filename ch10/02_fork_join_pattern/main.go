package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CodeDepth represents the depth of the code
type CodeDepth struct {
	file  string
	level int
}

func deepestNestedBlock(filename string) CodeDepth {
	code, _ := os.ReadFile(filename)
	max := 0
	level := 0
	for _, c := range code {
		if c == '{' {
			level++
			max = int(math.Max(float64(max), float64(level)))
		} else if c == '}' {
			level--
		}
	}
	return CodeDepth{filename, max}
}

func forkIfNeeded(path string, info os.FileInfo, wg *sync.WaitGroup, result chan CodeDepth) {
	if !info.IsDir() && strings.HasSuffix(path, ".go") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result <- deepestNestedBlock(path)
		}()
	}
}

func joinResults(partialResults chan CodeDepth) chan CodeDepth {
	finalResult := make(chan CodeDepth)
	max := CodeDepth{"", 0}
	go func() {
		for result := range partialResults {
			if result.level > max.level {
				max = result
			}
		}
		finalResult <- max
	}()
	return finalResult
}

func main() {
	dir := os.Args[1]
	partialResults := make(chan CodeDepth)
	wg := sync.WaitGroup{}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { // Walks the root directory, and for every file, calls the fork function, creating goroutines
		forkIfNeeded(path, info, &wg, partialResults)
		return nil
	})

	finalResult := joinResults(partialResults) // Calls the join function and gets the channel that will contain the final result

	wg.Wait()
	close(partialResults)

	result := <-finalResult

	fmt.Printf("%s has the deepest nested code block of %d\n", result.file, result.level)
}
