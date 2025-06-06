package main

import (
	"fmt"

	"golang-training/module-09/exercise-1/utils"
)

func main() {
	// --- Using string utility functions ---
	originalString := "Hello, Go!"
	reversed := utils.ReverseString(originalString)
	fmt.Printf("Original: \"%s\"\nReversed: \"%s\"\n", originalString, reversed)

	palindrome1 := "madam"
	palindrome2 := "racecar"
	notPalindrome := "golang"

	fmt.Printf("\"%s\" is a palindrome: %t\n", palindrome1, utils.IsPalindrome(palindrome1))
	fmt.Printf("\"%s\" is a palindrome: %t\n", palindrome2, utils.IsPalindrome(palindrome2))
	fmt.Printf("\"%s\" is a palindrome: %t\n", notPalindrome, utils.IsPalindrome(notPalindrome))

	fmt.Println("---")

	// --- Using mathematical utility functions ---
	num := 5
	fmt.Printf("Factorial of %d is: %d\n", num, utils.Factorial(num))

	a, b := 10, 25
	fmt.Printf("Max of %d and %d is: %d\n", a, b, utils.Max(a, b))
	fmt.Printf("Min of %d and %d is: %d\n", a, b, utils.Min(a, b))

	fmt.Println("---")

	numNeg := -3
	fmt.Printf("Factorial of %d is: %d\n", numNeg, utils.Factorial(numNeg)) // Should be 0 based on our implementation
}
