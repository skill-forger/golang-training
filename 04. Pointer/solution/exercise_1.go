package main

import (
	"errors"
	"fmt"
)

// Node represents an element in the stack
type Node struct {
	Value interface{}
	Next  *Node
}

// Stack is a basic LIFO stack implementation
type Stack struct {
	top  *Node
	size int
}

// Push adds a new value to the top of the stack
func (s *Stack) Push(value interface{}) {
	s.top = &Node{
		Value: value,
		Next:  s.top,
	}
	s.size++
}

// Pop removes and returns the top value from the stack
func (s *Stack) Pop() (interface{}, error) {
	if s.size == 0 {
		return nil, errors.New("stack is empty")
	}

	value := s.top.Value
	s.top = s.top.Next
	s.size--

	return value, nil
}

// Peek returns the top value without removing it
func (s *Stack) Peek() (interface{}, error) {
	if s.size == 0 {
		return nil, errors.New("stack is empty")
	}

	return s.top.Value, nil
}

// Size returns the number of elements in the stack
func (s *Stack) Size() int {
	return s.size
}

// IsEmpty returns true if the stack is empty
func (s *Stack) IsEmpty() bool {
	return s.size == 0
}

func main() {
	stack := Stack{}

	// Push elements
	stack.Push("first")
	stack.Push("second")
	stack.Push("third")

	fmt.Println("Stack size:", stack.Size())

	// Peek at top element
	top, _ := stack.Peek()
	fmt.Println("Top element:", top)

	// Pop elements
	for !stack.IsEmpty() {
		value, _ := stack.Pop()
		fmt.Println("Popped:", value)
	}

	// Try to pop from empty stack
	_, err := stack.Pop()
	fmt.Println("Error:", err)
}
