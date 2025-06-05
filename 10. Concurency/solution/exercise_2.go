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
	c1 := generator(10) // Generate: 1, 2, 3, ..., 10
	c2 := square(c1)    // Square: 1, 4, 9, ..., 100
	c3 := filter(c2)    // Filter: 4, 16, 36, 64, 100
	c4 := sum(c3)       // Sum: 220

	// Get the result from the end of the pipeline
	fmt.Println("Sum of squares of even numbers:", <-c4)
}
