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
