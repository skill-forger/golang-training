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
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Second version adds Age field
type UserV2 struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Age       int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Third version adds IsActive field and change Age to nullable
type UserV3 struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Age       *int   `gorm:"default:null"` // Changed to pointer to make nullable
	IsActive  bool   `gorm:"default:true"` // New field
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
