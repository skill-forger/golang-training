## Practical Exercises

### Exercise 1: Hello, Personalized World
Create a program that asks for the user's name and then greets them.

```go
// hello_user.go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    // Create a reader for reading from standard input
    reader := bufio.NewReader(os.Stdin)
    
    // Prompt the user for their name
    fmt.Print("Please enter your name: ")
    
    // Read the input until newline
    name, _ := reader.ReadString('\n')
    
    // Trim whitespace and newlines from the input
    name = strings.TrimSpace(name)
    
    // Greet the user
    fmt.Printf("Hello, %s! Welcome to Go programming!\n", name)
}
```

### Exercise 2: Simple Calculator
Create a program that performs basic arithmetic on two numbers.

```go
// simple_calculator.go
package main

import (
    "fmt"
)

func main() {
    // Declare variables
    var num1, num2 float64
    
    // Get input from user
    fmt.Print("Enter first number: ")
    fmt.Scanln(&num1)
    
    fmt.Print("Enter second number: ")
    fmt.Scanln(&num2)
    
    // Perform calculations
    sum := num1 + num2
    difference := num1 - num2
    product := num1 * num2
    
    // Handle division by zero
    var quotient float64
    if num2 != 0 {
        quotient = num1 / num2
    }
    
    // Display results
    fmt.Printf("Sum: %.2f\n", sum)
    fmt.Printf("Difference: %.2f\n", difference)
    fmt.Printf("Product: %.2f\n", product)
    
    if num2 != 0 {
        fmt.Printf("Quotient: %.2f\n", quotient)
    } else {
        fmt.Println("Cannot divide by zero")
    }
}
```

### Exercise 3: Type Explorer
Create a program that demonstrates different data types and their properties.

```go
// type_explorer.go
package main

import (
    "fmt"
    "math"
    "unsafe"
)

func main() {
    // Declare variables of different types
    var (
        intVar    int     = 42
        floatVar  float64 = 3.14159
        boolVar   bool    = true
        stringVar string  = "Hello, Go!"
        runeVar   rune    = 'A'
        byteVar   byte    = 255
    )
    
    // Display variable values and types
    fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Integer", intVar, intVar, unsafe.Sizeof(intVar))
    fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Float", floatVar, floatVar, unsafe.Sizeof(floatVar))
    fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Boolean", boolVar, boolVar, unsafe.Sizeof(boolVar))
    fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "String", stringVar, stringVar, unsafe.Sizeof(stringVar))
    fmt.Printf("%-10s: %v (%c)\t(Type: %T, Size: %d bytes)\n", "Rune", runeVar, runeVar, runeVar, unsafe.Sizeof(runeVar))
    fmt.Printf("%-10s: %v\t(Type: %T, Size: %d bytes)\n", "Byte", byteVar, byteVar, unsafe.Sizeof(byteVar))
    
    // Show limits of different numeric types
    fmt.Println("\n--- Numeric Type Limits ---")
    fmt.Printf("int8    : %d to %d\n", math.MinInt8, math.MaxInt8)
    fmt.Printf("uint8   : 0 to %d\n", math.MaxUint8)
    fmt.Printf("int16   : %d to %d\n", math.MinInt16, math.MaxInt16)
    fmt.Printf("uint16  : 0 to %d\n", math.MaxUint16)
    fmt.Printf("int32   : %d to %d\n", math.MinInt32, math.MaxInt32)
    fmt.Printf("uint32  : 0 to %d\n", math.MaxUint32)
    fmt.Printf("int64   : %d to %d\n", math.MinInt64, math.MaxInt64)
    // Note: MaxUint64 doesn't fit in a signed int64, so we'd need special handling
}
```
