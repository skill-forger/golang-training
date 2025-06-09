# Module 05: Collection Types in Go

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#arrays-fixed-size-sequences">Arrays: Fixed-Size Sequences</a></li>
    <li><a href="#slices-dynamic-and-flexible">Slices: Dynamic and Flexible</a></li>
    <li><a href="#maps-key-value-collections">Maps: Key-Value Collections</a></li>
    <li><a href="#collection-type-comparison">Collection Type Comparison</a></li>
    <li><a href="#common-patterns-and-idioms">Common Patterns and Idioms</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:

- Understand the differences between arrays, slices, and maps
- Master the creation and manipulation of fixed-size arrays
- Learn how to work with dynamic slices and their underlying memory model
- Use maps effectively for key-value storage
- Apply common patterns and idioms for data transformation
- Implement efficient algorithms using Go's collection types
- Recognize which collection type to use for different scenarios

## Overview

Collections are fundamental to nearly all programming tasks,
allowing us to store, organize, and manipulate groups of data efficiently.
Go provides three primary collection types: arrays, slices, and maps,
each with distinct characteristics and use cases.
Understanding these collection types and their behaviors is essential for writing effective Go programs.

## Arrays: Fixed-Size Sequences

Arrays in Go are fixed-length sequences of elements of a single type.
Their size is part of the type declaration and cannot change during execution.

### Basic Array Operations

```go
// array_basics.go
package main

import "fmt"

func main() {
	// Declaration and initialization
	var scores [5]int                    // Declare array of 5 integers (all initialized to 0)
	fmt.Println("Empty scores:", scores) // Output: [0 0 0 0 0]

	scores[0] = 95 // Assign values by index
	scores[1] = 89
	scores[2] = 78
	scores[3] = 92
	scores[4] = 85
	fmt.Println("Filled scores:", scores) // Output: [95 89 78 92 85]

	// Array literals
	names := [3]string{"Alice", "Bob", "Charlie"} // Initialize with values
	fmt.Println("Names:", names)                  // Output: [Alice Bob Charlie]

	// Array with size determined by initializer
	cities := [...]string{"New York", "London", "Tokyo", "Paris", "Beijing"}
	fmt.Println("Cities count:", len(cities)) // Output: 5

	// Accessing elements
	fmt.Println("First city:", cities[0])            // Output: New York
	fmt.Println("Last city:", cities[len(cities)-1]) // Output: Beijing

	// Arrays have fixed size
	// cities[5] = "Sydney"  // This would cause a compile-time error: index out of range
}
```

### Array Iteration Techniques

```go
// array_iteration.go
package main

import "fmt"

func main() {
	temperatures := [7]float64{23.5, 25.1, 24.8, 26.2, 28.0, 27.5, 26.8}

	fmt.Println("Daily temperatures:")

	// Method 1: Traditional for loop with index
	for i := 0; i < len(temperatures); i++ {
		fmt.Printf("Day %d: %.1f°C\n", i+1, temperatures[i])
	}

	// Method 2: Range-based iteration (more idiomatic in Go)
	fmt.Println("\nUsing range:")
	for index, temp := range temperatures {
		fmt.Printf("Day %d: %.1f°C\n", index+1, temp)
	}

	// Method 3: Range with only value (ignoring index)
	fmt.Println("\nAverage calculation:")
	sum := 0.0
	for _, temp := range temperatures {
		sum += temp
	}
	average := sum / float64(len(temperatures))
	fmt.Printf("Average temperature: %.1f°C\n", average)
}
```

### Multi-dimensional Arrays

Go supports multi-dimensional arrays, which are useful for grid-like data structures:

```go
// multi_dimensional.go
package main

import "fmt"

func main() {
	// 2D array - 3 rows, 4 columns
	var grid [3][4]int

	// Initialize with nested loops
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			grid[row][col] = row*4 + col
		}
	}

	// Print as a matrix
	fmt.Println("2D Grid:")
	for _, row := range grid {
		fmt.Println(row)
	}

	// 2D array with initialization
	chessboard := [8][8]string{
		{"r", "n", "b", "q", "k", "b", "n", "r"},
		{"p", "p", "p", "p", "p", "p", "p", "p"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"P", "P", "P", "P", "P", "P", "P", "P"},
		{"R", "N", "B", "Q", "K", "B", "N", "R"},
	}

	// Display a specific position
	fmt.Printf("Piece at e1: %s\n", chessboard[7][4]) // King
}
```

### When to Use Arrays

Arrays in Go are most suitable when:

- The collection size is known and fixed
- Memory efficiency is critical
- You need stack allocation
- You're working with a specific mathematical or algorithmic requirement

## Slices: Dynamic and Flexible

Slices are the most common collection type in Go, providing a flexible, dynamic view into an underlying array.

### Slice Basics

```go
// slice_basics.go
package main

import "fmt"

func main() {
	// Creating slices
	// Method 1: Slice literal
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Println("Numbers:", numbers)

	// Method 2: make function (with length and capacity)
	scores := make([]int, 5, 10) // length 5, capacity 10
	fmt.Printf("Scores: %v, Length: %d, Capacity: %d\n",
		scores, len(scores), cap(scores))

	// Method 3: Slice from an array
	fruits := [5]string{"apple", "banana", "cherry", "date", "elderberry"}
	someFruits := fruits[1:4] // elements 1 through 3
	fmt.Println("Some fruits:", someFruits)

	// Method 4: Empty slice
	var empty []int
	fmt.Printf("Empty: %v, Length: %d, isNil: %t\n",
		empty, len(empty), empty == nil)

	// Zero-length but non-nil slice
	noElements := []int{}
	fmt.Printf("No elements: %v, Length: %d, isNil: %t\n",
		noElements, len(noElements), noElements == nil)

	// Modifying slice elements
	someFruits[0] = "blueberry" // Changes the underlying array
	fmt.Println("Modified fruits array:", fruits)
	fmt.Println("Modified slice:", someFruits)
}
```

### Slice Operations

The real power of slices comes from operations like `append`, slicing, and access to the underlying array:

```go
// slice_operations.go
package main

import "fmt"

func main() {
	// Basic slice
	s := []int{2, 3, 5, 7, 11, 13}
	fmt.Println("Original slice:", s)

	// Slicing a slice
	fmt.Println("s[1:4] =", s[1:4]) // [3 5 7]
	fmt.Println("s[:3] =", s[:3])   // [2 3 5]
	fmt.Println("s[3:] =", s[3:])   // [7 11 13]

	// Appending elements
	s = append(s, 17, 19, 23)
	fmt.Println("After append:", s) // [2 3 5 7 11 13 17 19 23]

	// Appending another slice
	more := []int{29, 31}
	s = append(s, more...)
	fmt.Println("After append slice:", s) // [2 3 5 7 11 13 17 19 23 29 31]

	// Copying slices
	destination := make([]int, len(s))
	copied := copy(destination, s)
	fmt.Printf("Copied %d elements: %v\n", copied, destination)

	// Modifying destination doesn't affect source
	destination[0] = 999
	fmt.Println("Source after modifying copy:", s)
	fmt.Println("Modified copy:", destination)
}
```

### Slice Internals: Length and Capacity

Understanding length and capacity is crucial for effective slice usage:

```go
// slice_length_capacity.go
package main

import "fmt"

func main() {
	// Create a slice with make to demonstrate capacity
	s := make([]int, 3, 8)
	printSliceInfo("s", s)

	// Append elements up to capacity
	s = append(s, 1, 2, 3, 4, 5)
	printSliceInfo("s after append within capacity", s)

	// Append beyond capacity - triggers reallocation
	s = append(s, 6)
	printSliceInfo("s after capacity growth", s)

	// Creating a slice from another slice
	t := s[2:5]
	printSliceInfo("t := s[2:5]", t)

	// The full capacity view of t
	u := s[2:5:8] // slice with capacity limit
	printSliceInfo("u := s[2:5:8]", u)
}

func printSliceInfo(name string, slice []int) {
	fmt.Printf("%s: len=%d cap=%d %v\n",
		name, len(slice), cap(slice), slice)
}
```

### Common Slice Pitfalls

```go
// slice_gotchas.go
package main

import "fmt"

func main() {
	// Pitfall 1: Sharing underlying arrays
	original := []int{1, 2, 3, 4, 5}
	slice1 := original[1:3]
	slice2 := original[2:4]

	fmt.Println("Before modification:")
	fmt.Println("original:", original)
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)

	// Modifying one slice affects others that share the array
	slice1[1] = 999

	fmt.Println("\nAfter modification:")
	fmt.Println("original:", original)
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)

	// Pitfall 2: Append might detach from original
	fmt.Println("\nAppend behavior:")
	slice3 := original[1:3]
	fmt.Println("slice3:", slice3)

	// Append within capacity
	slice3 = append(slice3, 100)
	fmt.Println("After append within capacity:")
	fmt.Println("original:", original) // Original is affected
	fmt.Println("slice3:", slice3)

	// Append beyond capacity
	slice3 = append(slice3, 200, 300, 400)
	fmt.Println("After append beyond capacity:")
	fmt.Println("original:", original) // Original is not affected
	fmt.Println("slice3:", slice3)     // slice3 has new backing array
}
```

## Maps: Key-Value Collections

Maps provide an unordered collection of key-value pairs, with highly efficient lookups, insertions, and deletions.

### Map Basics

```go
// map_basics.go
package main

import "fmt"

func main() {
	// Creating maps
	// Method 1: Map literal
	population := map[string]int{
		"Tokyo":       37400068,
		"Delhi":       31399566,
		"Shanghai":    27058480,
		"São Paulo":   22043028,
		"Mexico City": 21782378,
	}
	fmt.Println("City populations:", population)

	// Method 2: make function
	scores := make(map[string]int)

	// Adding elements
	scores["Alice"] = 95
	scores["Bob"] = 82
	scores["Charlie"] = 88

	fmt.Println("Student scores:", scores)

	// Accessing elements
	fmt.Println("Bob's score:", scores["Bob"])

	// Accessing non-existent key returns zero value
	fmt.Println("David's score (not in map):", scores["David"])

	// Checking if a key exists
	davidsScore, exists := scores["David"]
	if exists {
		fmt.Println("David's score:", davidsScore)
	} else {
		fmt.Println("David is not in the system")
	}

	// Deleting an entry
	delete(scores, "Bob")
	fmt.Println("After deleting Bob:", scores)

	// Map length
	fmt.Println("Number of students:", len(scores))
}
```

### Map Iteration

Unlike in some languages, map iteration order in Go is not guaranteed:

```go
// map_iteration.go
package main

import (
	"fmt"
	"sort"
)

func main() {
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
		"black": "#000000",
		"white": "#FFFFFF",
	}

	// Basic iteration (unordered)
	fmt.Println("Unordered iteration:")
	for color, hex := range colors {
		fmt.Printf("%s: %s\n", color, hex)
	}

	// To iterate in a specific order, sort the keys first
	fmt.Println("\nOrdered iteration:")
	var keys []string
	for key := range colors {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, color := range keys {
		fmt.Printf("%s: %s\n", color, colors[color])
	}

	// Iterating over just keys or values
	fmt.Println("\nJust the colors:")
	for color := range colors {
		fmt.Println(color)
	}

	fmt.Println("\nJust the hex codes:")
	for _, hex := range colors {
		fmt.Println(hex)
	}
}
```

### Maps with Complex Values

Maps in Go can have values of any type, including structs, slices, or even other maps:

```go
// complex_maps.go
package main

import "fmt"

// Define a struct to use as a map value
type Student struct {
	Name  string
	Age   int
	Grade string
}

func main() {
	// Map with struct values
	students := map[int]Student{
		101: {Name: "Alice", Age: 21, Grade: "A"},
		102: {Name: "Bob", Age: 22, Grade: "B+"},
		103: {Name: "Charlie", Age: 20, Grade: "A-"},
	}

	fmt.Println("Student 102:", students[102])

	// Map with slice values
	classRooms := map[string][]string{
		"Science": {"Alice", "Bob", "Eve"},
		"Math":    {"Charlie", "Dave", "Frank"},
		"Art":     {"George", "Helen"},
	}

	fmt.Println("Math class students:", classRooms["Math"])

	// Adding a student to a class
	classRooms["Science"] = append(classRooms["Science"], "Ivan")
	fmt.Println("Updated Science class:", classRooms["Science"])

	// Map of maps
	schoolSystem := map[string]map[string]int{
		"Elementary": {
			"Grade 1": 25,
			"Grade 2": 28,
			"Grade 3": 30,
		},
		"HighSchool": {
			"Grade 9":  120,
			"Grade 10": 115,
			"Grade 11": 110,
			"Grade 12": 105,
		},
	}

	// Access nested map
	fmt.Println("Elementary Grade 2 students:",
		schoolSystem["Elementary"]["Grade 2"])

	// Add a new nested entry
	if _, exists := schoolSystem["MiddleSchool"]; !exists {
		schoolSystem["MiddleSchool"] = map[string]int{
			"Grade 6": 90,
			"Grade 7": 85,
			"Grade 8": 88,
		}
	}

	fmt.Println("Updated school system:", schoolSystem)
}
```

### Map Performance Considerations

Maps in Go are implemented as hash tables, providing:

- Fast lookups: O(1) average case
- Fast insertions and deletions: O(1) average case
- Optimized memory usage
- No guarantee of iteration order

## Collection Type Comparison

When deciding which collection type to use, consider these differences:

| Feature                | Array                | Slice               | Map                  |
|------------------------|----------------------|---------------------|----------------------|
| Size                   | Fixed                | Dynamic             | Dynamic              |
| Indexing               | Integer-only         | Integer-only        | Any comparable type  |
| Memory allocation      | Stack (small arrays) | Heap                | Heap                 |
| Zero value             | Zero-filled array    | nil                 | nil                  |
| Direct comparison      | Yes (==, !=)         | No                  | No                   |
| Pass-by-value behavior | Copies entire array  | Copies slice header | Copies map reference |
| Key lookup             | O(1)                 | O(1)                | O(1) average         |
| Memory overhead        | None                 | Small               | Moderate             |
| Iteration order        | Guaranteed           | Guaranteed          | Not guaranteed       |

## Common Patterns and Idioms

### Filtering Slices

```go
// filter_slice.go
package main

import "fmt"

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter even numbers
	var evenNumbers []int
	for _, num := range numbers {
		if num%2 == 0 {
			evenNumbers = append(evenNumbers, num)
		}
	}

	fmt.Println("Original numbers:", numbers)
	fmt.Println("Even numbers:", evenNumbers)

	// More efficient approach with pre-allocated slice
	// When you know the potential maximum size
	oddNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if num%2 != 0 {
			oddNumbers = append(oddNumbers, num)
		}
	}

	fmt.Println("Odd numbers:", oddNumbers)
}
```

### Transforming Maps

```go
// transform_map.go
package main

import (
	"fmt"
	"strings"
)

func main() {
	// Original map
	scientists := map[string]int{
		"Einstein": 1879,
		"Newton":   1643,
		"Galileo":  1564,
		"Darwin":   1809,
		"Tesla":    1856,
		"Curie":    1867,
	}

	// Transform: Calculate age at 1900
	ageIn1900 := make(map[string]int, len(scientists))
	for name, birthYear := range scientists {
		ageIn1900[name] = 1900 - birthYear
	}

	fmt.Println("Age in 1900:")
	for name, age := range ageIn1900 {
		fmt.Printf("%s: %d years old\n", name, age)
	}

	// Transform: Convert keys to uppercase
	uppercaseKeys := make(map[string]int, len(scientists))
	for name, year := range scientists {
		uppercaseKeys[strings.ToUpper(name)] = year
	}

	fmt.Println("\nUppercase names:")
	for name, year := range uppercaseKeys {
		fmt.Printf("%s: born in %d\n", name, year)
	}
}
```

### Counting with Maps

```go
// counting.go
package main

import (
	"fmt"
	"strings"
)

func main() {
	text := `Go is an open source programming language that makes it easy to build 
    simple, reliable, and efficient software. Go was designed at Google in 2007 
    by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar 
    to C, but with memory safety, garbage collection, structural typing, and 
    CSP-style concurrency.`

	// Count word frequencies
	words := strings.Fields(strings.ToLower(text))
	frequencies := make(map[string]int)

	for _, word := range words {
		// Remove punctuation (simplified approach)
		word = strings.Trim(word, ".,;:!?()[]{}\"'")
		frequencies[word]++
	}

	fmt.Println("Word frequencies:")
	for word, count := range frequencies {
		if count > 1 {
			fmt.Printf("%-12s: %d\n", word, count)
		}
	}
}
```

## Best Practices

1. **Choose the Right Collection Type**
    - Arrays: When size is fixed and known at compile time
    - Slices: For the vast majority of sequence needs
    - Maps: When you need key-value associations
2. **Memory Efficiency**
    - Pre-allocate slices when you know the approximate size: `make([]int, 0, capacity)`
    - Be cautious with very large arrays (they're copied when passed to functions)
    - Remember that `append()` might reallocate memory
3. **Safety First**
    - Always check for existence when accessing map elements with the two-value form
    - Check slice bounds before indexing
    - Watch for nil slices vs. empty slices (`nil` vs. `[]int{}`)
4. **Performance Considerations**
    - Avoid unnecessary allocations and copying
    - Reuse slices when possible
    - Use `copy()` for explicit slice duplication
    - Remember that map operations have small overhead compared to direct array access
5. **Idiomatic Go**
    - Use `range` for iteration
    - Leverage slice expressions for clean subsetting
    - Use maps for lookup tables and counting

## Practice Exercises

### Exercise 1: Student Grade Tracker

Create a program that tracks and analyzes student grades using maps and slices.
This exercise will help you understand how to work with collection types to store and process related data.

Your implementation should:

1. Create a map that associates student names (strings) with their grades (slices of integers)
2. Calculate each student's average grade and store it in a new map
3. Identify the student with the highest average grade
4. Create a ranked list of students based on their average grades
5. Display formatted output showing:
    - Individual student averages
    - The top performing student
    - A ranked list of all students

### Exercise 2: Word Frequency Counter

Develop a text analysis tool that counts word frequencies in a text document.
This exercise demonstrates how to process strings, use regular expressions, and manipulate maps for data analysis.

Your implementation should:

1. Process a sample text (or optionally read from a file)
2. Convert the text to lowercase and extract individual words
3. Count the frequency of each word in the text using a map
4. Filter out common stop words that don't add meaning
5. Sort the words by frequency in descending order
6. Display the top N most frequent words in a formatted table
7. Optionally write the complete results to a file

### Exercise 3: Contact Book Application

Build a contact management application that allows users to store and retrieve contact information.
This exercise combines structs with collection types to create a more complex data management system.

Your implementation should include:

1. A `Contact` struct that stores personal information (first name, last name, email, phone)
2. A `ContactBook` struct that manages a collection of contacts using a map
3. Methods for:
    - Adding a new contact to the book
    - Finding contacts by searching for a name or partial name
    - Deleting contacts
    - Listing all contacts in alphabetical order
    - Grouping contacts by their first letter
4. A demonstration that shows all the functionality of the contact book
5. Proper handling of case sensitivity in searches
6. Sorting capabilities for displaying contacts in a structured way
