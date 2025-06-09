package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	. "golang-training/module-09/exercise-3/math_rand"
	strUtil "golang-training/module-09/exercise-3/math_string"
	"golang-training/module-09/exercise-3/math_utils"
)

func main() {
	// --- Using string utility functions ---
	originalString := "Hello, Go!"
	reversed := strUtil.ReverseString(originalString)
	fmt.Printf("Original: \"%s\"\nReversed: \"%s\"\n", originalString, reversed)

	palindrome1 := "madam"
	palindrome2 := "racecar"
	notPalindrome := "golang"

	fmt.Printf("\"%s\" is a palindrome: %t\n", palindrome1, strUtil.IsPalindrome(palindrome1))
	fmt.Printf("\"%s\" is a palindrome: %t\n", palindrome2, strUtil.IsPalindrome(palindrome2))
	fmt.Printf("\"%s\" is a palindrome: %t\n", notPalindrome, strUtil.IsPalindrome(notPalindrome))

	fmt.Println("---")

	// --- Using mathematical utility functions ---
	num := 5
	fmt.Printf("Factorial of %d is: %d\n", num, math_utils.Factorial(num))

	a, b := 10, 25
	fmt.Printf("Max of %d and %d is: %d\n", a, b, math_utils.Max(a, b))
	fmt.Printf("Min of %d and %d is: %d\n", a, b, math_utils.Min(a, b))

	fmt.Println("---", RandomUuid())

	numNeg := -3
	fmt.Printf("Factorial of %d is: %d\n", numNeg, math_utils.Factorial(numNeg)) // Should be 0 based on our implementation
}
