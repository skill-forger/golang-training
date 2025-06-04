package main

import (
	"fmt"
	"time"
)

// Product represents an item in the inventory
type Product struct {
	SKU          string
	Name         string
	Description  string
	Category     string
	Price        float64
	Cost         float64
	StockLevel   int
	ReorderLevel int
	Supplier     string
	DateAdded    time.Time
}

// GetProfit returns the profit margin for a product
func (p Product) GetProfit() float64 {
	return p.Price - p.Cost
}

// GetProfitMargin returns the profit margin percentage
func (p Product) GetProfitMargin() float64 {
	if p.Price == 0 {
		return 0
	}
	return (p.GetProfit() / p.Price) * 100
}

// NeedsReorder checks if a product needs to be reordered
func (p Product) NeedsReorder() bool {
	return p.StockLevel <= p.ReorderLevel
}

// StockValue returns the total value of this product in stock
func (p Product) StockValue() float64 {
	return float64(p.StockLevel) * p.Cost
}

// Transaction represents an inventory transaction
type Transaction struct {
	ID         string
	ProductSKU string
	Type       string // "purchase", "sale", "adjustment"
	Quantity   int
	Date       time.Time
	Reference  string // invoice or order number
}

// Inventory manages the product catalog and transactions
type Inventory struct {
	Products     map[string]*Product
	Transactions []Transaction
}

// NewInventory creates a new inventory system
func NewInventory() *Inventory {
	return &Inventory{
		Products:     make(map[string]*Product),
		Transactions: []Transaction{},
	}
}

// AddProduct adds a new product to the inventory
func (i *Inventory) AddProduct(p Product) error {
	if _, exists := i.Products[p.SKU]; exists {
		return fmt.Errorf("product with SKU %s already exists", p.SKU)
	}

	p.DateAdded = time.Now()
	i.Products[p.SKU] = &p
	return nil
}

// RecordPurchase records a product purchase
func (i *Inventory) RecordPurchase(sku string, quantity int, reference string) error {
	product, exists := i.Products[sku]
	if !exists {
		return fmt.Errorf("product with SKU %s not found", sku)
	}

	// Update stock level
	product.StockLevel += quantity

	// Record transaction
	transaction := Transaction{
		ID:         fmt.Sprintf("T%d", len(i.Transactions)+1),
		ProductSKU: sku,
		Type:       "purchase",
		Quantity:   quantity,
		Date:       time.Now(),
		Reference:  reference,
	}

	i.Transactions = append(i.Transactions, transaction)
	return nil
}

// RecordSale records a product sale
func (i *Inventory) RecordSale(sku string, quantity int, reference string) error {
	product, exists := i.Products[sku]
	if !exists {
		return fmt.Errorf("product with SKU %s not found", sku)
	}

	if product.StockLevel < quantity {
		return fmt.Errorf("insufficient stock: have %d, need %d", product.StockLevel, quantity)
	}

	// Update stock level
	product.StockLevel -= quantity

	// Record transaction
	transaction := Transaction{
		ID:         fmt.Sprintf("T%d", len(i.Transactions)+1),
		ProductSKU: sku,
		Type:       "sale",
		Quantity:   quantity,
		Date:       time.Now(),
		Reference:  reference,
	}

	i.Transactions = append(i.Transactions, transaction)
	return nil
}

// AdjustStock adjusts stock level (e.g., for inventory count)
func (i *Inventory) AdjustStock(sku string, newLevel int, reason string) error {
	product, exists := i.Products[sku]
	if !exists {
		return fmt.Errorf("product with SKU %s not found", sku)
	}

	// Calculate adjustment amount
	adjustment := newLevel - product.StockLevel

	// Update stock level
	product.StockLevel = newLevel

	// Record transaction
	transaction := Transaction{
		ID:         fmt.Sprintf("T%d", len(i.Transactions)+1),
		ProductSKU: sku,
		Type:       "adjustment",
		Quantity:   adjustment,
		Date:       time.Now(),
		Reference:  reason,
	}

	i.Transactions = append(i.Transactions, transaction)
	return nil
}

// GetLowStockProducts returns all products that need reordering
func (i *Inventory) GetLowStockProducts() []*Product {
	var lowStock []*Product

	for _, product := range i.Products {
		if product.NeedsReorder() {
			lowStock = append(lowStock, product)
		}
	}

	return lowStock
}

// GetInventoryValue returns the total value of inventory
func (i *Inventory) GetInventoryValue() float64 {
	var total float64

	for _, product := range i.Products {
		total += product.StockValue()
	}

	return total
}

// GetProductTransactions returns all transactions for a specific product
func (i *Inventory) GetProductTransactions(sku string) []Transaction {
	var transactions []Transaction

	for _, t := range i.Transactions {
		if t.ProductSKU == sku {
			transactions = append(transactions, t)
		}
	}

	return transactions
}

func main() {
	// Create a new inventory
	inventory := NewInventory()

	// Add products
	products := []Product{
		{
			SKU:          "LAPTOP001",
			Name:         "Pro Laptop 15\"",
			Description:  "High-performance laptop with 16GB RAM",
			Category:     "Electronics",
			Price:        1299.99,
			Cost:         950.00,
			StockLevel:   0,
			ReorderLevel: 5,
			Supplier:     "TechSuppliers Inc.",
		},
		{
			SKU:          "PHONE001",
			Name:         "Smartphone X",
			Description:  "Latest model smartphone",
			Category:     "Electronics",
			Price:        799.99,
			Cost:         550.00,
			StockLevel:   0,
			ReorderLevel: 10,
			Supplier:     "MobileTech Inc.",
		},
		{
			SKU:          "CHAIR001",
			Name:         "Ergonomic Office Chair",
			Description:  "Adjustable office chair with lumbar support",
			Category:     "Furniture",
			Price:        249.99,
			Cost:         125.00,
			StockLevel:   0,
			ReorderLevel: 3,
			Supplier:     "Office Furnishings Co.",
		},
	}

	for _, product := range products {
		err := inventory.AddProduct(product)
		if err != nil {
			fmt.Printf("Error adding product: %s\n", err)
		}
	}

	// Perform inventory operations
	fmt.Println("Inventory Operations:")

	// Record purchases
	fmt.Println("\n1. Recording purchases:")
	purchases := map[string]int{
		"LAPTOP001": 10,
		"PHONE001":  20,
		"CHAIR001":  8,
	}

	for sku, quantity := range purchases {
		err := inventory.RecordPurchase(sku, quantity, "PO-12345")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			product := inventory.Products[sku]
			fmt.Printf("- Purchased %d units of %s (New stock: %d)\n", quantity, product.Name, product.StockLevel)
		}
	}

	// Record sales
	fmt.Println("\n2. Recording sales:")
	sales := map[string]int{
		"LAPTOP001": 3,
		"PHONE001":  7,
		"CHAIR001":  2,
	}

	for sku, quantity := range sales {
		err := inventory.RecordSale(sku, quantity, "SO-67890")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			product := inventory.Products[sku]
			fmt.Printf("- Sold %d units of %s (New stock: %d)\n", quantity, product.Name, product.StockLevel)
		}
	}

	// Check for low stock
	fmt.Println("\n3. Low stock report:")
	lowStock := inventory.GetLowStockProducts()

	if len(lowStock) == 0 {
		fmt.Println("No products need reordering.")
	} else {
		for _, product := range lowStock {
			fmt.Printf("- %s: Current stock: %d, Reorder level: %d\n", product.Name, product.StockLevel, product.ReorderLevel)
		}
	}

	// Display inventory value
	fmt.Printf("\n4. Total inventory value: $%.2f\n", inventory.GetInventoryValue())

	// Display product profitability
	fmt.Println("\n5. Product profitability:")
	for _, product := range inventory.Products {
		fmt.Printf("- %s: Cost: $%.2f, Price: $%.2f, Margin: %.1f%%\n",
			product.Name, product.Cost, product.Price, product.GetProfitMargin())
	}

	// Adjust stock (e.g., after inventory count)
	fmt.Println("\n6. Stock adjustment:")
	err := inventory.AdjustStock("LAPTOP001", 8, "Inventory count adjustment")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		product := inventory.Products["LAPTOP001"]
		fmt.Printf("- Adjusted %s stock to %d units\n", product.Name, product.StockLevel)
	}

	// Display transaction history for a product
	fmt.Println("\n7. Transaction history for Pro Laptop 15\":")
	transactions := inventory.GetProductTransactions("LAPTOP001")
	for _, t := range transactions {
		fmt.Printf("- %s: %s %d units on %s (Ref: %s)\n",
			t.ID, t.Type, t.Quantity, t.Date.Format("2006-01-02"), t.Reference)
	}
}
