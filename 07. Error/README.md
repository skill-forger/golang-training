# Module 07: Error Handling in Go

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#go-error-handling-philosophy">Go Error Handling Philosophy</a></li>
    <li><a href="#basic-error-handling">Basic Error Handling</a></li>
    <li><a href="#custom-error-types">Custom Error Types</a></li>
    <li><a href="#error-wrapping-and-context">Error Wrapping and Context</a></li>
    <li><a href="#error-handling-patterns">Error Handling Patterns</a></li>
    <li><a href="#panic-and-recover">Panic and Recover</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#error-handling-in-concurrency">Error Handling in Concurrency</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand Go's error handling philosophy and approach
- Master the error interface and its implementations
- Learn to create and use custom error types
- Apply error wrapping techniques for better context
- Implement error handling patterns and best practices
- Use advanced error handling including panic and recover
- Avoid common error handling pitfalls

## Overview

Error handling is a critical aspect of writing reliable, maintainable Go programs.
Unlike many modern languages that use exceptions for error handling,
Go takes a different approach by treating errors as values that are explicitly returned, checked, and handled.
This philosophy aligns with Go's emphasis on explicitness and simplicity.

## Go Error Handling Philosophy

### Errors as Values

The cornerstone of Go's error handling is that errors are just values - not special exceptions or control flow
mechanisms. This means errors are:

- Returned from functions like any other value
- Explicitly checked by the caller
- Handled using standard control flow constructs
- Can be created, stored, and manipulated like any other value

### The Error Interface

At the core of Go's error handling is the built-in `error` interface:

```go
// The error interface from the standard library
type error interface {
Error() string
}
```

This simple interface has just one method that returns a string description of the error. Any type that implements this
method is considered an error in Go.

## Basic Error Handling

### Returning and Checking Errors

```go
// basic_errors.go
package main

import (
	"errors"
	"fmt"
)

// Function that returns a result and potential error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// Create and return a simple error
		return 0, errors.New("division by zero")
	}
	return a / b, nil // nil indicates no error
}

func main() {
	// Call the function and get both result and error
	result, err := divide(10, 2)

	// Check if an error occurred
	if err != nil {
		fmt.Println("Error:", err)
		return // Handle the error by returning early
	}

	// Continue with normal processing if no error
	fmt.Println("Result:", result)

	// Another call with different parameters
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err) // This will print the error
		return
	}
	fmt.Println("Result:", result) // This won't execute
}
```

### Creating Errors

Go offers several ways to create errors:

```go
// creating_errors.go
package main

import (
	"errors" // For simple errors
	"fmt"    // For formatted errors
)

func main() {
	// Method 1: Using errors.New for simple static error messages
	err1 := errors.New("something went wrong")

	// Method 2: Using fmt.Errorf for formatted error messages
	name := "file.txt"
	err2 := fmt.Errorf("failed to open %s", name)

	fmt.Println(err1) // Prints: something went wrong
	fmt.Println(err2) // Prints: failed to open file.txt
}
```

## Custom Error Types

For more sophisticated error handling, you can create custom error types:

```go
// custom_errors.go
package main

import (
	"fmt"
)

// Define a custom error type with additional context
type ValidationError struct {
	Field   string
	Message string
}

// Implement the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Function that returns our custom error
func validateUsername(username string) error {
	if len(username) < 3 {
		// Return a custom error with specific context
		return &ValidationError{
			Field:   "username",
			Message: "must be at least 3 characters",
		}
	}
	return nil
}

func main() {
	// Test with invalid username
	err := validateUsername("ab")
	if err != nil {
		fmt.Println("Error:", err)

		// Type assertion to access custom error fields
		if validationErr, ok := err.(*ValidationError); ok {
			fmt.Printf("Field: %s, Message: %s\n",
				validationErr.Field, validationErr.Message)
		}
	}
}
```

### Benefits of Custom Error Types

1. **Richer Context**: Include fields that provide additional information
2. **Type Safety**: Handle specific error types differently
3. **Behavior**: Add methods for special error handling
4. **Semantic Clarity**: Make error types that match your domain

## Error Wrapping and Context

Go 1.13 introduced better error wrapping capabilities:

```go
// error_wrapping.go
package main

import (
	"errors"
	"fmt"
)

// Define some sentinel errors
var (
	ErrNotFound   = errors.New("resource not found")
	ErrPermission = errors.New("permission denied")
)

// Function that uses a sentinel error
func findResource(id string) error {
	// Simulate a not found error
	return ErrNotFound
}

// Function that wraps errors with context
func processResource(id string) error {
	err := findResource(id)
	if err != nil {
		// Wrap the original error with additional context
		// The %w verb preserves the original error for unwrapping
		return fmt.Errorf("processing resource %s: %w", id, err)
	}
	return nil
}

func main() {
	err := processResource("resource-123")
	if err != nil {
		fmt.Println("Error:", err)

		// Check if the wrapped error is a specific error
		if errors.Is(err, ErrNotFound) {
			fmt.Println("The resource was not found!")
		}

		// Get the wrapped error directly
		fmt.Println("Unwrapped:", errors.Unwrap(err))
	}
}
```

### The Wrapping Mechanism

Error wrapping lets you:

1. **Add Context**: Enrich errors as they travel up the call stack
2. **Preserve Original Errors**: Keep underlying errors for precise checks
3. **Create Error Chains**: Build a chain of wrapped errors for detailed tracking

### Checking Wrapped Errors

Go 1.13 added two key functions for working with wrapped errors:

- `errors.Is(err, target)`: Checks if `err` or any error it wraps equals `target`
- `errors.As(err, target)`: Finds the first error in the chain that matches the type of `target`

```go
// error_checking.go
package main

import (
	"errors"
	"fmt"
	"os"
)

// Custom error type
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query error for %q: %v", e.Query, e.Err)
}

// Implement Unwrap to make it part of the error chain
func (e *QueryError) Unwrap() error {
	return e.Err
}

func searchDatabase(query string) error {
	// Simulate a file not found error wrapped in our custom error
	return &QueryError{
		Query: query,
		Err:   os.ErrNotExist,
	}
}

func main() {
	err := searchDatabase("SELECT * FROM users")

	// Check for specific error type
	var queryErr *QueryError
	if errors.As(err, &queryErr) {
		fmt.Printf("Query error occurred: %s\n", queryErr.Query)
	}

	// Check for specific error value, even if wrapped
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("The underlying error is file not found")
	}
}
```

## Error Handling Patterns

### The Sentinel Error Pattern

Sentinel errors are predefined error values that can be compared directly:

```go
// sentinel_errors.go
package main

import (
	"errors"
	"fmt"
)

// Sentinel errors - predefined error values
var (
	ErrInsufficientFunds = errors.New("insufficient funds in account")
	ErrAccountLocked     = errors.New("account is locked")
	ErrInvalidAmount     = errors.New("invalid transaction amount")
)

func WithdrawMoney(account string, amount float64) error {
	// Validate amount
	if amount <= 0 {
		return ErrInvalidAmount
	}

	// Check if account is locked (simulated)
	isLocked := (account == "locked123")
	if isLocked {
		return ErrAccountLocked
	}

	// Check if sufficient funds (simulated)
	balance := 100.0 // Pretend this comes from a database
	if balance < amount {
		return ErrInsufficientFunds
	}

	// Process withdrawal
	fmt.Printf("Withdrew $%.2f from account %s\n", amount, account)
	return nil
}

func main() {
	// Try different scenarios
	accounts := []string{"valid123", "locked123", "empty456"}
	amounts := []float64{50.0, -10.0, 200.0}

	for _, acc := range accounts {
		for _, amt := range amounts {
			err := WithdrawMoney(acc, amt)

			if err != nil {
				switch {
				case errors.Is(err, ErrInvalidAmount):
					fmt.Printf("Cannot withdraw $%.2f: invalid amount\n", amt)

				case errors.Is(err, ErrAccountLocked):
					fmt.Printf("Account %s is locked\n", acc)

				case errors.Is(err, ErrInsufficientFunds):
					fmt.Printf("Not enough funds to withdraw $%.2f\n", amt)

				default:
					fmt.Printf("Unknown error: %v\n", err)
				}
			}
		}
	}
}
```

### Error Handling in HTTP Servers

Error handling in web applications has specific patterns:

```go
// http_errors.go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Custom error type with HTTP status code
type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// User-related functions that can fail
func getUserByID(id string) (string, error) {
	if id == "admin" {
		return "Administrator", nil
	}
	return "", &HTTPError{
		Code:    http.StatusNotFound,
		Message: "user not found",
	}
}

// HTTP handler with error handling
func userHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	username, err := getUserByID(userID)
	if err != nil {
		// Extract HTTP error if possible
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.Code)
		} else {
			// Generic server error for unexpected errors
			log.Printf("Unexpected error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "User: %s", username)
}

func main() {
	http.HandleFunc("/user", userHandler)
	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Panic and Recover

While Go emphasizes explicit error handling, it provides `panic` and `recover` for exceptional situations:

```go
// panic_recover.go
package main

import (
	"fmt"
)

// Function that might panic
func riskyOperation(data string) (result string) {
	// Set up recovery mechanism
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			result = "fallback result after panic"
		}
	}()

	// This could cause a panic
	if len(data) == 0 {
		panic("cannot process empty data")
	}

	return "processed: " + data
}

func main() {
	// Test with valid data
	result := riskyOperation("hello")
	fmt.Println("Result:", result)

	// Test with data that will cause panic
	result = riskyOperation("")
	fmt.Println("Result after panic:", result)

	fmt.Println("Program continues executing")
}
```

### When to Use Panic and Recover

- **Use panic**: For truly exceptional conditions that should not be handled normally
- **Use recover**: To prevent a panic from crashing your entire program
- **Prefer error returns**: For most error conditions that can be reasonably expected

## Common Mistakes

### 1. Ignoring Errors

```go
// BAD: Error ignored
file, _ := os.Open("config.txt") // Don't do this!

// GOOD: Error handled
file, err := os.Open("config.txt")
if err != nil {
log.Printf("Failed to open config: %v", err)
// Handle error appropriately
return
}
```

### 2. Using String Comparison for Error Checking

```go
// BAD: Comparing error strings
if err != nil && err.Error() == "file not found" {
// Unreliable and fragile
}

// GOOD: Using errors.Is
if err != nil && errors.Is(err, os.ErrNotExist) {
// Reliable even with wrapped errors
}
```

### 3. Insufficient Error Context

```go
// BAD: Generic error
return errors.New("operation failed")

// GOOD: Contextual error
return fmt.Errorf("user update failed for ID %s: %w", userID, err)
```

### 4. Excessive Panic Usage

```go
// BAD: Using panic for normal error conditions
func getData(id string) []byte {
data, err := ioutil.ReadFile(id + ".dat")
if err != nil {
panic(err) // Don't do this for expected errors!
}
return data
}

// GOOD: Return errors normally
func getData(id string) ([]byte, error) {
return ioutil.ReadFile(id + ".dat")
}
```

## Error Handling in Concurrency

Error handling in goroutines requires special attention:

```go
// concurrent_errors.go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Process URLs concurrently and collect errors
func processURLs(urls []string) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Simulate processing
			if url == "https://example.com/error" {
				errChan <- fmt.Errorf("failed to process %s", url)
			} else if url == "https://example.com/timeout" {
				time.Sleep(100 * time.Millisecond)
				errChan <- fmt.Errorf("timeout processing %s", url)
			}
			// Success case doesn't send an error
		}(url)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Collect all errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	return errs
}

func main() {
	urls := []string{
		"https://example.com/success",
		"https://example.com/error",
		"https://example.com/timeout",
	}

	fmt.Println("Processing URLs...")
	errors := processURLs(urls)

	if len(errors) > 0 {
		fmt.Println("Encountered errors:")
		for i, err := range errors {
			fmt.Printf("%d: %v\n", i+1, err)
		}
	} else {
		fmt.Println("All URLs processed successfully")
	}
}
```

## Best Practices

### 1. Be Explicit and Check Errors

Always check returned errors and handle them appropriately:

```go
file, err := os.Open("file.txt")
if err != nil {
// Handle the error - never ignore it!
log.Printf("Error opening file: %v", err)
return
}
defer file.Close()
```

### 2. Add Context to Errors

Make errors more informative by adding context:

```go
func processConfig(path string) error {
data, err := readFile(path)
if err != nil {
return fmt.Errorf("config processing failed: %w", err)
}
// Process data...
return nil
}
```

### 3. Design Error Types for Your Domain

Create error types that match your application's domain:

```go
type DatabaseError struct {
Query  string
Err    error
Table  string
UserID string
}

func (e *DatabaseError) Error() string {
return fmt.Sprintf("database error on table %s: %v", e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
return e.Err
}
```

### 4. Handle Errors at the Appropriate Level

Not every function needs to handle every error:

```go
// Low-level function just returns the error
func readUserData(id string) ([]byte, error) {
return ioutil.ReadFile("users/" + id + ".json")
}

// Mid-level function adds context but doesn't handle
func getUserProfile(id string) (*Profile, error) {
data, err := readUserData(id)
if err != nil {
return nil, fmt.Errorf("get profile %s failed: %w", id, err)
}
// Process data...
return profile, nil
}

// High-level function handles the error appropriately
func handleUserRequest(w http.ResponseWriter, id string) {
profile, err := getUserProfile(id)
if err != nil {
if errors.Is(err, os.ErrNotExist) {
http.Error(w, "User not found", http.StatusNotFound)
} else {
http.Error(w, "Internal error", http.StatusInternalServerError)
log.Printf("User request failed: %v", err)
}
return
}
// Use profile...
}
```

### 5. Only Use Panic for Truly Unrecoverable Situations

Reserve panic for programmer errors or truly exceptional conditions:

```go
func MustCompileRegex(pattern string) *regexp.Regexp {
re, err := regexp.Compile(pattern)
if err != nil {
// This is a programmer error - the pattern should be valid
panic(fmt.Sprintf("Invalid regex pattern %q: %v", pattern, err))
}
return re
}

// Expected to be used in initialization, not regular operation
var validNameRe = MustCompileRegex(`^[a-zA-Z][a-zA-Z0-9_]{2,29}$`)
```

## Practice Exercises

### Exercise 1: Building a Robust API Client

Create a resilient API client that demonstrates comprehensive error handling practices in Go.
This exercise will teach you how to design clear error types,
use the error wrapping features, and implement error-based control flow.

Your implementation should include:

1. Custom error types that:
    - Include relevant contextual information (status code, URL, message)
    - Implement the `Error()` method for clear error messages
    - Support unwrapping via the `Unwrap()` method
2. Sentinel (predefined) errors for common error conditions:
    - `ErrNotFound` for missing resources
    - `ErrUnauthorized` for authentication failures
    - `ErrTimeout` for request timeouts
3. An `APIClient` struct with methods for:
    - Making HTTP requests with proper error handling
    - Handling different error cases (timeouts, HTTP error codes)
    - Wrapping underlying errors with context
4. A demonstration showing:
    - Proper error checking with specific error types
    - Using `errors.Is()` to check for sentinel errors
    - Using `errors.As()` to extract information from custom errors
    - Providing appropriate feedback based on error types

### Exercise 2: Database Connection with Error Recovery

Develop a database connection manager that handles various error scenarios and implements automatic recovery strategies.
This exercise shows how to use errors to make robust systems that can recover from failures.

Your implementation should include:

1. Custom error types for different database issues:
    - Connection errors
    - Query execution errors
    - Transaction errors
2. A `DBConnector` that manages database connections with:
    - Connection retry logic with exponential backoff
    - Transaction handling with proper rollback on errors
    - Query execution with timeout handling
3. A `QueryExecutor` interface with implementations for:
    - Basic queries
    - Prepared statements
    - Transactions
4. A demonstration showing:
    - Handling temporary network failures
    - Automatic reconnection after errors
    - Transaction rollback on errors
    - Proper resource cleanup

### Exercise 3: File Processing with Error Logging

Build a file processing system that demonstrates advanced error handling techniques,
including logging, error wrapping, and recovery mechanisms.
This exercise shows how to handle I/O errors and create user-friendly error messages.

Your implementation should include:

1. A hierarchical error handling system with:
    - Base error types for I/O operations
    - Specialized errors for different file processing stages
    - Error aggregation for batch operations
2. A structured logging system that:
    - Records errors with appropriate context
    - Categorizes errors by severity
    - Provides detailed debugging information
3. Recovery mechanisms that:
    - Skip problematic files and continue processing others
    - Attempt alternative processing methods when primary methods fail
    - Clean up resources even when errors occur
4. A demonstration that processes multiple files and shows how the system handles various error conditions
