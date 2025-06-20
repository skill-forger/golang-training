# Module 01: Introduction to Go Programming

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview-of-go">Overview of Go</a></li>
    <li><a href="#why-learn-go">Why Learn Go</a></li>
    <li><a href="#understanding-go-workspace">Understanding Go Workspace</a></li>
    <li><a href="#your-first-go-program">Your First Go Program</a></li>
    <li><a href="#understanding-the-basic-structure">Understanding the Basic Structure</a></li>
    <li><a href="#variables-and-data-types">Variables and Data Types</a></li>
    <li><a href="#type-conversion">Type Conversion</a></li>
    <li><a href="#common-pitfalls-and-strict-rules">Common Pitfalls and Strict Rules</a></li>
    <li><a href="#practice-exercise">Practice Exercise</a></li>
</ol>

## Objectives

- Understand Go's design philosophy and use cases
- Write and run your first Go programs
- Understand package organization and import system
- Learn about Go's variable declaration styles and basic data types
- Get familiar with Go's code formatting conventions

## Overview of Go

Go (or Golang) is an open-source programming language developed by Google in 2007 and released publicly in 2009.
It was designed by Robert Griesemer, Rob Pike, and Ken Thompson with the following goals:

- **Simplicity**: Clean syntax and minimal language features
- **Efficiency**: Fast compilation, efficient execution
- **Safety**: Strong static typing and memory safety
- **Concurrency**: Built-in support for concurrent programming
- **Modern**: Designed for modern multicore computers and networked systems
  Go combines the performance of compiled languages like C++ with the simplicity and readability of languages like
  Python.

## Why Learn Go

- **Industry Adoption**: Used by companies like Google, Uber, Dropbox, Netflix
- **Performance**: Excellent for high-performance applications
- **Concurrency**: Powerful concurrency model with goroutines and channels
- **Standard Library**: Rich standard library for common tasks
- **Tooling**: Excellent built-in tools for formatting, testing, and documentation
- **Growing Ecosystem**: Increasing number of libraries and frameworks

## Understanding Go Workspace

### Go Modules (Modern Approach)

Since Go 1.11, the recommended way to manage dependencies is with Go Modules:

```bash
# Create a new project directory anywhere you like
mkdir my-go-project
cd my-go-project

# Initialize a new module
go mod init github.com/yourusername/my-go-project
```

The `go.mod` file will track your dependencies and versions.

### Traditional GOPATH (Legacy)

Before Go modules, all Go code had to reside in a specific workspace structure:

```
$GOPATH/
  ├── bin/    # Compiled executable programs
  ├── pkg/    # Compiled package objects
  └── src/    # Source code organized by repository
```

Modern Go development uses modules, so you don't need to worry about GOPATH much.

## Your First Go Program

### Hello World

Create a file named `hello_world.go`:

```go
// hello_world.go - A simple Hello World program in Go
package main // Declares this file belongs to the main package

import "fmt" // Import the formatting package from standard library

// The main function is the entry point of the program
func main() {
	// Print a message to the console
	fmt.Println("Hello, Go!")

	// We can print multiple lines
	fmt.Println("My name is Gopher!")
	fmt.Println("I'm learning Go programming!")
}
```

### Running Your Program

```bash
# Method 1: Directly run the file
go run hello_world.go

# Method 2: Run all Go files in the current directory
go run .

# Method 3: Build an executable and then run it
go build -o hello
./hello  # On Windows: hello.exe
```

## Understanding the Basic Structure

Every Go program consists of:

1. **Package declaration**: `package main`
    - Executable programs must have a `main` package
    - Libraries use other package names

2. **Import statements**: `import "fmt"`
    - Import packages from standard library or third-party sources
    - Multiple imports can be grouped:
      ```go
      import (
          "fmt"
          "strings"
          "time"
      )
      ```

3. **Functions**: `func main() { ... }`
    - `main()` is the entry point for executable programs
    - Other functions can be defined and called

4. **Comments**:
    - Single-line comments: `// This is a comment`
    - Multi-line comments: `/* This is a multi-line comment */`

## Variables and Data Types

### Variable Declaration Styles

Go provides several ways to declare variables:

```go
package main

import "fmt"

func main() {
	// 1. Full declaration with type and initial value
	var name string = "Gopher"

	// 2. Type inference - type is determined from the value
	var age = 25

	// 3. Short declaration (most common inside functions)
	company := "Google" // := is only used for initial declaration

	// 4. Multiple variable declaration
	var (
		username = "developer"
		isActive = true
		score    = 95.5
	)

	// 5. Zero-value initialization (without initial value)
	var counter int    // Initialized to 0
	var message string // Initialized to "" (empty string)
	var valid bool     // Initialized to false

	// Printing all variables
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Company:", company)
	fmt.Println("Username:", username)
	fmt.Println("Is Active:", isActive)
	fmt.Println("Score:", score)
	fmt.Println("Counter:", counter)
	fmt.Println("Message:", message)
	fmt.Println("Valid:", valid)
}
```

### Constants

Constants are declared using the `const` keyword:

```go
package main

import "fmt"

func main() {
	// Single constant
	const pi = 3.14159

	// Multiple constants
	const (
		statusOK       = 200
		statusNotFound = 404
		appName        = "GoLearner"
	)

	fmt.Println("Pi:", pi)
	fmt.Println("HTTP OK:", statusOK)
	fmt.Println("App:", appName)
}
```

### Basic Data Types

Go has several built-in data types:

#### Numeric Types

```go
package main

import "fmt"

func main() {
	// Integer types
	var intValue int = 42        // Platform dependent (32 or 64 bit)
	var int8Value int8 = 127     // -128 to 127
	var int16Value int16 = 32767 // -32768 to 32767
	var int32Value int32 = 2147483647
	var int64Value int64 = 9223372036854775807

	// Unsigned integer types
	var uintValue uint = 42    // Platform dependent (32 or 64 bit)
	var uint8Value uint8 = 255 // 0 to 255
	var uint16Value uint16 = 65535

	// Alias types
	var byteValue byte = 255 // alias for uint8
	var runeValue rune = 'A' // alias for int32, represents a Unicode code point

	// Floating-point types
	var float32Value float32 = 3.14159265358979323846264338327950288419716939937510
	var float64Value float64 = 3.14159265358979323846264338327950288419716939937510

	// Print with type information
	fmt.Printf("int: %d (type: %T)\n", intValue, intValue)
	fmt.Printf("int8: %d (type: %T)\n", int8Value, int8Value)
	fmt.Printf("float32: %f (type: %T)\n", float32Value, float32Value)
	fmt.Printf("float64: %.10f (type: %T)\n", float64Value, float64Value)
	fmt.Printf("byte: %d (type: %T)\n", byteValue, byteValue)
	fmt.Printf("rune: %c (type: %T)\n", runeValue, runeValue)
}
```

#### Boolean Type

```go
package main

import "fmt"

func main() {
	var isTrue bool = true
	var isFalse bool = false

	// Boolean operators
	andResult := isTrue && isFalse // Logical AND
	orResult := isTrue || isFalse  // Logical OR
	notResult := !isTrue           // Logical NOT

	fmt.Printf("isTrue: %t\n", isTrue)
	fmt.Printf("isFalse: %t\n", isFalse)
	fmt.Printf("AND: %t\n", andResult)
	fmt.Printf("OR: %t\n", orResult)
	fmt.Printf("NOT: %t\n", notResult)
}
```

#### String Type

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// String declaration
	var message string = "Hello, Go!"

	// String operations
	length := len(message)
	upperCase := strings.ToUpper(message)
	hasPrefix := strings.HasPrefix(message, "Hello")

	// String concatenation
	firstName := "Go"
	lastName := "pher"
	fullName := firstName + " " + lastName

	// Multi-line string using backticks
	multiLine := `This is a multi-line string
It can span multiple lines
Without escape characters`

	fmt.Println("Message:", message)
	fmt.Println("Length:", length)
	fmt.Println("Uppercase:", upperCase)
	fmt.Println("Has 'Hello' prefix:", hasPrefix)
	fmt.Println("Full name:", fullName)
	fmt.Println("Multi-line:", multiLine)
}
```

## Type Conversion

### Explicit Type Conversion

Go requires explicit type conversion (no automatic type conversion):

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Numeric type conversion
	var intValue int = 42
	var float64Value float64 = float64(intValue) // int to float64
	var uint8Value uint8 = uint8(intValue)       // int to uint8

	// String to number conversion
	str := "100"
	parsedInt, err := strconv.Atoi(str) // string to int
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	parsedFloat, err := strconv.ParseFloat("3.14", 64) // string to float64
	if err != nil {
		fmt.Println("Error converting string to float:", err)
	}

	// Number to string conversion
	intStr := strconv.Itoa(intValue)                     // int to string
	floatStr := strconv.FormatFloat(3.14159, 'f', 2, 64) // float64 to string with 2 decimal places

	// Display results
	fmt.Printf("intValue: %d (type: %T)\n", intValue, intValue)
	fmt.Printf("float64Value: %f (type: %T)\n", float64Value, float64Value)
	fmt.Printf("uint8Value: %d (type: %T)\n", uint8Value, uint8Value)
	fmt.Printf("parsedInt: %d (type: %T)\n", parsedInt, parsedInt)
	fmt.Printf("parsedFloat: %f (type: %T)\n", parsedFloat, parsedFloat)
	fmt.Printf("intStr: %s (type: %T)\n", intStr, intStr)
	fmt.Printf("floatStr: %s (type: %T)\n", floatStr, floatStr)
}
```

### Using fmt package

Go's `fmt` package provides powerful string formatting capabilities:

```go
package main

import "fmt"

func main() {
	name := "Gopher"
	age := 10
	height := 1.23
	isAwesome := true

	// Basic formatting verbs
	fmt.Printf("Name: %s\n", name)              // %s for strings
	fmt.Printf("Age: %d\n", age)                // %d for integers
	fmt.Printf("Height: %.2f meters\n", height) // %f for floats, .2 specifies decimal places
	fmt.Printf("Is awesome? %t\n", isAwesome)   // %t for booleans

	// Type formatting
	fmt.Printf("Type of name: %T\n", name) // %T shows the type

	// Width and alignment
	fmt.Printf("|%-10s|%10s|\n", "Left", "Right") // Left/right alignment with width

	// Custom formats for different bases
	fmt.Printf("Decimal: %d, Binary: %b, Octal: %o, Hex: %x\n", 42, 42, 42, 42)
}
```

## Common Pitfalls and Strict Rules

Go enforces several rules that new developers should be aware of:

1. **Unused variables** cause compilation errors:
   ```go
   func main() {
       x := 10  // Error: x declared but not used
       fmt.Println("Hello")
   }
   ```

2. **Unused imports** cause compilation errors:
   ```go
   import (
       "fmt"
       "time"  // Error if time package is not used
   )
   ```

3. **No implicit type conversion**:
   ```go
   var a int = 10
   var b float64 = a  // Error: cannot use a (type int) as type float64
   // Correct: var b float64 = float64(a)
   ```

4. **Variable shadowing** can cause unexpected behavior:
   ```go
   x := 10
   if true {
       x := 20  // This creates a new variable x, doesn't modify outer x
       fmt.Println(x)  // Prints 20
   }
   fmt.Println(x)  // Prints 10
   ```

5. **Strict error handling**:
   ```go
   // Ignoring errors is bad practice
   result, err := strconv.Atoi("123")
   if err != nil {
       // Always handle errors!
       fmt.Println("Error:", err)
       return
   }
   // Use result safely knowing it's valid
   ```

## Practice Exercise

### Exercise 1: Hello, Personalized World

Create a program that asks for the user's name and then greets them.

### Exercise 2: Simple Calculator

Create a program that performs basic arithmetic (sum, difference, product, quotient) on two numbers.

### Exercise 3: Type Explorer

Create a program that demonstrates different data types and their properties.