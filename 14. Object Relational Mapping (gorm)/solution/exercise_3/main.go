package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"size:100;not null"`
	Email     string  `gorm:"size:100;uniqueIndex;not null"`
	Profile   Profile // Has One relationship
	Orders    []Order // Has Many relationship
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Profile model - belongs to User
type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex"`
	Address   string `gorm:"size:200"`
	Phone     string `gorm:"size:20"`
	BirthDate time.Time
}

// Product model
type Product struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `gorm:"size:100;not null"`
	Description  string        `gorm:"type:text"`
	Price        float64       `gorm:"type:decimal(10,2);not null"`
	Categories   []Category    `gorm:"many2many:product_categories;"`
	OrderDetails []OrderDetail // Has Many relationship
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Category model
type Category struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:50;uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"many2many:product_categories;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Order model - belongs to User
type Order struct {
	ID           uint          `gorm:"primaryKey"`
	UserID       uint          `gorm:"index;not null"`
	User         User          // Belongs To relationship
	OrderDetails []OrderDetail // Has Many relationship
	Status       string        `gorm:"size:20;default:'pending'"`
	TotalAmount  float64       `gorm:"type:decimal(10,2);not null"`
	OrderDate    time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OrderDetail model - belongs to Order and Product
type OrderDetail struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"index;not null"`
	Order     Order   // Belongs To relationship
	ProductID uint    `gorm:"index;not null"`
	Product   Product // Belongs To relationship
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Database represents the application's data layer
type Database struct {
	db *gorm.DB
}

// NewDatabase creates a new database connection and migrates the schema
func NewDatabase() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate all models
	err = db.AutoMigrate(&User{}, &Profile{}, &Product{}, &Category{}, &Order{}, &OrderDetail{})
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// SeedDatabase populates the database with initial data
func (d *Database) SeedDatabase() error {
	// Create categories
	categories := []Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Books", Description: "Books and publications"},
	}

	for i := range categories {
		if err := d.db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// Create products
	products := []Product{
		{
			Name:        "Smartphone",
			Description: "Latest smartphone model",
			Price:       799.99,
		},
		{
			Name:        "Laptop",
			Description: "Professional laptop",
			Price:       1299.99,
		},
		{
			Name:        "T-Shirt",
			Description: "Cotton T-shirt",
			Price:       19.99,
		},
		{
			Name:        "Programming Guide",
			Description: "Comprehensive programming book",
			Price:       49.99,
		},
	}

	for i := range products {
		if err := d.db.Create(&products[i]).Error; err != nil {
			return err
		}
	}

	// Associate products with categories
	if err := d.db.Model(&products[0]).Association("Categories").Append(&categories[0]); err != nil {
		return err
	}
	if err := d.db.Model(&products[1]).Association("Categories").Append(&categories[0]); err != nil {
		return err
	}
	if err := d.db.Model(&products[2]).Association("Categories").Append(&categories[1]); err != nil {
		return err
	}
	if err := d.db.Model(&products[3]).Association("Categories").Append(&categories[2]); err != nil {
		return err
	}

	// Create a user with profile
	user := User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Profile: Profile{
			Address:   "123 Main St",
			Phone:     "555-123-4567",
			BirthDate: time.Date(1990, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	if err := d.db.Create(&user).Error; err != nil {
		return err
	}

	// Create an order for the user
	order := Order{
		UserID:      user.ID,
		Status:      "completed",
		TotalAmount: 2119.97, // Sum of all order items
		OrderDate:   time.Now(),
	}

	if err := d.db.Create(&order).Error; err != nil {
		return err
	}

	// Add order details
	orderDetails := []OrderDetail{
		{
			OrderID:   order.ID,
			ProductID: products[0].ID, // Smartphone
			Quantity:  1,
			Price:     products[0].Price,
		},
		{
			OrderID:   order.ID,
			ProductID: products[1].ID, // Laptop
			Quantity:  1,
			Price:     products[1].Price,
		},
		{
			OrderID:   order.ID,
			ProductID: products[2].ID, // T-Shirt
			Quantity:  1,
			Price:     products[2].Price,
		},
	}

	for i := range orderDetails {
		if err := d.db.Create(&orderDetails[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// QueryExamples demonstrates various GORM queries with relationships
func (d *Database) QueryExamples() {
	// 1. Find user with profile (Preload)
	var user User
	fmt.Println("1. User with Profile:")
	if err := d.db.Preload("Profile").First(&user).Error; err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}
	fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
	fmt.Printf("Profile - Address: %s, Phone: %s\n\n", user.Profile.Address, user.Profile.Phone)

	// 2. Find user with orders (Preload nested relationships)
	fmt.Println("2. User with Orders and Order Details:")
	if err := d.db.Preload("Orders.OrderDetails.Product").First(&user).Error; err != nil {
		log.Fatalf("Failed to fetch user with orders: %v", err)
	}
	fmt.Printf("User: %s has %d orders\n", user.Name, len(user.Orders))
	for i, order := range user.Orders {
		fmt.Printf("Order #%d - Total: $%.2f, Status: %s, Items: %d\n",
			i+1, order.TotalAmount, order.Status, len(order.OrderDetails))
		for j, detail := range order.OrderDetails {
			fmt.Printf("  Item %d: %s, Qty: %d, Price: $%.2f\n",
				j+1, detail.Product.Name, detail.Quantity, detail.Price)
		}
	}
	fmt.Println()

	// 3. Products with their categories
	fmt.Println("3. Products with Categories:")
	var products []Product
	if err := d.db.Preload("Categories").Find(&products).Error; err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}
	for _, product := range products {
		fmt.Printf("Product: %s, Price: $%.2f, Categories: ", product.Name, product.Price)
		for i, category := range product.Categories {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(category.Name)
		}
		fmt.Println()
	}
	fmt.Println()

	// 4. Join query - Find all orders with user info using joins
	fmt.Println("4. Orders with User Info (Join):")
	type OrderWithUser struct {
		OrderID     uint
		Status      string
		TotalAmount float64
		UserName    string
		UserEmail   string
	}

	var ordersWithUsers []OrderWithUser
	if err := d.db.Table("orders").
		Select("orders.id as order_id, orders.status, orders.total_amount, users.name as user_name, users.email as user_email").
		Joins("left join users on users.id = orders.user_id").
		Scan(&ordersWithUsers).Error; err != nil {
		log.Fatalf("Failed to execute join query: %v", err)
	}

	for _, o := range ordersWithUsers {
		fmt.Printf("Order #%d - Status: %s, Total: $%.2f, Customer: %s (%s)\n",
			o.OrderID, o.Status, o.TotalAmount, o.UserName, o.UserEmail)
	}
	fmt.Println()

	// 5. Complex query - Find products in a specific category
	fmt.Println("5. Products in Electronics category:")
	var electronicsProducts []Product
	if err := d.db.Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Joins("JOIN categories ON categories.id = product_categories.category_id").
		Where("categories.name = ?", "Electronics").
		Find(&electronicsProducts).Error; err != nil {
		log.Fatalf("Failed to execute complex query: %v", err)
	}

	for _, p := range electronicsProducts {
		fmt.Printf("Product: %s, Price: $%.2f\n", p.Name, p.Price)
	}

	// 6. Transaction example - Create order with rollback on error
	fmt.Println("\n6. Transaction Example (Commit / Rollback):")

	var txUser User
	if err := d.db.First(&txUser).Error; err != nil {
		log.Fatalf("Failed to fetch user for transaction demo: %v", err)
	}

	var txProduct Product
	if err := d.db.First(&txProduct).Error; err != nil {
		log.Fatalf("Failed to fetch product for transaction demo: %v", err)
	}

	err := d.db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("Starting transaction...")

		// Create order
		order := Order{
			UserID:      txUser.ID,
			Status:      "pending",
			TotalAmount: txProduct.Price,
			OrderDate:   time.Now(),
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// ❌ Intentionally invalid quantity to trigger rollback
		orderDetail := OrderDetail{
			OrderID:   order.ID,
			ProductID: txProduct.ID,
			Quantity:  -1, // invalid quantity
			Price:     txProduct.Price,
		}

		if orderDetail.Quantity <= 0 {
			return fmt.Errorf("invalid quantity, rolling back transaction")
		}

		if err := tx.Create(&orderDetail).Error; err != nil {
			return err
		}

		fmt.Println("Transaction committed")
		return nil
	})

	if err != nil {
		fmt.Println("Transaction rolled back:", err)
	}

	// Verify rollback
	var orderCount int64
	d.db.Model(&Order{}).Where("status = ?", "pending").Count(&orderCount)
	fmt.Printf("Pending orders after transaction: %d\n", orderCount)
}

func main() {
	// Initialize database
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Check if data needs to be seeded
	var count int64
	db.db.Model(&User{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding database with initial data...")
		if err := db.SeedDatabase(); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		fmt.Println("Database seeded successfully!")
	} else {
		fmt.Println("Database already contains data")
	}

	// Run query examples
	db.QueryExamples()
}
