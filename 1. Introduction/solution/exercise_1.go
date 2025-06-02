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
