package models

type Order struct {
	OrderID     string
	Items       []Item
	TotalAmount float64
	Status      string // e.g., "Pending", "Completed", "Cancelled"
}
