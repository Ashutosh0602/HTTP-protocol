package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int
var mu sync.Mutex

func incrementCounter(wg *sync.WaitGroup) {
	// fmt.Println("Incrementing counter in goroutine: ", counter)

	mu.Lock()
	counter++
	mu.Unlock()

	// fmt.Println("Counter incremented in goroutine: ", counter)
	wg.Done()
}

func main() {

	fmt.Println("Number of CPU cores:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS value:", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go incrementCounter(&wg)
	}

	// wg.Wait()

	// fmt.Println("Final counter value:", counter)

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Final counter value after all goroutines have finished:", counter)
}
