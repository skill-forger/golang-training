## Practical Exercises

### Exercise 1: Custom Stack Implementation

Implement a stack data structure using pointers. A stack is a Last-In-First-Out (LIFO) data structure where elements are added and removed from the same end.

Your implementation should include:
1. A `Node` struct that holds a value of any type (using `interface{}`) and a pointer to the next node
2. A `Stack` struct that tracks the top node and the size of the stack
3. The following stack operations:
   - `Push`: Add a new element to the top of the stack
   - `Pop`: Remove and return the top element from the stack
   - `Peek`: View the top element without removing it
   - `Size`: Return the number of elements in the stack
   - `IsEmpty`: Check if the stack is empty
4. Error handling for operations on an empty stack
5. A demonstration in the `main` function that shows all stack operations

```go
// stack.go
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
```

### Exercise 2: Swap Function

Implement a generic swap function using pointers. This exercise demonstrates how pointers allow you to modify variables passed to functions.

Your implementation should:
1. Create a generic `Swap` function that exchanges the values of two variables of any type
2. Use Go's generics (`[T any]`) to ensure type safety
3. Demonstrate the function with different data types:
   - Swap integers
   - Swap strings
   - Swap custom struct types (e.g., a Person struct with Name and Age fields)
4. Print the values before and after swapping to show the effect

```go
// swap.go
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
```

### Exercise 3: Binary Tree Implementation

Implement a simple binary search tree (BST) using pointers. A binary search tree is a hierarchical data structure where each node has at most two children, with values less than the node to the left and values greater than the node to the right.

Your implementation should include:
1. A `TreeNode` struct that contains an integer value and pointers to left and right child nodes
2. A `BinarySearchTree` struct that maintains a pointer to the root node
3. The following tree operations:
   - `Insert`: Add a new value to the tree while maintaining the BST property
   - `Find`: Check if a value exists in the tree
   - `InOrderTraversal`: Visit all nodes in ascending order and apply a function to each value
4. Helper functions using recursion for tree operations
5. A demonstration in the `main` function that:
   - Creates a tree with several values
   - Prints the values in sorted order
   - Searches for values that exist and don't exist in the tree

```go
// binary_tree.go
package main

import "fmt"

// TreeNode represents a node in a binary search tree
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

// BinarySearchTree manages a binary search tree
type BinarySearchTree struct {
    Root *TreeNode
}

// Insert adds a new value to the tree
func (bst *BinarySearchTree) Insert(value int) {
    if bst.Root == nil {
        bst.Root = &TreeNode{Value: value}
        return
    }
    
    insertRecursive(bst.Root, value)
}

// insertRecursive is a helper function for Insert
func insertRecursive(node *TreeNode, value int) {
    if value < node.Value {
        if node.Left == nil {
            node.Left = &TreeNode{Value: value}
        } else {
            insertRecursive(node.Left, value)
        }
    } else {
        if node.Right == nil {
            node.Right = &TreeNode{Value: value}
        } else {
            insertRecursive(node.Right, value)
        }
    }
}

// Find checks if a value exists in the tree
func (bst *BinarySearchTree) Find(value int) bool {
    return findRecursive(bst.Root, value)
}

// findRecursive is a helper function for Find
func findRecursive(node *TreeNode, value int) bool {
    if node == nil {
        return false
    }
    
    if value == node.Value {
        return true
    }
    
    if value < node.Value {
        return findRecursive(node.Left, value)
    }
    
    return findRecursive(node.Right, value)
}

// InOrderTraversal visits all nodes in order and applies a function
func (bst *BinarySearchTree) InOrderTraversal(visit func(int)) {
    inOrderRecursive(bst.Root, visit)
}

// inOrderRecursive is a helper function for InOrderTraversal
func inOrderRecursive(node *TreeNode, visit func(int)) {
    if node != nil {
        inOrderRecursive(node.Left, visit)
        visit(node.Value)
        inOrderRecursive(node.Right, visit)
    }
}

func main() {
    bst := BinarySearchTree{}
    
    // Insert values
    values := []int{50, 30, 70, 20, 40, 60, 80}
    for _, v := range values {
        bst.Insert(v)
    }
    
    // Print all values in order
    fmt.Println("Tree contents (in-order):")
    bst.InOrderTraversal(func(value int) {
        fmt.Print(value, " ")
    })
    fmt.Println()
    
    // Search for values
    for _, v := range []int{40, 90} {
        if bst.Find(v) {
            fmt.Printf("Value %d found in tree\n", v)
        } else {
            fmt.Printf("Value %d NOT found in tree\n", v)
        }
    }
}
```
