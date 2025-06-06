package models

// Product represents an item available for sale.
type Product struct {
	ID    string
	Name  string
	Price float64
}

// Item represents a specific product with a quantity in a cart or order.
type Item struct {
	ProductID string
	Quantity  int
}
