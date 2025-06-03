# Module 02: Control Flow

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#conditional-statements">Conditional Statements</a></li>
    <li><a href="#loops">Loops</a></li>
    <li><a href="#switch">Switch</a></li>
    <li><a href="#break-and-continue">Break and Continue</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#common-pitfallsaand-mistakes">Common Pitfalls and Mistakes</a></li>
    <li><a href="#practice-exercise">Practice Exercise</a></li>
</ol>

## Objectives

- Master Go's conditional statements (if, if-else, nested if)
- Understand Go's unique loop structure and its various forms
- Use switch statements efficiently for multi-way branching
- Apply break and continue statements to control loop execution
- Implement labeled control flow statements for complex scenarios
- Recognize common control flow patterns and best practices

## Overview

Control flow structures determine the order in which statements are executed in a program. 
Go provides clean, efficient control structures that allow you to:
- Make decisions (conditional statements)
- Repeat operations (loops)
- Direct program execution based on multiple conditions (switch statements)
- Break out of or continue loops (break/continue)

Understanding control flow is essential for writing dynamic, 
responsive programs that can make decisions and adapt to different inputs or conditions.

## Conditional Statements (if-else)

Conditional statements allow your program to make decisions based on certain conditions.

### Basic Syntax
```go
// Basic if statement
if condition {
    // Code executed when condition is true
}

// if-else statement
if condition {
    // Code executed when condition is true
} else {
    // Code executed when condition is false
}

// if-else if-else chain
if condition1 {
    // Code executed when condition1 is true
} else if condition2 {
    // Code executed when condition1 is false but condition2 is true
} else {
    // Code executed when both condition1 and condition2 are false
}
```

### Key Features of Go's if Statements
- **No Parentheses Required**: Unlike many languages, Go doesn't require parentheses around conditions
   ```go
   // Go style (correct)
   if x > 10 {
       fmt.Println("x is greater than 10")
   }
   
   // C/Java style (unnecessary in Go)
   if (x > 10) {
       fmt.Println("x is greater than 10")
   }
   ```
- **Braces Are Required**: Unlike some languages, the braces `{}` are mandatory even for single-statement blocks
- **Statement Before Condition**: Go allows you to execute a short statement before the condition
   ```go
   if value := calculateValue(); value > threshold {
       fmt.Println("Value exceeds threshold:", value)
   }
   // Note: 'value' is only accessible within this if block
   ```
- **No Implicit Type Conversion**: Go requires boolean conditions - it won't interpret non-zero values as true

### Example: Simple Grade Calculator
```go
// grade_calculator.go
package main

import "fmt"

func main() {
    var score int
    
    fmt.Print("Enter your score (0-100): ")
    fmt.Scanln(&score)
    
    // Simple validation
    if score < 0 || score > 100 {
        fmt.Println("Invalid score. Please enter a value between 0 and 100.")
        return
    }
    
    // Grade calculation
    if score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    } else if score >= 70 {
        fmt.Println("Grade: C")
    } else if score >= 60 {
        fmt.Println("Grade: D")
    } else {
        fmt.Println("Grade: F")
    }
    
    // Additional feedback
    if score >= 60 {
        fmt.Println("You passed!")
    } else {
        fmt.Println("You need to study more.")
    }
}
```

### Example: Statement Before Condition
```go
// login_example.go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Get username
    var username string
    fmt.Print("Enter username: ")
    fmt.Scanln(&username)
    
    // If statement with initialization
    if normalizedName := strings.ToLower(strings.TrimSpace(username)); normalizedName == "admin" {
        fmt.Println("Welcome administrator!")
    } else if len(normalizedName) < 3 {
        fmt.Println("Username too short, must be at least 3 characters.")
    } else {
        fmt.Printf("Welcome user %s!\n", normalizedName)
    }
    
    // Note: normalizedName is not accessible here
    // fmt.Println(normalizedName) // This would cause a compilation error
}
```

## Loops

Go simplifies looping by providing a single, flexible `for` loop construct that can be used in different ways.

### Basic Syntax
```go
// 1. Traditional C-style for loop
for initialization; condition; post {
    // Code to repeat
}

// 2. While-like loop
for condition {
    // Code to repeat while condition is true
}

// 3. Infinite loop
for {
    // Code to repeat indefinitely (until break)
}

// 4. For-each loop with range
for index, value := range collection {
    // Code to execute for each element
}
```

### Key Features of Go's Loops
- **Single Loop Construct**: Unlike many other languages, Go has only one loop keyword: `for`
- **Multiple Loop Forms**: Despite having only one loop keyword, it handles all looping scenarios
- **Range-based Iteration**: Go provides an elegant way to iterate over arrays, slices, maps, and strings using `range`
- **No Do-While Loop**: Go doesn't have a do-while loop, but you can emulate it with a for loop and a break statement

### Example: Traditional For Loop
```go
// countdown.go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Counting down:")
    
    // Traditional for loop with initialization, condition, and post statement
    for i := 10; i > 0; i-- {
        fmt.Println(i)
        time.Sleep(500 * time.Millisecond) // Pause for half a second
    }
    
    fmt.Println("Liftoff!")
}
```

### Example: While-Like Loop
```go
// sum_digits.go
package main

import "fmt"

func main() {
    var number, sum int
    
    fmt.Print("Enter a positive number: ")
    fmt.Scanln(&number)
    
    // While-like loop
    for number > 0 {
        digit := number % 10  // Get the last digit
        sum += digit          // Add to sum
        number /= 10          // Remove the last digit
    }
    
    fmt.Println("Sum of digits:", sum)
}
```

### Example: Infinite Loop with Break
```go
// guess_number.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    
    // Generate a random number between 1 and 100
    secretNumber := rand.Intn(100) + 1
    
    fmt.Println("I've selected a random number between 1 and 100.")
    fmt.Println("Can you guess it?")
    
    attempts := 0
    
    // Infinite loop
    for {
        attempts++
        
        var guess int
        fmt.Print("Enter your guess: ")
        fmt.Scanln(&guess)
        
        if guess < secretNumber {
            fmt.Println("Too low! Try a higher number.")
        } else if guess > secretNumber {
            fmt.Println("Too high! Try a lower number.")
        } else {
            fmt.Printf("Congratulations! You guessed the number in %d attempts!\n", attempts)
            break // Exit the loop when the guess is correct
        }
    }
}
```

### Example: For-Range Loop
```go
// word_analyzer.go
package main

import (
    "fmt"
    "strings"
)

func main() {
    var text string
    
    fmt.Print("Enter a sentence: ")
    fmt.Scanln(&text)
    
    // Convert to lowercase for easier analysis
    text = strings.ToLower(text)
    
    // Count vowels
    vowels := 0
    consonants := 0
    
    // Loop through each character using range
    for _, char := range text {
        if char >= 'a' && char <= 'z' {
            if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
                vowels++
            } else {
                consonants++
            }
        }
    }
    
    fmt.Printf("Your text contains %d vowels and %d consonants.\n", vowels, consonants)
    
    // Demonstrate range with slice
    words := strings.Fields(text) // Split into words
    fmt.Println("\nWords in your sentence:")
    
    for index, word := range words {
        fmt.Printf("%d: %s\n", index+1, word)
    }
}
```

## Switch

Switch statements provide a cleaner way to express multiple conditions compared to long if-else chains.

### Basic Syntax
```go
// 1. Expression switch
switch expression {
case value1:
    // Code for value1
case value2, value3:
    // Code for value2 or value3
default:
    // Code for other values
}

// 2. Type switch (for interfaces)
switch x.(type) {
case type1:
    // Code for type1
case type2:
    // Code for type2
default:
    // Code for other types
}

// 3. Condition switch (like if-else)
switch {
case condition1:
    // Code for condition1
case condition2:
    // Code for condition2
default:
    // Code for other conditions
}
```

### Key Features of Go's Switch Statement
- **No Fall-Through**: Unlike C and similar languages, Go doesn't automatically fall through to the next case
- **Optional Expression**: The switch expression is optional, allowing for condition-based cases
- **Multiple Values Per Case**: A case can handle multiple values, separated by commas
- **Type Switching**: Go supports switching based on the type of an interface value
- **Explicit Fallthrough**: The `fallthrough` keyword can be used when fall-through behavior is desired

### Example: Expression Switch
```go
// day_type.go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Get the current time
    now := time.Now()
    
    // Get the day of the week (Sunday = 0, Monday = 1, etc.)
    weekday := now.Weekday()
    
    fmt.Printf("Today is %s. ", weekday)
    
    // Switch on the weekday
    switch weekday {
    case time.Saturday, time.Sunday:
        fmt.Println("It's the weekend!")
    case time.Friday:
        fmt.Println("It's almost the weekend.")
    case time.Monday:
        fmt.Println("It's the start of the work week.")
    default:
        fmt.Println("It's a regular weekday.")
    }
}
```

### Example: Type Switch
```go
// type_checker.go
package main

import "fmt"

func describeValue(value interface{}) string {
    // Switch on the type of the interface value
    switch v := value.(type) {
    case nil:
        return "nil value"
    case int:
        return fmt.Sprintf("integer with value %d", v)
    case float64:
        return fmt.Sprintf("float with value %f", v)
    case bool:
        if v {
            return "boolean true"
        }
        return "boolean false"
    case string:
        return fmt.Sprintf("string with value '%s' and length %d", v, len(v))
    case []int:
        return fmt.Sprintf("integer slice with %d elements", len(v))
    case map[string]int:
        return fmt.Sprintf("string->int map with %d pairs", len(v))
    default:
        return fmt.Sprintf("unknown type: %T", v)
    }
}

func main() {
    values := []interface{}{
        42,
        3.14,
        true,
        "hello",
        []int{1, 2, 3},
        map[string]int{"one": 1, "two": 2},
        struct{ name string }{"Go"},
    }
    
    for _, value := range values {
        description := describeValue(value)
        fmt.Printf("Value: %v is %s\n", value, description)
    }
}
```

### Example: Condition Switch
```go
// score_feedback.go
package main

import "fmt"

func main() {
    var score int
    
    fmt.Print("Enter your score (0-100): ")
    fmt.Scanln(&score)
    
    // Using switch with conditions instead of if-else chain
    switch {
    case score < 0 || score > 100:
        fmt.Println("Invalid score. Please enter a value between 0 and 100.")
    case score >= 90:
        fmt.Println("Excellent! You got an A.")
    case score >= 80:
        fmt.Println("Good job! You got a B.")
    case score >= 70:
        fmt.Println("Not bad. You got a C.")
    case score >= 60:
        fmt.Println("You passed with a D.")
    default:
        fmt.Println("You need to study more. You got an F.")
    }
}
```

### Example: Fallthrough
```go
// fallthrough_demo.go
package main

import "fmt"

func main() {
    var level int
    
    fmt.Print("Enter your access level (1-3): ")
    fmt.Scanln(&level)
    
    fmt.Println("Your access permissions:")
    
    switch level {
    case 3:
        fmt.Println("- Administrator access")
        fallthrough
    case 2:
        fmt.Println("- File editing privileges")
        fallthrough
    case 1:
        fmt.Println("- Read-only access")
    default:
        fmt.Println("- No access")
    }
}
```

## Break and Continue

`break` and `continue` statements provide additional control within loops.

### Break
The `break` statement terminates the current loop or switch statement and transfers control to the statement following the terminated statement.

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // Exit the loop when i equals 5
    }
    fmt.Println(i)
}
// Output: 0 1 2 3 4
```

### Continue
The `continue` statement skips the rest of the current iteration and continues with the next iteration of the loop.

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // Skip even numbers
    }
    fmt.Println(i)
}
// Output: 1 3 5 7 9
```

### Labels with Break and Continue
For nested loops, Go allows you to specify which loop to break from or continue to:

```go
OuterLoop:
    for i := 0; i < 5; i++ {
        for j := 0; j < 5; j++ {
            if i*j > 10 {
                fmt.Println("Breaking out of outer loop")
                break OuterLoop
            }
            fmt.Printf("i=%d, j=%d\n", i, j)
        }
    }
```

### Example: Finding Prime Numbers
```go
// find_primes.go
package main

import "fmt"

func main() {
    var limit int
    
    fmt.Print("Find prime numbers up to: ")
    fmt.Scanln(&limit)
    
    if limit < 2 {
        fmt.Println("There are no prime numbers less than 2.")
        return
    }
    
    fmt.Printf("Prime numbers up to %d:\n", limit)
    
    // Check each number from 2 to the limit
    for num := 2; num <= limit; num++ {
        isPrime := true
        
        // Check if num is divisible by any number from 2 to sqrt(num)
        for divisor := 2; divisor*divisor <= num; divisor++ {
            if num%divisor == 0 {
                isPrime = false
                break  // Exit inner loop once we find a divisor
            }
        }
        
        if isPrime {
            fmt.Print(num, " ")
        }
    }
    fmt.Println()
}
```

### Example: Labeled Break
```go
// nested_loop_example.go
package main

import "fmt"

func main() {
    fmt.Println("Multiplication table (up to products of 50):")
    
    // Label for the outer loop
OuterLoop:
    for i := 1; i <= 10; i++ {
        for j := 1; j <= 10; j++ {
            product := i * j
            
            // If product exceeds 50, break out of both loops
            if product > 50 {
                fmt.Printf("\nStopping at i=%d, j=%d because %d × %d = %d exceeds 50\n", 
                    i, j, i, j, product)
                break OuterLoop
            }
            
            fmt.Printf("%3d ", product)
        }
        fmt.Println()  // New line after each row
    }
}
```

### Example: Skipping Iterations
```go
// skip_multiples.go
package main

import "fmt"

func main() {
    fmt.Println("Numbers from 1 to 20, skipping multiples of 3:")
    
    for i := 1; i <= 20; i++ {
        // Skip multiples of 3
        if i%3 == 0 {
            continue
        }
        fmt.Print(i, " ")
    }
    fmt.Println()
}
```

## Best Practices

1. **Keep it Simple**
   - Don't nest too many control structures
   - Extract complex logic into separate functions
   - Consider early returns to reduce nesting
2. **Prefer switch to long if-else chains**
   - More readable and maintainable
   - Better performance in many cases
   - Clearer intent
3. **Variable Scoping**
   - Use the short statement form `if x := getValue(); x > 10 { ... }` to limit variable scope
   - Variables defined in if statements are only accessible within that block
4. **Infinite Loops**
   - Use `for { ... }` for intentional infinite loops
   - Always ensure there's a way to exit (break, return, os.Exit)
5. **Error Handling**
   - Check for errors immediately after they might occur
   - Use early returns for error cases to avoid nesting
   ```go
   // Good error handling
   file, err := os.Open(filename)
   if err != nil {
       return err  // Return early on error
   }
   // Continue with the file...
   ```
6. **Break & Continue**
   - Use sparingly and with clear intent
   - Consider refactoring if you have many break/continue statements
   - Always document labeled breaks for clarity

## Common Pitfalls and Mistakes

1. **Forgetting Braces**
   - Unlike some languages, Go always requires braces for control structures
2. **Off-by-One Errors**
   - Be careful with loop boundaries, especially when iterating through arrays or slices
3. **Unintended Shadowing**
   ```go
   x := 10
   if x := 5; x > 0 {  // This creates a new 'x' variable
       fmt.Println(x)  // Prints 5
   }
   fmt.Println(x)      // Prints 10
   ```
4. **Misplaced Else Clauses**
   - In Go, the `else` must be on the same line as the closing brace of the if statement
5. **Infinity Loops Without Exit Conditions**
   - Always ensure a way to break out of infinite loops
6. **Misusing fallthrough**
   - Remember that fallthrough passes control to the next case unconditionally

## Practice Exercise

### Exercise 1: Password Validator
Create a program that validates a password based on the following rules:
- At least 8 characters long
- Contains at least one uppercase letter
- Contains at least one lowercase letter
- Contains at least one digit
- Contains at least one special character (!, @, #, $, %, etc.)

### Exercise 2: FizzBuzz
Implement the classic FizzBuzz program:
- Print numbers from 1 to n
- For multiples of 3, print "Fizz" instead of the number
- For multiples of 5, print "Buzz" instead of the number
- For multiples of both 3 and 5, print "FizzBuzz"

### Exercise 3: Number Guessing Game
Create a more advanced number guessing game where the player has to guess a random number within a specified range, 
with hints and limited attempts.
