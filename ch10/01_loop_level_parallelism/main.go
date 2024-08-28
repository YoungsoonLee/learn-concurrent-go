package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// FHash calculates the SHA-256 hash of a file
func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()

	sha := sha256.New()
	io.Copy(sha, file)

	return sha.Sum(nil)
}

// loop-carried dependency
func loopCarriedDependency() {
	// dir := os.Args[1]
	// files, _ := os.ReadDir(dir)
	// sha := sha256.New()

	// for _, file := range files {
	// 	if !file.IsDir() {
	// 		fPath := filepath.Join(dir, file.Name())
	// 		hashOnFile := FHash(fPath)
	// 		sha.Write(hashOnFile) // Concatenates the computed hash code to the directory one
	// 	}
	// }
	// fmt.Printf("%x\n", sha.Sum(nil))

	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int
	for _, file := range files {
		if !file.IsDir() {
			next = make(chan int)
			go func(prev, next chan int, filename string) {
				fpath := filepath.Join(dir, filename)
				hashOnFile := FHash(fpath)
				if prev != nil {
					<-prev // If the goroutine is not in the first iteration, waits until the previous iteration sends a signal
				}
				sha.Write(hashOnFile)
				next <- 0 // Signals to the next iteration that it’s done
			}(prev, next, file.Name())
			prev = next // Assigns the next channel to be previous; the next goroutine will wait on a signal from the current iteration
		}
	}
	<-next
	fmt.Printf("%x\n", sha.Sum(nil)) // Waits for the last iteration to be complete before outputting the result
}

// loop-level parallelism
func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)

	wg := sync.WaitGroup{}
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				fPath := filepath.Join(dir, filename)
				hash := FHash(fPath)
				fmt.Printf("%s: %x\n", filename, hash)
			}(file.Name())
		}
	}
	wg.Wait()

	loopCarriedDependency()
}
