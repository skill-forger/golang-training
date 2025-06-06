package orderprocessor

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"golang-training/module-09/exercise-2/cart"
	"golang-training/module-09/exercise-2/inventory"
	"golang-training/module-09/exercise-2/models"
)

// ProcessOrder handles the logic for converting a cart into a completed order.
// It interacts with the inventory and models packages.
func ProcessOrder(c *cart.Cart, productPrices map[string]float64) (*models.Order, error) {
	if c == nil || len(c.GetItems()) == 0 {
		return nil, errors.New("cannot process an empty cart")
	}

	// Generate a unique order ID
	orderID := uuid.New().String()

	// Prepare for order creation
	orderItems := make([]models.Item, 0, len(c.GetItems()))
	totalAmount := 0.0

	// Temporary slice to track stock changes before committing, for rollback in case of error
	// For this simple example, we're not doing a full transaction rollback.
	// In a real system, you'd use database transactions.

	for _, item := range c.GetItems() {
		currentStock := inventory.GetStock(item.ProductID)
		if currentStock < item.Quantity {
			return nil, errors.New(fmt.Sprintf("insufficient stock for %s. Available: %d, Requested: %d", item.ProductID, currentStock, item.Quantity))
		}

		// Simulate deducting stock. If this was a real database, this would be part of a transaction.
		err := inventory.RemoveStock(item.ProductID, item.Quantity)
		if err != nil {
			// This error should ideally not happen if GetStock check passed, but good to handle.
			return nil, fmt.Errorf("failed to deduct stock for %s: %w", item.ProductID, err)
		}

		// Add item to the order and calculate total
		orderItems = append(orderItems, item)
		if price, ok := productPrices[item.ProductID]; ok {
			totalAmount += price * float64(item.Quantity)
		} else {
			return nil, fmt.Errorf("price not found for product %s", item.ProductID)
		}
	}

	order := &models.Order{
		OrderID:     orderID,
		Items:       orderItems,
		TotalAmount: totalAmount,
		Status:      "Completed",
	}

	fmt.Printf("Order %s processed successfully!\n", orderID)
	return order, nil
}
