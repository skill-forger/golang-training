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
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:100;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Stock       int     `gorm:"default:0"`
	Category    string  `gorm:"size:50;index"`
	IsActive    bool    `gorm:"default:true"`
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
