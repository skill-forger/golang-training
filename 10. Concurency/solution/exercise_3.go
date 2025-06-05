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
