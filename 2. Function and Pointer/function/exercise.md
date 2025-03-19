## Practical Exercises

### Exercise 1: Utility Functions Library

Create a library of simple utility functions that can be used for common tasks. This exercise will help you practice defining and organizing functions in Go.

Specifically, you need to:
1. Create a `StringUtils` type with methods for string manipulation:
   - A `Reverse` function that reverses the characters in a string
   - An `IsPalindrome` function that checks if a string reads the same forward and backward

2. Create a `MathUtils` type with methods for mathematical operations:
   - A `Factorial` function that calculates n! (factorial) for a given number
   - An `IsPrime` function that determines whether a number is prime

3. Test each function in the `main` function with example inputs to demonstrate how they work

```go
// utils.go
package main

import (
    "fmt"
    "strings"
)

// StringUtils contains utility functions for string manipulation
type StringUtils struct{}

// Reverse returns the reversed version of the input string
func (su StringUtils) Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// IsPalindrome checks if a string reads the same backward as forward
func (su StringUtils) IsPalindrome(s string) bool {
    // Convert to lowercase and remove spaces for case-insensitive comparison
    s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
    return s == su.Reverse(s)
}

// MathUtils contains utility functions for mathematical operations
type MathUtils struct{}

// Factorial calculates the factorial of a number
func (mu MathUtils) Factorial(n uint) uint {
    if n == 0 {
        return 1
    }
    return n * mu.Factorial(n-1)
}

// IsPrime determines if a number is prime
func (mu MathUtils) IsPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    
    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}

func main() {
    su := StringUtils{}
    mu := MathUtils{}
    
    // Test string utilities
    fmt.Println("Reversed 'hello':", su.Reverse("hello"))
    fmt.Println("Is 'radar' a palindrome?", su.IsPalindrome("radar"))
    fmt.Println("Is 'A man a plan a canal Panama' a palindrome?", 
        su.IsPalindrome("A man a plan a canal Panama"))
    
    // Test math utilities
    fmt.Println("Factorial of 5:", mu.Factorial(5))
    fmt.Println("Is 17 prime?", mu.IsPrime(17))
    fmt.Println("Is 20 prime?", mu.IsPrime(20))
}
```

### Exercise 2: Function Composition

Practice creating higher-order functions that compose various operations. This exercise will introduce you to functional programming concepts in Go where functions can be passed as arguments and returned as values.

Implement the following:
1. Create function types (`StringProcessor` and `IntProcessor`) that represent functions which transform strings and integers
2. Implement a `ComposeStringProcessors` function that combines two string processing functions into a single function
3. Implement a `Chain` function that applies multiple string processors in sequence
4. Create example string processors (trim spaces, convert to uppercase, reverse string)
5. Demonstrate function composition by creating and applying composed functions to test data

```go
// function_composition.go
package main

import (
    "fmt"
    "strings"
)

// Function types
type StringProcessor func(string) string
type IntProcessor func(int) int

// Compose two string processors into a single processor
func ComposeStringProcessors(f, g StringProcessor) StringProcessor {
    return func(s string) string {
        return g(f(s))
    }
}

// Chain applies a series of string processors in sequence
func Chain(processors ...StringProcessor) StringProcessor {
    return func(s string) string {
        result := s
        for _, processor := range processors {
            result = processor(result)
        }
        return result
    }
}

func main() {
    // Define some string processors
    trim := func(s string) string { return strings.TrimSpace(s) }
    upper := func(s string) string { return strings.ToUpper(s) }
    reverse := func(s string) string {
        runes := []rune(s)
        for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
            runes[i], runes[j] = runes[j], runes[i]
        }
        return string(runes)
    }
    
    // Compose processors
    trimAndUpper := ComposeStringProcessors(trim, upper)
    fmt.Println(trimAndUpper("  hello world  "))  // Output: HELLO WORLD
    
    // Chain multiple processors
    processAll := Chain(trim, upper, reverse)
    fmt.Println(processAll("  hello world  "))  // Output: DLROW OLLEH
    
    // Create and apply custom chains
    emphasize := Chain(
        trim,
        upper,
        func(s string) string { return "*** " + s + " ***" },
    )
    fmt.Println(emphasize("  important message  "))  // Output: *** IMPORTANT MESSAGE ***
}
```

### Exercise 3: Advanced Calculator with Function Types

Build a flexible calculator application that can be extended with new operations. This exercise demonstrates how to use function types to create a plugin architecture.

Your implementation should:
1. Define an `Operation` function type that takes two float64 values and returns a result with a possible error
2. Implement basic arithmetic operations (Add, Subtract, Multiply, Divide, Power)
3. Create a `Calculator` struct that stores operations mapped to symbol strings
4. Include a method to register new operations at runtime
5. Implement a Calculate method that performs the specified operation by symbol
6. Demonstrate the calculator by:
   - Performing basic operations
   - Adding a custom operation (e.g., average)
   - Handling errors (e.g., division by zero)

```go
// calculator.go
package main

import (
    "errors"
    "fmt"
    "math"
)

// Define operation function type
type Operation func(float64, float64) (float64, error)

// Basic operations
func Add(a, b float64) (float64, error) {
    return a + b, nil
}

func Subtract(a, b float64) (float64, error) {
    return a - b, nil
}

func Multiply(a, b float64) (float64, error) {
    return a * b, nil
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func Power(a, b float64) (float64, error) {
    return math.Pow(a, b), nil
}

// Calculator holds operations and provides methods to use them
type Calculator struct {
    operations map[string]Operation
}

// NewCalculator creates a new calculator with standard operations
func NewCalculator() *Calculator {
    calc := &Calculator{
        operations: make(map[string]Operation),
    }
    
    // Register basic operations
    calc.RegisterOperation("+", Add)
    calc.RegisterOperation("-", Subtract)
    calc.RegisterOperation("*", Multiply)
    calc.RegisterOperation("/", Divide)
    calc.RegisterOperation("^", Power)
    
    return calc
}

// RegisterOperation adds a new operation to the calculator
func (c *Calculator) RegisterOperation(symbol string, op Operation) {
    c.operations[symbol] = op
}

// Calculate performs the specified operation
func (c *Calculator) Calculate(a, b float64, symbol string) (float64, error) {
    operation, found := c.operations[symbol]
    if !found {
        return 0, fmt.Errorf("unknown operation: %s", symbol)
    }
    
    return operation(a, b)
}

func main() {
    calc := NewCalculator()
    
    // Use the calculator
    result, err := calc.Calculate(10, 5, "+")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("10 + 5 =", result)  // Output: 10 + 5 = 15
    }
    
    // Add a custom operation
    calc.RegisterOperation("avg", func(a, b float64) (float64, error) {
        return (a + b) / 2, nil
    })
    
    result, _ = calc.Calculate(10, 20, "avg")
    fmt.Println("Average of 10 and 20:", result)  // Output: Average of 10 and 20: 15
    
    // Try division by zero
    result, err = calc.Calculate(10, 0, "/")
    if err != nil {
        fmt.Println("Error:", err)  // Output: Error: division by zero
    }
}
```
