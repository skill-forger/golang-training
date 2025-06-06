package cart

import (
	"golang-training/module-09/exercise-2/models"
)

// Cart represents a user's shopping cart.
type Cart struct {
	Items map[string]models.Item // Using map for easy item lookup/update by ProductID
}

// NewCart creates and returns a new empty Cart.
func NewCart() *Cart {
	return &Cart{
		Items: make(map[string]models.Item),
	}
}

// AddItem adds a product to the cart or updates its quantity if already present.
func (c *Cart) AddItem(productID string, quantity int) {
	if quantity <= 0 {
		return // Do not add zero or negative quantity
	}

	item, exists := c.Items[productID]
	if exists {
		item.Quantity += quantity
	} else {
		item = models.Item{
			ProductID: productID,
			Quantity:  quantity,
		}
	}
	c.Items[productID] = item
}

// RemoveItem removes a product from the cart.
func (c *Cart) RemoveItem(productID string) {
	delete(c.Items, productID)
}

// GetItems returns a slice of items currently in the cart.
func (c *Cart) GetItems() []models.Item {
	itemsSlice := make([]models.Item, 0, len(c.Items))
	for _, item := range c.Items {
		itemsSlice = append(itemsSlice, item)
	}
	return itemsSlice
}

// CalculateTotal calculates the total price of all items in the cart.
// It requires a map of product prices to look up individual item prices.
func (c *Cart) CalculateTotal(productPrices map[string]float64) float64 {
	total := 0.0
	for _, item := range c.Items {
		if price, ok := productPrices[item.ProductID]; ok {
			total += price * float64(item.Quantity)
		}
	}
	return total
}
