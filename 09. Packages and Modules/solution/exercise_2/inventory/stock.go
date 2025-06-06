package inventory

import (
	"errors"
	"fmt"
)

// stock holds the quantity of each product.
// It's a private variable within the inventory package.
var stock = make(map[string]int)

// InitializeProducts sets up initial stock for given products.
func InitializeProducts(initialStock map[string]int) {
	for productID, quantity := range initialStock {
		stock[productID] = quantity
	}
	fmt.Println("Inventory initialized.")
}

// AddStock increases the quantity of a product in stock.
func AddStock(productID string, quantity int) {
	stock[productID] += quantity
	fmt.Printf("Added %d to %s. New stock: %d\n", quantity, productID, stock[productID])
}

// RemoveStock decreases the quantity of a product in stock.
// Returns an error if there's insufficient stock.
func RemoveStock(productID string, quantity int) error {
	if stock[productID] < quantity {
		return errors.New(fmt.Sprintf("insufficient stock for product %s. Available: %d, Requested: %d", productID, stock[productID], quantity))
	}
	stock[productID] -= quantity
	fmt.Printf("Removed %d from %s. New stock: %d\n", quantity, productID, stock[productID])
	return nil
}

// GetStock returns the current stock quantity for a product.
func GetStock(productID string) int {
	return stock[productID]
}
