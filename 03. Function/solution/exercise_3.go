package main

import (
	"errors"
	"fmt"
	"math"
)

// Define operation function type
type Operation func(float64, float64) (float64, error)

// Basic operations
func Add(a, b float64) (float64, error) {
	return a + b, nil
}

func Subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func Multiply(a, b float64) (float64, error) {
	return a * b, nil
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func Power(a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

// Calculator holds operations and provides methods to use them
type Calculator struct {
	operations map[string]Operation
}

// NewCalculator creates a new calculator with standard operations
func NewCalculator() *Calculator {
	calc := &Calculator{
		operations: make(map[string]Operation),
	}

	// Register basic operations
	calc.RegisterOperation("+", Add)
	calc.RegisterOperation("-", Subtract)
	calc.RegisterOperation("*", Multiply)
	calc.RegisterOperation("/", Divide)
	calc.RegisterOperation("^", Power)

	return calc
}

// RegisterOperation adds a new operation to the calculator
func (c *Calculator) RegisterOperation(symbol string, op Operation) {
	c.operations[symbol] = op
}

// Calculate performs the specified operation
func (c *Calculator) Calculate(a, b float64, symbol string) (float64, error) {
	operation, found := c.operations[symbol]
	if !found {
		return 0, fmt.Errorf("unknown operation: %s", symbol)
	}

	return operation(a, b)
}

func main() {
	calc := NewCalculator()

	// Use the calculator
	result, err := calc.Calculate(10, 5, "+")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 + 5 =", result) // Output: 10 + 5 = 15
	}

	// Add a custom operation
	calc.RegisterOperation("avg", func(a, b float64) (float64, error) {
		return (a + b) / 2, nil
	})

	result, _ = calc.Calculate(10, 20, "avg")
	fmt.Println("Average of 10 and 20:", result) // Output: Average of 10 and 20: 15

	// Try division by zero
	result, err = calc.Calculate(10, 0, "/")
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: division by zero
	}
}
