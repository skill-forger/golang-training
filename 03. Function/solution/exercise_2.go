package main

import (
	"fmt"
	"strings"
)

// Function types
type StringProcessor func(string) string

type IntProcessor func(int) int

// ComposeStringProcessors two string processors into a single processor
func ComposeStringProcessors(f, g StringProcessor) StringProcessor {
	return func(s string) string {
		return g(f(s))
	}
}

// Chain applies a series of string processors in sequence
func Chain(processors ...StringProcessor) StringProcessor {
	return func(s string) string {
		result := s
		for _, processor := range processors {
			result = processor(result)
		}
		return result
	}
}

func main() {
	// Define some string processors
	trim := func(s string) string { return strings.TrimSpace(s) }
	upper := func(s string) string { return strings.ToUpper(s) }
	reverse := func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	// Compose processors
	trimAndUpper := ComposeStringProcessors(trim, upper)
	fmt.Println(trimAndUpper("  hello world  ")) // Output: HELLO WORLD

	// Chain multiple processors
	processAll := Chain(trim, upper, reverse)
	fmt.Println(processAll("  hello world  ")) // Output: DLROW OLLEH

	// Create and apply custom chains
	emphasize := Chain(
		trim,
		upper,
		func(s string) string { return "*** " + s + " ***" },
	)
	fmt.Println(emphasize("  important message  ")) // Output: *** IMPORTANT MESSAGE ***
}
