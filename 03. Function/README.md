# Module 03: Functions in Go

## Table of Contents

<ol>
    <li><a href="#objectives">Objective</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#function-eclaration">Function Declaration</a></li>
    <li><a href="#multiple-return-values">Multiple Return Values</a></li>
    <li><a href="#variadic-functions">Variadic Functions</a></li>
    <li><a href="#anonymous-functions-and-closures">Anonymous Functions and Closures</a></li>
    <li><a href="#defer-panic-and-recover">Defer Panic and Recover</a></li>
    <li><a href="#function-types-and-higher-order-functions">Function Types and Higher-Order Functions</a></li>
    <li><a href="#recursive-functions">Recursive Functions</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercise">Practice Exercise</a></li>
</ol>

## Objectives

By the end of this module, you will:
- Understand Go's function declaration syntax and type system
- Master various parameter and return value patterns
- Learn how to use named return values and variadic functions
- Work with anonymous functions and closures
- Understand defer, panic, and recover mechanisms
- Apply function types and higher-order functions
- Recognize common function design patterns and best practices

## Overview

Functions are the building blocks of Go programs, providing modularity, re-usability, and abstraction. 
Go's approach to functions emphasizes simplicity and clarity while offering powerful features like multiple return values, 
variadic parameters, and first-class function support.

## Function Declaration

Go's function declaration syntax is designed for readability and type safety. 
Unlike some languages, Go uses the `func` keyword and places types after parameter names.

### Simple Function Syntax
```go
// Basic function declaration
func functionName(parameter1 type1, parameter2 type2) returnType {
    // Function body - code goes here
    return returnValue
}
```

### Example: Basic Function
```go
// simple_function.go
package main

import "fmt"

// The add function takes two integers and returns their sum
func add(a int, b int) int {
    return a + b
}

func main() {
    result := add(5, 3)
    fmt.Println("5 + 3 =", result) // Output: 5 + 3 = 8
}
```

### Compact Parameter Type Declaration
When consecutive parameters share the same type, you can declare the type just once after the last parameter in the group:

```go
// Instead of this:
func multiply(a int, b int) int {
    return a * b
}

// You can use this more concise form:
func multiply(a, b int) int {
    return a * b
}
```

## Multiple Return Values

One of Go's most distinctive features is its ability to return multiple values from a function, 
which is especially useful for returning results along with error information.

### Basic Syntax
```go
func functionName(parameters) (returnType1, returnType2) {
    // Function body
    return value1, value2
}
```

### Example: Returning a Result and Error
```go
// divide.go
package main

import (
    "errors"
    "fmt"
)

// The divide function returns both a result and an error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero is not allowed")
    }
    return a / b, nil
}

func main() {
    // Calling a function with multiple return values
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("10 ÷ 2 =", result) // Output: 10 ÷ 2 = 5
    }
    
    // Division by zero example
    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err) // Output: Error: division by zero is not allowed
    } else {
        fmt.Println("Result:", result)
    }
}
```

### Named Return Values
Go allows you to name your return values, 
which creates variables that you can assign to within the function and implicitly return with a "naked" return statement:

```go
// calculate.go
package main

import "fmt"

// Named return values
func calculateRectangleProperties(width, height float64) (area, perimeter float64) {
    // The return variables area and perimeter are pre-declared
    area = width * height
    perimeter = 2 * (width + height)
    
    // Naked return - implicitly returns the named return values
    return
}

func main() {
    a, p := calculateRectangleProperties(5.0, 3.0)
    fmt.Printf("Rectangle with width=5.0, height=3.0:\n")
    fmt.Printf("Area: %.2f\n", a)         // Output: Area: 15.00
    fmt.Printf("Perimeter: %.2f\n", p)     // Output: Perimeter: 16.00
}
```

#### Benefits of Named Return Values:
- Self-documenting return values
- Pre-initialized to zero values
- Useful for documenting the purpose of each return value
- Helpful in long functions with multiple return paths

#### When to Use Named Returns:
- For functions with multiple return values of the same type
- When return value names add clarity
- In longer functions where the return values are modified throughout

#### Best Practice:
While named returns with naked returns are convenient, they should be used judiciously. 
For simple, short functions, explicit returns often provide better clarity.

## Variadic Functions

Variadic functions can accept a variable number of arguments, making them flexible for different calling scenarios.

### Syntax
```go
func functionName(param1 type1, params ...type) returnType {
    // params is a slice of type
}
```

### Example: Basic Variadic Function
```go
// sum_variadic.go
package main

import "fmt"

// sum can take any number of integers
func sum(numbers ...int) int {
    total := 0
    
    // numbers is treated as a slice inside the function
    for _, num := range numbers {
        total += num
    }
    
    return total
}

func main() {
    // Different ways to call a variadic function
    fmt.Println(sum(1, 2))                  // Output: 3
    fmt.Println(sum(1, 2, 3, 4, 5))         // Output: 15
    fmt.Println(sum())                      // Output: 0 (empty slice)
    
    // Passing a slice to a variadic function
    numbers := []int{10, 20, 30, 40}
    fmt.Println(sum(numbers...))           // Output: 100 (unpacking the slice)
}
```

### Combining Regular and Variadic Parameters
Variadic parameters must be the last parameter in the function signature:

```go
// process_values.go
package main

import "fmt"

func processValues(prefix string, values ...int) {
    fmt.Println("Prefix:", prefix)
    fmt.Print("Values: ")
    for _, val := range values {
        fmt.Print(val, " ")
    }
    fmt.Println()
}

func main() {
    processValues("Numbers", 1, 2, 3, 4)
    // Output:
    // Prefix: Numbers
    // Values: 1 2 3 4
}
```

## Anonymous Functions and Closures
Go supports anonymous functions (functions without names) that can be assigned to variables, passed as arguments, or returned from other functions.

### Anonymous Function Syntax
```go
// Assigning to a variable
functionVar := func(parameters) returnType {
    // Function body
}

// Immediate execution
func(parameters) returnType {
    // Function body
}(arguments)
```

### Example: Basic Anonymous Function
```go
// anonymous_functions.go
package main

import "fmt"

func main() {
    // Anonymous function assigned to a variable
    greet := func(name string) {
        fmt.Println("Hello,", name)
    }
    
    // Calling the function through the variable
    greet("Alice")   // Output: Hello, Alice
    
    // Anonymous function with immediate execution
    func(message string) {
        fmt.Println("Message:", message)
    }("This is an anonymous function")
    // Output: Message: This is an anonymous function
}
```

### Closures
A closure is a function that references variables from outside its body. The function may access and assign to the referenced variables.

```go
// closure_example.go
package main

import "fmt"

func main() {
    // Variable outside the anonymous function
    counter := 0
    
    // This function forms a closure by capturing the counter variable
    increment := func() int {
        counter++
        return counter
    }
    
    fmt.Println(increment())  // Output: 1
    fmt.Println(increment())  // Output: 2
    fmt.Println(increment())  // Output: 3
    
    // counter has been modified by the closure
    fmt.Println("Counter:", counter)  // Output: Counter: 3
}
```

### Example: Creating a Counter Function
```go
// counter_generator.go
package main

import "fmt"

// Function that returns another function
func createCounter(start int) func() int {
    // The returned function closes over the start variable
    return func() int {
        start++
        return start
    }
}

func main() {
    // Create two separate counters
    counter1 := createCounter(0)
    counter2 := createCounter(10)
    
    fmt.Println(counter1())  // Output: 1
    fmt.Println(counter1())  // Output: 2
    
    fmt.Println(counter2())  // Output: 11
    fmt.Println(counter2())  // Output: 12
    
    // Each counter maintains its own state
    fmt.Println(counter1())  // Output: 3
}
```

## Defer Panic and Recover
Go provides mechanisms for controlling execution flow in exceptional situations and for cleanup operations.

### Defer
The `defer` statement schedules a function call to be executed immediately before the surrounding function returns, 
regardless of whether it returns normally or with an error.

```go
// defer_example.go
package main

import "fmt"

func main() {
    // defer statements are executed in LIFO order (last in, first out)
    defer fmt.Println("This executes last")
    defer fmt.Println("This executes second")
    defer fmt.Println("This executes first")
    
    fmt.Println("Regular execution")
}
/* Output:
Regular execution
This executes first
This executes second
This executes last
*/
```

#### Common use cases for Defer:
- Resource cleanup (closing files, network connections)
- Unlocking mutexes
- Executing required operations regardless of which path a function takes to return

### Example: File Operations with Defer
```go
// file_with_defer.go
package main

import (
    "fmt"
    "os"
)

func readFile(filename string) (string, error) {
    // Open the file
    file, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    
    // Schedule the file to close when the function returns
    defer file.Close()
    
    // Read the file (simplified for example)
    buffer := make([]byte, 100)
    count, err := file.Read(buffer)
    if err != nil {
        return "", err
    }
    
    return string(buffer[:count]), nil
}

func main() {
    content, err := readFile("example.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("File content:", content)
}
```

### Panic
Panic is a built-in function that stops the normal execution of the current goroutine. 
When a function calls `panic`, normal execution stops, deferred functions are executed, and control returns to the caller.

```go
// panic_example.go
package main

import "fmt"

func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

func main() {
    fmt.Println("Starting program")
    
    defer fmt.Println("This will still execute")
    
    // This will cause a panic
    result := divide(10, 0)
    
    // This code won't execute
    fmt.Println("Result:", result)
}
/* Output:
Starting program
This will still execute
panic: division by zero
goroutine 1 [running]:
main.divide(...)
...
*/
```

### Recover
Recover is a built-in function that regains control of a panicking goroutine. 
It's only useful inside deferred functions.

```go
// recover_example.go
package main

import "fmt"

func recoverExample() {
    // A deferred function that will handle any panic
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    
    fmt.Println("About to panic")
    panic("something went wrong")
    fmt.Println("This won't execute")
}

func main() {
    fmt.Println("Starting program")
    
    recoverExample()
    
    fmt.Println("Program continues normally")
}
/* Output:
Starting program
About to panic
Recovered from panic: something went wrong
Program continues normally
*/
```

## Function Types and Higher-Order Functions

In Go, functions are first-class citizens, which means:
- Functions can be assigned to variables
- Functions can be passed as arguments to other functions
- Functions can be returned from other functions

### Function Types
A function type defines the signature of a function without specifying its implementation:

```go
// function_type.go
package main

import "fmt"

// Declare a function type
type MathFunc func(int, int) int

// Function that takes a MathFunc type as a parameter
func applyMathFunc(f MathFunc, a, b int) int {
    return f(a, b)
}

func main() {
    // Define functions that match the MathFunc type
    add := func(x, y int) int { return x + y }
    multiply := func(x, y int) int { return x * y }
    
    // Use the functions with applyMathFunc
    fmt.Println("10 + 5 =", applyMathFunc(add, 10, 5))      // Output: 10 + 5 = 15
    fmt.Println("10 × 5 =", applyMathFunc(multiply, 10, 5)) // Output: 10 × 5 = 50
}
```

### Higher-Order Functions
Higher-order functions either take functions as arguments or return functions:

#### Example: Function That Returns a Function
```go
// higher_order.go
package main

import "fmt"

// Returns a greeter function customized with the greeting parameter
func createGreeter(greeting string) func(string) string {
    // Return a function that uses the captured greeting
    return func(name string) string {
        return greeting + ", " + name + "!"
    }
}

func main() {
    // Create specialized greeter functions
    englishGreeter := createGreeter("Hello")
    spanishGreeter := createGreeter("Hola")
    
    // Use the greeter functions
    fmt.Println(englishGreeter("John"))  // Output: Hello, John!
    fmt.Println(spanishGreeter("Maria")) // Output: Hola, Maria!
}
```

#### Example: Functional Programming
```go
// functional_example.go
package main

import "fmt"

// Function that applies a transformation to each element in a slice
func mapSlice(input []int, transformFunc func(int) int) []int {
    result := make([]int, len(input))
    for i, val := range input {
        result[i] = transformFunc(val)
    }
    return result
}

// Function that filters elements based on a condition
func filterSlice(input []int, keepFunc func(int) bool) []int {
    var result []int
    for _, val := range input {
        if keepFunc(val) {
            result = append(result, val)
        }
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Map: Square all numbers
    squared := mapSlice(numbers, func(x int) int {
        return x * x
    })
    fmt.Println("Squared:", squared)
    // Output: Squared: [1 4 9 16 25 36 49 64 81 100]
    
    // Filter: Keep only even numbers
    evens := filterSlice(numbers, func(x int) bool {
        return x%2 == 0
    })
    fmt.Println("Even numbers:", evens)
    // Output: Even numbers: [2 4 6 8 10]
    
    // Combine operations: square all numbers, then keep only those greater than 50
    bigSquares := filterSlice(
        mapSlice(numbers, func(x int) int { return x * x }),
        func(x int) bool { return x > 50 },
    )
    fmt.Println("Squares > 50:", bigSquares)
    // Output: Squares > 50: [64 81 100]
}
```

## Recursive Functions
A recursive function is one that calls itself directly or indirectly:

```go
// recursion_example.go
package main

import "fmt"

// Factorial function using recursion
func factorial(n uint) uint {
    // Base case
    if n == 0 {
        return 1
    }
    
    // Recursive case
    return n * factorial(n-1)
}

// Fibonacci function using recursion
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    fmt.Println("Factorial of 5:", factorial(5))  // Output: Factorial of 5: 120
    
    fmt.Println("Fibonacci sequence:")
    for i := 0; i < 10; i++ {
        fmt.Print(fibonacci(i), " ")  // Output: 0 1 1 2 3 5 8 13 21 34
    }
    fmt.Println()
}
```

### Note on Recursive Functions in Go
While Go supports recursion, it doesn't optimize tail recursion like some functional languages do. 
For deep recursion, consider using iteration or implementing continuation-passing style to avoid stack overflow.

## Best Practices

### 1. Single Responsibility
Each function should have a single, well-defined purpose.

```go
// Bad: Function doing too many things
func processUserData(user User) (bool, error) {
    // Validate user
    // Update database
    // Send notification email
    // Log audit trail
}

// Better: Split into focused functions
func validateUser(user User) error { /* ... */ }
func saveUser(user User) error { /* ... */ }
func notifyUser(user User) error { /* ... */ }
func auditUserChange(user User) error { /* ... */ }
```

### 2. Keep Functions Small
Shorter functions are easier to understand, test, and maintain.

### 3. Descriptive Function Names
Function names should clearly indicate what the function does:

```go
// Too vague
func process(s string) string

// Better - clear about purpose
func sanitizeUserInput(input string) string
```

### 4. Parameter Count
Limit the number of parameters (ideally 3 or fewer). If you need many parameters, consider using a struct:

```go
// Too many parameters
func createUser(firstName, lastName, email, phone, address, city, state, country, zip string, age int) User

// Better approach
type UserInfo struct {
    FirstName, LastName string
    Email, Phone string
    Address, City, State, Country, Zip string
    Age int
}

func createUser(info UserInfo) User
```

### 5. Error Handling
Always check and handle errors returned by functions:

```go
result, err := someFunction()
if err != nil {
    // Handle the error
    return err
}
// Use result safely
```

### 6. Return Early
Return as soon as you know the answer to reduce nesting and complexity:

```go
// Early return pattern
func processData(data []int) (int, error) {
    if len(data) == 0 {
        return 0, errors.New("empty data set")
    }
    
    if !isValid(data) {
        return 0, errors.New("invalid data")
    }
    
    // Process valid data...
    return result, nil
}
```

### 7. Document Your Functions
Use comments to explain what functions do, especially for exported functions:

```go
// CalculateDiscount determines the discount amount based on the 
// purchase total and customer type. It returns the discount amount
// and an error if the calculation cannot be performed.
func CalculateDiscount(total float64, customerType string) (float64, error) {
    // Implementation...
}
```

## Practice Exercise

### Exercise 1: Utility Functions Library
Create a library of simple utility functions that can be used for common tasks. 
This exercise will help you practice defining and organizing functions in Go.

Specifically, you need to:
1. Create a `StringUtils` type with methods for string manipulation:
    - A `Reverse` function that reverses the characters in a string
    - An `IsPalindrome` function that checks if a string reads the same forward and backward

2. Create a `MathUtils` type with methods for mathematical operations:
    - A `Factorial` function that calculates n! (factorial) for a given number
    - An `IsPrime` function that determines whether a number is prime

3. Test each function in the `main` function with example inputs to demonstrate how they work

### Exercise 2: Function Composition
Practice creating higher-order functions that compose various operations. 
This exercise will introduce you to functional programming concepts in Go 
where functions can be passed as arguments and returned as values.

Implement the following:
1. Create function types (`StringProcessor` and `IntProcessor`) that represent functions which transform strings and integers
2. Implement a `ComposeStringProcessors` function that combines two string processing functions into a single function
3. Implement a `Chain` function that applies multiple string processors in sequence
4. Create example string processors (trim spaces, convert to uppercase, reverse string)
5. Demonstrate function composition by creating and applying composed functions to test data

### Exercise 3: Advanced Calculator with Function Types
Build a flexible calculator application that can be extended with new operations. 
This exercise demonstrates how to use function types to create a plugin architecture.

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
