## Module 02: Practice Exercises

### Exercise 1: Password Validator
Create a program that validates a password based on the following rules:
- At least 8 characters long
- Contains at least one uppercase letter
- Contains at least one lowercase letter
- Contains at least one digit
- Contains at least one special character (!, @, #, $, %, etc.)

```go
// password_validator.go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

func main() {
    var password string
    
    fmt.Print("Enter a password: ")
    fmt.Scanln(&password)
    
    // Track which requirements are met
    var (
        hasMinLength  bool = len(password) >= 8
        hasUpper      bool
        hasLower      bool
        hasDigit      bool
        hasSpecial    bool
    )
    
    // Special characters list
    specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/"
    
    // Check each character
    for _, char := range password {
        if unicode.IsUpper(char) {
            hasUpper = true
        } else if unicode.IsLower(char) {
            hasLower = true
        } else if unicode.IsDigit(char) {
            hasDigit = true
        } else if strings.ContainsRune(specialChars, char) {
            hasSpecial = true
        }
    }
    
    // Check if all requirements are met
    isValid := hasMinLength && hasUpper && hasLower && hasDigit && hasSpecial
    
    // Print the result
    if isValid {
        fmt.Println("Password is valid!")
    } else {
        fmt.Println("Password is invalid. Please ensure it meets the following criteria:")
        if !hasMinLength {
            fmt.Println("- At least 8 characters long")
        }
        if !hasUpper {
            fmt.Println("- Contains at least one uppercase letter")
        }
        if !hasLower {
            fmt.Println("- Contains at least one lowercase letter")
        }
        if !hasDigit {
            fmt.Println("- Contains at least one digit")
        }
        if !hasSpecial {
            fmt.Println("- Contains at least one special character")
        }
    }
}
```

### Exercise 2: FizzBuzz
Implement the classic FizzBuzz program:
- Print numbers from 1 to n
- For multiples of 3, print "Fizz" instead of the number
- For multiples of 5, print "Buzz" instead of the number
- For multiples of both 3 and 5, print "FizzBuzz"

```go
// fizzbuzz.go
package main

import "fmt"

func main() {
    var n int
    
    fmt.Print("Enter a number: ")
    fmt.Scanln(&n)
    
    for i := 1; i <= n; i++ {
        switch {
        case i%3 == 0 && i%5 == 0:
            fmt.Println("FizzBuzz")
        case i%3 == 0:
            fmt.Println("Fizz")
        case i%5 == 0:
            fmt.Println("Buzz")
        default:
            fmt.Println(i)
        }
    }
}
```

### Exercise 3: Number Guessing Game
Create a more advanced number guessing game where the player has to guess a random number within a specified range, with hints and limited attempts.

```go
// advanced_guessing_game.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    
    // Game configuration
    minNumber := 1
    maxNumber := 100
    maxAttempts := 7
    
    // Generate the secret number
    secretNumber := rand.Intn(maxNumber-minNumber+1) + minNumber
    
    // Game introduction
    fmt.Printf("Welcome to the Number Guessing Game!\n")
    fmt.Printf("I'm thinking of a number between %d and %d.\n", minNumber, maxNumber)
    fmt.Printf("You have %d attempts to guess it.\n\n", maxAttempts)
    
    // Previous guesses
    var previousGuesses []int
    
    // Main game loop
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        // Show attempts remaining
        attemptsLeft := maxAttempts - attempt + 1
        fmt.Printf("Attempt %d/%d: ", attempt, maxAttempts)
        
        // Get the player's guess
        var guess int
        fmt.Scan(&guess)
        
        // Check if the guess is valid
        if guess < minNumber || guess > maxNumber {
            fmt.Printf("Please enter a number between %d and %d.\n", minNumber, maxNumber)
            attempt--  // Don't count this as an attempt
            continue
        }
        
        // Check if the guess was already made
        alreadyGuessed := false
        for _, prevGuess := range previousGuesses {
            if guess == prevGuess {
                alreadyGuessed = true
                break
            }
        }
        
        if alreadyGuessed {
            fmt.Printf("You already guessed %d. Try a different number.\n", guess)
            attempt--  // Don't count this as an attempt
            continue
        }
        
        // Add to previous guesses
        previousGuesses = append(previousGuesses, guess)
        
        // Check the guess
        if guess == secretNumber {
            fmt.Printf("\nCongratulations! You guessed the number %d in %d attempts!\n", 
                secretNumber, attempt)
            return
        } else if guess < secretNumber {
            fmt.Println("Too low!")
        } else {
            fmt.Println("Too high!")
        }
        
        // Show previous guesses
        fmt.Print("Previous guesses: ")
        for i, prevGuess := range previousGuesses {
            if i > 0 {
                fmt.Print(", ")
            }
            fmt.Print(prevGuess)
        }
        fmt.Printf("\nAttempts left: %d\n\n", attemptsLeft-1)
        
        // Game over check
        if attempt == maxAttempts {
            fmt.Printf("Game over! The number was %d.\n", secretNumber)
        }
    }
}
```
