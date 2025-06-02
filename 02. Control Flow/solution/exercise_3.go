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
			attempt-- // Don't count this as an attempt
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
			attempt-- // Don't count this as an attempt
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
