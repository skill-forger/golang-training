package main

import (
	"fmt"
	"sync"
	"time"
)

func countUp(name string, count int, wg *sync.WaitGroup) {
	// Ensure we mark this goroutine as done when the function completes
	defer wg.Done()

	for i := 1; i <= count; i++ {
		fmt.Printf("%s: %d\n", name, i)
		// Simulate work with a small delay
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("%s completed counting to %d\n", name, count)
}

func main() {
	// Create a WaitGroup to manage our goroutines
	var wg sync.WaitGroup

	fmt.Println("Starting counters...")

	// Launch three goroutines with different counts
	wg.Add(3)
	go countUp("Counter A", 5, &wg)
	go countUp("Counter B", 3, &wg)
	go countUp("Counter C", 7, &wg)

	// Wait for all goroutines to complete
	wg.Wait()

	fmt.Println("All counters have finished!")
}
