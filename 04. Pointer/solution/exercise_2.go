package main

import "fmt"

// Swap exchanges the values at two memory locations
func Swap[T any](a, b *T) {
	*a, *b = *b, *a
}

func main() {
	// Swap integers
	x, y := 5, 10
	fmt.Printf("Before swap: x=%d, y=%d\n", x, y)
	Swap(&x, &y)
	fmt.Printf("After swap: x=%d, y=%d\n", x, y)

	// Swap strings
	first, second := "hello", "world"
	fmt.Printf("Before swap: first=%s, second=%s\n", first, second)
	Swap(&first, &second)
	fmt.Printf("After swap: first=%s, second=%s\n", first, second)

	// Swap custom types
	type Person struct {
		Name string
		Age  int
	}

	alice := Person{Name: "Alice", Age: 30}
	bob := Person{Name: "Bob", Age: 25}

	fmt.Printf("Before swap: alice=%v, bob=%v\n", alice, bob)
	Swap(&alice, &bob)
	fmt.Printf("After swap: alice=%v, bob=%v\n", alice, bob)
}
