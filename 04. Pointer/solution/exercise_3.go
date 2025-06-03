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
