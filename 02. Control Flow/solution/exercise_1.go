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
		hasMinLength bool = len(password) >= 8
		hasUpper     bool
		hasLower     bool
		hasDigit     bool
		hasSpecial   bool
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
