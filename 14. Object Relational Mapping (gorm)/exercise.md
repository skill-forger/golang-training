## Practical Exercises

### Exercise 1: Basic CRUD Operations with GORM

Create a program that implements basic CRUD operations using GORM:

```go
// gorm_crud.go
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Product model
type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"default:0"`
	Category    string    `gorm:"size:50;index"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProductService handles database operations for products
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a new product service with the provided database connection
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// Create adds a new product to the database
func (s *ProductService) Create(product *Product) error {
	return s.db.Create(product).Error
}

// FindByID retrieves a product by its ID
func (s *ProductService) FindByID(id uint) (*Product, error) {
	var product Product
	err := s.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindAll retrieves all products
func (s *ProductService) FindAll() ([]Product, error) {
	var products []Product
	err := s.db.Find(&products).Error
	return products, err
}

// FindByCategory retrieves products by category
func (s *ProductService) FindByCategory(category string) ([]Product, error) {
	var products []Product
	err := s.db.Where("category = ?", category).Find(&products).Error
	return products, err
}

// Update modifies an existing product
func (s *ProductService) Update(product *Product) error {
	return s.db.Save(product).Error
}

// Delete removes a product by ID
func (s *ProductService) Delete(id uint) error {
	return s.db.Delete(&Product{}, id).Error
}

// SearchProducts searches for products by name or description
func (s *ProductService) SearchProducts(query string) ([]Product, error) {
	var products []Product
	err := s.db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products).Error
	return products, err
}

func main() {
	// Set up the logger for GORM
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	// Connect to a SQLite database
	db, err := gorm.Open(sqlite.Open("products.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create a product service
	productService := NewProductService(db)

	// Create products
	products := []Product{
		{
			Name:        "Laptop",
			Description: "High-performance laptop with 16GB RAM",
			Price:       1299.99,
			Stock:       10,
			Category:    "Electronics",
			IsActive:    true,
		},
		{
			Name:        "Smartphone",
			Description: "Latest smartphone with advanced camera",
			Price:       799.99,
			Stock:       15,
			Category:    "Electronics",
		},
		{
			Name:        "Coffee Maker",
			Description: "Automatic coffee maker with timer",
			Price:       89.99,
			Stock:       5,
			Category:    "Home Appliances",
		},
	}

	// Demonstrate CRUD operations
	fmt.Println("--- Create Products ---")
	for _, product := range products {
		if err := productService.Create(&product); err != nil {
			log.Printf("Failed to create product: %v", err)
		} else {
			fmt.Printf("Product created: %s (ID: %d)\n", product.Name, product.ID)
		}
	}

	fmt.Println("\n--- Find All Products ---")
	allProducts, err := productService.FindAll()
	if err != nil {
		log.Printf("Failed to retrieve products: %v", err)
	} else {
		for _, p := range allProducts {
			fmt.Printf("ID: %d, Name: %s, Price: $%.2f, Category: %s\n",
				p.ID, p.Name, p.Price, p.Category)
		}
	}

	fmt.Println("\n--- Find Products by Category ---")
	electronicsProducts, err := productService.FindByCategory("Electronics")
	if err != nil {
		log.Printf("Failed to retrieve electronics products: %v", err)
	} else {
		fmt.Printf("Found %d products in Electronics category:\n", len(electronicsProducts))
		for _, p := range electronicsProducts {
			fmt.Printf("ID: %d, Name: %s, Price: $%.2f\n", p.ID, p.Name, p.Price)
		}
	}

	fmt.Println("\n--- Update Product ---")
	if len(allProducts) > 0 {
		productToUpdate := allProducts[0]
		// Increase price by 10%
		productToUpdate.Price = productToUpdate.Price * 1.1
		if err := productService.Update(&productToUpdate); err != nil {
			log.Printf("Failed to update product: %v", err)
		} else {
			fmt.Printf("Product updated: %s (New price: $%.2f)\n",
				productToUpdate.Name, productToUpdate.Price)
		}
	}

	fmt.Println("\n--- Search Products ---")
	searchResults, err := productService.SearchProducts("coffee")
	if err != nil {
		log.Printf("Search failed: %v", err)
	} else {
		fmt.Printf("Found %d products matching 'coffee':\n", len(searchResults))
		for _, p := range searchResults {
			fmt.Printf("ID: %d, Name: %s, Description: %s\n",
				p.ID, p.Name, p.Description)
		}
	}

	fmt.Println("\n--- Delete Product ---")
	if len(allProducts) > 0 {
		idToDelete := allProducts[len(allProducts)-1].ID
		if err := productService.Delete(idToDelete); err != nil {
			log.Printf("Failed to delete product: %v", err)
		} else {
			fmt.Printf("Product with ID %d deleted\n", idToDelete)
		}
	}

	fmt.Println("\n--- Final Product List ---")
	finalProducts, _ := productService.FindAll()
	for _, p := range finalProducts {
		fmt.Printf("ID: %d, Name: %s, Price: $%.2f\n", p.ID, p.Name, p.Price)
	}
}
```

### Exercise 2: Relationships with GORM

Implement a database schema with relationships using GORM:

```go
// gorm_relationships.go
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
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Profile   Profile   // Has One relationship
	Orders    []Order   // Has Many relationship
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
		fmt.Println("Database seeded successfully!\n")
	} else {
		fmt.Println("Database already contains data\n")
	}

	// Run query examples
	db.QueryExamples()
}
```

### Exercise 3: GORM Migration and Schema Management

Create a program to demonstrate database migration and schema changes:

```go
// gorm_migration.go
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initial version of the User model
type UserV1 struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Second version adds Age field
type UserV2 struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Age       int       `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Third version adds IsActive field and change Age to nullable
type UserV3 struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"size:100;not null"`
	Email     string     `gorm:"size:100;uniqueIndex;not null"`
	Age       *int       `gorm:"default:null"` // Changed to pointer to make nullable
	IsActive  bool       `gorm:"default:true"` // New field
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Add soft delete
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *gorm.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager() (*MigrationManager, error) {
	db, err := gorm.Open(sqlite.Open("migration_demo.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MigrationManager{db: db}, nil
}

// RunMigrations demonstrates migration steps
func (m *MigrationManager) RunMigrations() {
	// Step 1: Initial schema
	fmt.Println("Step 1: Creating initial schema (UserV1)")
	if err := m.db.AutoMigrate(&UserV1{}); err != nil {
		log.Fatalf("Failed to migrate to UserV1: %v", err)
	}

	// Add some initial users
	initialUsers := []UserV1{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}

	for _, user := range initialUsers {
		if err := m.db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
		}
	}

	// Display users after initial migration
	var usersV1 []UserV1
	m.db.Find(&usersV1)
	fmt.Println("Users after initial migration:")
	for _, u := range usersV1 {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
	}
	fmt.Println()

	// Step 2: Migrate to the second version (add Age field)
	fmt.Println("Step 2: Migrating to UserV2 (adding Age field)")
	if err := m.db.AutoMigrate(&UserV2{}); err != nil {
		log.Fatalf("Failed to migrate to UserV2: %v", err)
	}

	// Update existing users with age
	m.db.Exec("UPDATE users SET age = ? WHERE name = ?", 30, "Alice")
	m.db.Exec("UPDATE users SET age = ? WHERE name = ?", 25, "Bob")

	// Add a new user with the age field
	newUserV2 := UserV2{
		Name:  "Charlie",
		Email: "charlie@example.com",
		Age:   35,
	}
	m.db.Create(&newUserV2)

	// Display users after second migration
	var usersV2 []UserV2
	m.db.Find(&usersV2)
	fmt.Println("Users after adding Age field:")
	for _, u := range usersV2 {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d\n", u.ID, u.Name, u.Email, u.Age)
	}
	fmt.Println()

	// Step 3: Migrate to the third version (make Age nullable and add IsActive)
	fmt.Println("Step 3: Migrating to UserV3 (nullable Age, add IsActive, add soft delete)")
	if err := m.db.AutoMigrate(&UserV3{}); err != nil {
		log.Fatalf("Failed to migrate to UserV3: %v", err)
	}

	// Set IsActive for existing users
	m.db.Exec("UPDATE users SET is_active = ? WHERE name IN (?, ?)", true, "Alice", "Bob")
	m.db.Exec("UPDATE users SET is_active = ? WHERE name = ?", false, "Charlie")

	// Remove age for one user to demonstrate NULL value
	m.db.Exec("UPDATE users SET age = NULL WHERE name = ?", "Bob")

	// Add a new user with the full schema
	age := 40
	newUserV3 := UserV3{
		Name:     "Dave",
		Email:    "dave@example.com",
		Age:      &age,
		IsActive: true,
	}
	m.db.Create(&newUserV3)

	// Display users after third migration
	var usersV3 []UserV3
	m.db.Find(&usersV3)
	fmt.Println("Users after final migration:")
	for _, u := range usersV3 {
		ageStr := "NULL"
		if u.Age != nil {
			ageStr = fmt.Sprintf("%d", *u.Age)
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %s, IsActive: %v\n",
			u.ID, u.Name, u.Email, ageStr, u.IsActive)
	}
	fmt.Println()

	// Demonstrate soft delete
	fmt.Println("Soft deleting user 'Charlie'")
	m.db.Where("name = ?", "Charlie").Delete(&UserV3{})

	// Show all users including soft deleted
	var allUsers []UserV3
	m.db.Unscoped().Find(&allUsers)
	fmt.Println("All users (including soft deleted):")
	for _, u := range allUsers {
		ageStr := "NULL"
		if u.Age != nil {
			ageStr = fmt.Sprintf("%d", *u.Age)
		}
		deletedStr := "Active"
		if u.DeletedAt != nil {
			deletedStr = "Deleted"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %s, Status: %s\n",
			u.ID, u.Name, u.Email, ageStr, deletedStr)
	}
}

func main() {
	// Initialize migration manager
	mgr, err := NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to initialize migration manager: %v", err)
	}

	// Run migrations
	mgr.RunMigrations()
}
```
