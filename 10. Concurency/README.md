# Module 10: Concurrency

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#concurrency-vs-parallelism">Concurrency vs Parallelism</a></li>
    <li><a href="#concurrency-in-golang">Concurrency in Golang</a></li>
    <li><a href="#mutual-exclusion">Mutual Exclusion</a></li>
    <li><a href="#goroutine-lifecycle">Goroutine Lifecycle</a></li>
    <li><a href="#popular-concurrency-patterns">Popular Concurrency Patterns</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand the concept of concurrency and parallelism
- Learn how Go utilize goroutines for concurrency
- Be familiar with the concept of mutual exclusion and how to utilize shared memory in goroutines
- Learn to control goroutines lifecycles
- Know common best practices and common use cases for goroutines
- Avoid common mistakes when working with interfaces

## Overview

Concurrency is the ability of a program to have multiple computations in progress at the same time.
These computations can be running on a single processor, where the operating system switches between them,
or on multiple processors, where they can run simultaneously.

In modern computing, with the prevalence of multi-core processors, concurrency is no longer a niche concept.
It's a fundamental tool for building responsive and scalable applications.
Whether you're building a web server that needs to handle thousands of simultaneous connections
or a data processing pipeline that needs to chew through terabytes of data, understanding concurrency is key.
Go was designed from the ground up with concurrency in mind,
making it a powerful language for today's computing landscape.

## Concurrency vs Parallelism

It's crucial to understand the distinction between concurrency and parallelism.

- **Concurrency** is about dealing with lots of things at once. It's a way of structuring your program to handle
  multiple tasks.
- **Parallelism** is about doing lots of things at once. It's the simultaneous execution of those tasks, typically on
  multiple processor cores.

Think of it this way: a chef chopping vegetables while also keeping an eye on a simmering pot is concurrent. If that
chef had a helper and they were both chopping vegetables at the same time, that would be parallel.

### Comparison between Concurrency and Parallelism

| Feature        | Concurrency                                                      | Parallelism                                                                    |
|----------------|------------------------------------------------------------------|--------------------------------------------------------------------------------|
| **Concept**    | Dealing with multiple tasks at the same time.                    | Doing multiple tasks at the same time.                                         |
| **Execution**  | Can be achieved on a single processor through context switching. | Requires multiple processors or cores for simultaneous execution.              |
| **Goal**       | To structure a program to handle multiple independent tasks.     | To speed up computations by running them simultaneously.                       |
| **Go Analogy** | Managing multiple goroutines.                                    | Multiple goroutines running on different OS threads across multiple CPU cores. |

## Concurrency in Golang

Go's approach to concurrency is one of its most celebrated features, primarily through **goroutines** and **channels**.

### Goroutines

To start a goroutine, you simply use the `go` keyword before a function call.

```go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from the goroutine!")
}

func main() {
	go sayHello()
	fmt.Println("Hello from the main function!")
	// Wait for the goroutine to finish (in a real app, use a sync mechanism)
	time.Sleep(100 * time.Millisecond)
}

```

### Characteristics of Goroutines

- **Lightweight**: They start with a small stack size (a few kilobytes) that can grow and shrink as needed.
- **Cheap to create**: The overhead of creating a goroutine is very low.
- **Multiplexed**: The Go runtime's scheduler maps goroutines onto a smaller number of OS threads, handling the context
  switching efficiently.

### Channels: The Artery of Goroutines

Goroutines communicate with each other using channels.
A channel is a typed conduit through which you can send and receive values with the channel operator, `<-`.

#### Channel Types:

- **Unbuffered Channels**: The default type. They require both the sender and receiver to be ready to communicate at the
  same time. If one is ready and the other is not, it will block until the other side is ready.
- **Buffered Channels**: These have a capacity and will only block the sender if the buffer is full. The receiver will
  block if the buffer is empty.

```go
// Unbuffered channel
ch := make(chan int)

// Buffered channel with a capacity of 10
bufferedCh := make(chan string, 10)
```

#### Send and Receive Operations:

To send and receive the data from and to channels, the `<-` operator is used:

```go
// Send a value to a channel
ch <- 10

// Receive a value from a channel
value := <-ch
```

#### Multiplexing with `select`

The `select` statement lets a goroutine wait on multiple communication operations.
A `select` blocks until one of its cases can run, then it executes that case.
It chooses one at random if multiple are ready.

```go
package main

import (
	"fmt"
)

func main() {
	select {
	case msg1 := <-ch1:
		fmt.Println("received", msg1)
	case msg2 := <-ch2:
		fmt.Println("received", msg2)
	default:
		// This default case makes the select non-blocking
		fmt.Println("no communication")
	}
}
```

#### Ranging Over Channels and Closing

You can iterate over the values received from a channel using a `for...range` loop.
The loop will automatically break when the channel is closed.
It's important for the sender to close the channel to signal that no more values will be sent.

```go
package main

import (
	"fmt"
)

func produce(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch) // Close the channel when done sending
}

func main() {
	ch := make(chan int)
	go produce(ch)

	for value := range ch {
		fmt.Println("Received:", value)
	}
}

```

## Mutual Exclusion

While channels are the preferred way to handle communication between goroutines, sometimes you need to share memory. To
prevent data races (where multiple goroutines access the same memory location concurrently and at least one of the
accesses is a write), you need to use mutual exclusion locks.

The sync package provides this functionality, primarily with sync.Mutex. A Mutex (mutual exclusion lock) ensures that
only one goroutine can access a critical section of code at a time.

### Types of Mutexes

- **sync.Mutex**: The standard mutex.
- **sync.RWMutex**: A reader/writer mutex. It allows any number of readers to access the shared resource simultaneously,
  but
  only one writer at a time. This is beneficial when you have many more reads than writes.

### Common Use Cases for Mutexes

- Protecting access to a shared data structure like a map or a slice that multiple goroutines need to modify.
- Ensuring an operation is performed atomically.
    ```go
    package main
    
    import (
        "fmt"
        "sync"
    )
    
    type SafeCounter struct {
        mu      sync.Mutex
        counter int
    }
    
    func (c *SafeCounter) Increment() {
        c.mu.Lock()
        defer c.mu.Unlock()
        c.counter++
    }
    
    func (c *SafeCounter) Value() int {
        c.mu.Lock()
        defer c.mu.Unlock()
        return c.counter
    }
    
    func main() {
        counter := SafeCounter{}
        var wg sync.WaitGroup
    
        for i := 0; i < 1000; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                counter.Increment()
            }()
        }
    
        wg.Wait()
        fmt.Println("Final counter value:", counter.Value())
    }
    ```

## Goroutine Lifecycle

The lifecycle of a goroutine in Go can be described through the following stages:

1. Creation:
    - A goroutine is created when a function or method call is prefixed with the `go` keyword.
    - The Go runtime allocates a small initial stack for the goroutine, typically a few kilobytes. This stack can grow
      or shrink dynamically as needed.
2. Running:
    - The goroutine starts executing concurrently with other goroutines.
    - The Go runtime scheduler manages the execution of goroutines, switching between them as needed.
3. Blocked:
    - A goroutine may enter a blocked state when it is waiting for a resource, such as:
    - Waiting for I/O operations.
    - Waiting to acquire a lock.
    - Waiting for a channel to receive or send data.
    - When the resource becomes available, the goroutine will be unblocked and resume execution.
4. Terminated:
    - A goroutine terminates when its function or method completes execution.
    - Once terminated, the goroutine is removed from the scheduler.

### Important Take-away:

- The main function is also a goroutine, and when it returns, the program exits.
- The program does not wait for other non-main goroutines to complete.
- Explicit synchronization using channels or other concurrency primitives is required to coordinate and wait for
  goroutines if needed.
- The context package can be used to manage the lifecycle of goroutines, enabling cancellation and timeouts.
- Each goroutine is a separated light-weighted thread which should has its own `defer` function to handle `panic` and
  clean up tasks.

### Context for Lifecycle

The `context` package provides a way to manage the life-cycle of a goroutine, especially for cancellation, timeouts, and
passing request-scoped values.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: stopping\n", id)
			return
		default:
			fmt.Printf("Worker %d: working\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 3; i++ {
		go worker(ctx, i)
	}

	time.Sleep(2 * time.Second)
	cancel()                    // Signal all workers to stop
	time.Sleep(1 * time.Second) // Give workers time to print their stopping message
}
```

### Wait Group Mechanism

A WaitGroup has three main methods you need to know:

1. **Add(delta int)**: This increments the WaitGroup counter by the given delta. You typically call Add(1) before you
   start a goroutine to signal that there is one more task to wait for.
2. **Done()**: This decrements the WaitGroup counter by one. You call this from within the goroutine when it has
   finished its work. A common and robust practice is to use defer wg.Done() at the beginning of the goroutine's
   function to ensure it's always called, even if the function panics.
3. **Wait()**: This blocks the execution of the goroutine that calls it until the internal counter becomes zero. This is
   usually called from the main goroutine to wait for all the spawned workers to complete.

#### Using Wait Group in Goroutines

Here is a step-by-step implementation for a practical use case for `sync.WaitGroup`:

1. Create an instance of sync.WaitGroup.
2. Loop to spawn your worker goroutines.
3. Inside the loop, before starting each goroutine, call wg.Add(1).
4. Start the goroutine. The goroutine function should call wg.Done() when it's finished (preferably using defer).
5. After the loop, call wg.Wait() to block until the counter is back to zero.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// worker simulates a task that takes some time to complete.
// It accepts the WaitGroup pointer so it can signal when it's done.
func worker(id int, wg *sync.WaitGroup) {
	// Defer Done decrements the counter when the function returns.
	// This is the idiomatic way to ensure the WaitGroup is notified.
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)

	// Simulate some work
	time.Sleep(time.Second)

	fmt.Printf("Worker %d done\n", id)
}

func main() {
	// 1. Create a new WaitGroup.
	var wg sync.WaitGroup

	// Number of workers we want to run
	const numWorkers = 5

	fmt.Println("Starting workers...")

	for i := 1; i <= numWorkers; i++ {
		// 2. Increment the WaitGroup counter for each goroutine.
		// It's important to do this *before* launching the goroutine
		// to avoid a race condition where Wait() could be called before Add().
		wg.Add(1)

		// 3. Launch the goroutine.
		go worker(i, &wg)
	}

	fmt.Println("Main goroutine is now waiting for workers to finish...")

	// 4. Block until the WaitGroup counter is zero.
	wg.Wait()

	fmt.Println("All workers have finished. Main goroutine is exiting.")
}
```

## Popular Concurrency Patterns

### Worker Pool

A worker pool is a fixed number of goroutines (workers) that process tasks from a queue (a channel). This pattern is
useful for controlling the number of concurrent operations and preventing resource exhaustion.

```go
package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
```

### Timeout Handling

You can use `select` with a `time.After` channel to implement timeouts for operations.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
```

### Rate Limiting

Rate limiting is essential for controlling resource utilization and maintaining service quality. Go's tickers can be
used to implement a simple rate limiter.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.NewTicker(200 * time.Millisecond)

	for req := range requests {
		<-limiter.C
		fmt.Println("request", req, time.Now())
	}
}
```

### Pipeline

A pipeline is a series of stages connected by channels, where each stage is a goroutine that processes data and passes
it to the next stage.

```go
package main

import "fmt"

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	c := generator(2, 3)
	out := square(c)

	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9
}
```

### Publish-Subscribe

In the Publish-Subscribe (Pub-Sub) pattern, a publisher sends messages to a topic, and multiple subscribers can listen
to that topic without knowing about each other.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Broker struct {
	mu          sync.Mutex
	subscribers map[string][]chan string
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan string),
	}
}

func (b *Broker) Subscribe(topic string) <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 1)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}

func (b *Broker) Publish(topic, msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subscribers[topic] {
		ch <- msg
	}
}

func main() {
	broker := NewBroker()
	sub1 := broker.Subscribe("topic1")
	sub2 := broker.Subscribe("topic1")

	go func() {
		for {
			msg := <-sub1
			fmt.Println("Subscriber 1 received:", msg)
		}
	}()

	go func() {
		for {
			msg := <-sub2
			fmt.Println("Subscriber 2 received:", msg)
		}
	}()

	broker.Publish("topic1", "Hello World!")
	broker.Publish("topic1", "Go is awesome!")

	time.Sleep(time.Second)
}

```

## Common Mistakes

1. **Race Conditions**
    - Occur when multiple goroutines access shared data
    - Solved using mutexes or channel-based synchronization

2. **Deadlocks**
    - Situation where goroutines are waiting for each other
    - Prevented by careful channel and mutex management

3. **Resource Leaks**
    - Goroutines that don't terminate
    - Managed through proper channel closing and context usage

## Best Practices

1. Never start a goroutine that cannot be stopped
2. Prefer channels for communication
3. When memory is used for sharing, use mutexes carefully
4. Control the number of goroutines
5. Handle panics within goroutines
6. Use `context` for cancellation and timeouts
7. Use the Race Detector!

## Practice Exercises

### Exercise 1: Simple Goroutine Counter

Create a program that demonstrates basic goroutine usage with a counter:

### Exercise 2: Channel-Based Number Processor

Build a pipeline using goroutines and channels to process numbers:

### Exercise 3: Worker Pool

Implement a worker pool to distribute tasks among multiple goroutines:

