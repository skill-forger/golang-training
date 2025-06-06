package main

import (
	"fmt"

	"golang-training/module-09/exercise-2/cart"
	"golang-training/module-09/exercise-2/inventory"
	"golang-training/module-09/exercise-2/models"
	processor "golang-training/module-09/exercise-2/order"
)

func main() {
	fmt.Println("--- Starting E-commerce Simulation ---")

	// 1. Initialize Product Data (Simulated Database/Catalog)
	products := map[string]models.Product{
		"P001": {ID: "P001", Name: "Laptop Pro", Price: 1200.00},
		"P002": {ID: "P002", Name: "Mechanical Keyboard", Price: 150.00},
		"P003": {ID: "P003", Name: "Wireless Mouse", Price: 50.00},
		"P004": {ID: "P004", Name: "USB-C Hub", Price: 75.00},
	}

	// Extract product prices for easy lookup by other packages
	productPrices := make(map[string]float64)
	for _, p := range products {
		productPrices[p.ID] = p.Price
	}

	// 2. Initialize Inventory Stock
	initialStock := map[string]int{
		"P001": 5,  // 5 Laptops
		"P002": 10, // 10 Keyboards
		"P003": 20, // 20 Mouses
		"P004": 8,  // 8 USB-C Hubs
	}
	inventory.InitializeProducts(initialStock)

	fmt.Println("\n--- First Customer Order ---")
	// 3. Simulate a User's Shopping Journey (Cart 1)
	customerCart1 := cart.NewCart()
	customerCart1.AddItem("P001", 1) // 1 Laptop
	customerCart1.AddItem("P002", 2) // 2 Keyboards
	customerCart1.AddItem("P003", 3) // 3 Mouses
	customerCart1.AddItem("P001", 1) // Add another Laptop (should update quantity)

	fmt.Printf("Cart 1 Items: %v\n", customerCart1.GetItems())
	fmt.Printf("Cart 1 Total: $%.2f\n", customerCart1.CalculateTotal(productPrices))

	fmt.Println("\n--- Processing Order 1 ---")
	fmt.Println("Stock before Order 1:")
	fmt.Printf("P001 (Laptop Pro) stock: %d\n", inventory.GetStock("P001"))
	fmt.Printf("P002 (Mechanical Keyboard) stock: %d\n", inventory.GetStock("P002"))
	fmt.Printf("P003 (Wireless Mouse) stock: %d\n", inventory.GetStock("P003"))

	order1, err := processor.ProcessOrder(customerCart1, productPrices)
	if err != nil {
		fmt.Printf("Error processing Order 1: %v\n", err)
	} else {
		fmt.Println("Order 1 Details:")
		fmt.Printf("  Order ID: %s\n", order1.OrderID)
		fmt.Printf("  Total Amount: $%.2f\n", order1.TotalAmount)
		fmt.Printf("  Status: %s\n", order1.Status)
		fmt.Println("  Items:")
		for _, item := range order1.Items {
			fmt.Printf("    - Product ID: %s, Quantity: %d\n", item.ProductID, item.Quantity)
		}
	}

	fmt.Println("\nStock after Order 1:")
	fmt.Printf("P001 (Laptop Pro) stock: %d\n", inventory.GetStock("P001"))
	fmt.Printf("P002 (Mechanical Keyboard) stock: %d\n", inventory.GetStock("P002"))
	fmt.Printf("P003 (Wireless Mouse) stock: %d\n", inventory.GetStock("P003"))

	fmt.Println("\n--- Second Customer Order (Insufficient Stock Scenario) ---")
	customerCart2 := cart.NewCart()
	customerCart2.AddItem("P001", 5) // Request 5 Laptops (we only have 3 left)
	customerCart2.AddItem("P004", 1) // 1 USB-C Hub

	fmt.Printf("Cart 2 Items: %v\n", customerCart2.GetItems())
	fmt.Printf("Cart 2 Total: $%.2f\n", customerCart2.CalculateTotal(productPrices))

	fmt.Println("\n--- Processing Order 2 ---")
	fmt.Println("Stock before Order 2:")
	fmt.Printf("P001 (Laptop Pro) stock: %d\n", inventory.GetStock("P001"))
	fmt.Printf("P004 (USB-C Hub) stock: %d\n", inventory.GetStock("P004"))

	order2, err := processor.ProcessOrder(customerCart2, productPrices)
	if err != nil {
		fmt.Printf("Error processing Order 2 (expected): %v\n", err)
	} else {
		fmt.Printf("Order 2 processed unexpectedly: %+v\n", order2)
	}

	fmt.Println("\nStock after attempted Order 2:")
	fmt.Printf("P001 (Laptop Pro) stock: %d\n", inventory.GetStock("P001"))
	fmt.Printf("P004 (USB-C Hub) stock: %d\n", inventory.GetStock("P004")) // Should remain unchanged for P004

	fmt.Println("\n--- End of Simulation ---")
}
