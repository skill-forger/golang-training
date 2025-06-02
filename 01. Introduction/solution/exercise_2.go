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
