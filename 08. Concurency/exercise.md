## Practical Exercises

### Exercise 1: Simple Goroutine Counter

Create a program that demonstrates basic goroutine usage with a counter:

```go
// goroutine_counter.go
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
```

### Exercise 2: Channel-Based Number Processor

Build a pipeline using goroutines and channels to process numbers:

```go
// number_pipeline.go
package main

import (
	"fmt"
)

// generator creates a channel and sends numbers 1 to max on it
func generator(max int) <-chan int {
	out := make(chan int)
	
	go func() {
		for i := 1; i <= max; i++ {
			out <- i
		}
		close(out)
	}()
	
	return out
}

// square receives numbers from a channel, squares them, and sends results to a new channel
func square(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()
	
	return out
}

// filter receives numbers from a channel, filters out odd numbers, and sends results to a new channel
func filter(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		for num := range in {
			if num%2 == 0 { // Only keep even numbers
				out <- num
			}
		}
		close(out)
	}()
	
	return out
}

// sum adds up all numbers from the input channel and sends the total on the output channel
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		total := 0
		for num := range in {
			total += num
		}
		out <- total
		close(out)
	}()
	
	return out
}

func main() {
	// Create a pipeline
	c1 := generator(10)          // Generate: 1, 2, 3, ..., 10
	c2 := square(c1)             // Square: 1, 4, 9, ..., 100
	c3 := filter(c2)             // Filter: 4, 16, 36, 64, 100
	c4 := sum(c3)                // Sum: 220
	
	// Get the result from the end of the pipeline
	fmt.Println("Sum of squares of even numbers:", <-c4)
}
```

### Exercise 3: Worker Pool

Implement a worker pool to distribute tasks among multiple goroutines:

```go
// worker_pool.go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID      int
	Content string
}

// Worker represents a worker that processes tasks
func worker(id int, tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for task := range tasks {
		// Simulate processing time
		processingTime := time.Duration(task.ID%3+1) * 100 * time.Millisecond
		time.Sleep(processingTime)
		
		// Process the task
		result := fmt.Sprintf("Worker %d processed task %d (%s) in %v",
			id, task.ID, task.Content, processingTime)
		
		// Send the result
		results <- result
	}
}

func main() {
	// Create channels for tasks and results
	tasksChan := make(chan Task, 10)
	resultsChan := make(chan string, 10)
	
	// Create a WaitGroup to manage workers
	var wg sync.WaitGroup
	
	// Launch 3 workers
	numWorkers := 3
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, tasksChan, resultsChan, &wg)
	}
	
	// Send 10 tasks
	go func() {
		for i := 1; i <= 10; i++ {
			task := Task{
				ID:      i,
				Content: fmt.Sprintf("Task content %d", i),
			}
			tasksChan <- task
		}
		close(tasksChan) // Close tasks channel when all tasks are sent
	}()
	
	// Launch a goroutine to close the results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Collect and display results
	for result := range resultsChan {
		fmt.Println(result)
	}
	
	fmt.Println("All tasks have been processed!")
}
```
