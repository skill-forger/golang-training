# Module 11: GORM in Go - Working with Databases

## Introduction to ORM and GORM

Object-Relational Mapping (ORM) provides a more intuitive, object-oriented approach to database interactions in Go applications. GORM is one of the most popular ORM libraries in the Go ecosystem, offering a powerful and developer-friendly way to work with relational databases while maintaining Go's simplicity and performance.

### The ORM Landscape in Go

Before diving into GORM specifics, let's understand why an ORM is valuable:

1. **Raw SQL vs ORM**
    - **Development Speed**: ORMs automate repetitive CRUD operations
    - **Code Organization**: Database interactions through Go structs
    - **Safety**: Reduced risk of SQL injection and type-related errors

2. **Why GORM Matters for Modern Applications**
    - Intuitive API with chainable method calls
    - Excellent support for migrations and schema changes
    - Comprehensive feature set including hooks, transactions, and relationships
    - Support for multiple databases (MySQL, PostgreSQL, SQLite, SQL Server)

### Getting Started with GORM

To begin working with GORM, you'll need to install it first:

```go
// Install GORM and the MySQL driver
// In your terminal:
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql

// Basic GORM example with MySQL
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// Define a model
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:255;not null"`
    Email     string `gorm:"size:255;uniqueIndex;not null"`
    Age       uint8
    CreatedAt time.Time
    UpdatedAt time.Time
}

func main() {
    // Connection string
    dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    
    // Open connection to database
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Auto migrate the schema
    db.AutoMigrate(&User{})
    
    // Create a user
    user := User{Name: "John Doe", Email: "john@example.com", Age: 25}
    result := db.Create(&user)
    
    if result.Error != nil {
        log.Fatalf("Failed to create user: %v", result.Error)
    }
    
    fmt.Printf("Created user with ID: %d\n", user.ID)
}
```

#### GORM Core Components
- **DB**: The central component that manages database connections
- **Model**: Defines the structure and behavior of database tables
- **Session**: Provides a context for database operations
- **Hooks**: Functions that run before/after specific operations

### Defining Models in GORM

GORM uses Go structs to define database tables:

```go
// Basic model definition
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:200;not null"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"type:decimal(10,2);not null"`
    Stock       int       `gorm:"default:0;not null"`
    CategoryID  uint      `gorm:"index"` // Foreign key
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"` // Soft delete support
}

// Model with composite primary key
type OrderItem struct {
    OrderID   uint    `gorm:"primaryKey"`
    ProductID uint    `gorm:"primaryKey"`
    Quantity  int     `gorm:"not null"`
    Price     float64 `gorm:"type:decimal(10,2);not null"`
}

// Model with embedded struct
type Customer struct {
    gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
    Name      string  `gorm:"size:100;not null"`
    Email     string  `gorm:"size:100;uniqueIndex;not null"`
    // Address details
    Address   string  `gorm:"size:255"`
    City      string  `gorm:"size:100"`
    State     string  `gorm:"size:100"`
    ZipCode   string  `gorm:"size:20"`
}
```

#### Field Tags and Modifiers

GORM uses struct tags to customize how fields are mapped to database columns:

```go
type BlogPost struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"` // Explicit primary key
    Title     string    `gorm:"size:200;not null;index"`  // Indexed field
    Content   string    `gorm:"type:longtext"`            // Specific column type
    Views     int64     `gorm:"default:0"`                // Default value
    Status    string    `gorm:"size:20;check:status IN ('draft', 'published', 'archived')"` // Check constraint
    AuthorID  uint      `gorm:"not null"`                 // Required field
    Tags      string    `gorm:"-"`                        // Ignore this field in database
    CreatedAt time.Time `gorm:"index:idx_created_updated,priority:1"` // Composite index
    UpdatedAt time.Time `gorm:"index:idx_created_updated,priority:2"` // Composite index
}
```

### Establishing a Database Connection

Proper database connection setup is crucial for GORM applications:

```go
package database

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// DB instance
var DB *gorm.DB

// ConnectDB connects to the database and initializes the global DB variable
func ConnectDB() {
    // Get database credentials from environment variables
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    
    // Construct DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        username, password, host, port, dbName)
    
    // Custom logger configuration
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second,  // Slow SQL threshold
            LogLevel:                  logger.Info,  // Log level
            IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error
            Colorful:                  true,         // Enable color
        },
    )
    
    // Open connection
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time {
            return time.Now().UTC() // Use UTC time for CreatedAt/UpdatedAt
        },
    })
    
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Get generic database object SQL DB to use its functions
    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("Failed to get database: %v", err)
    }
    
    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)           // Maximum idle connections
    sqlDB.SetMaxOpenConns(100)          // Maximum open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection
    
    log.Println("Database connection established")
}
```

To use this in your main application:

```go
package main

import (
    "myapp/database"
    "myapp/models"
)

func main() {
    // Set environment variables or load from .env file
    // ...
    
    // Connect to database
    database.ConnectDB()
    
    // Auto-migrate models
    database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
    
    // Start application
    // ...
}
```

### Basic CRUD Operations

GORM provides intuitive methods for Create, Read, Update, and Delete operations:

```go
// User model for examples
type User struct {
    ID      uint   `gorm:"primaryKey"`
    Name    string `gorm:"not null"`
    Email   string `gorm:"uniqueIndex;not null"`
    Age     int
    Active  bool   `gorm:"default:true"`
}

// ==================== CREATE ====================

// Create a single record
func CreateUser(db *gorm.DB) {
    user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
    
    result := db.Create(&user)
    if result.Error != nil {
        log.Fatalf("Error creating user: %v", result.Error)
    }
    
    fmt.Printf("User created with ID: %d\n", user.ID)
    fmt.Printf("Rows affected: %d\n", result.RowsAffected)
}

// Create multiple records in a batch
func CreateBatchUsers(db *gorm.DB) {
    users := []User{
        {Name: "Bob", Email: "bob@example.com", Age: 25},
        {Name: "Charlie", Email: "charlie@example.com", Age: 35},
        {Name: "Dave", Email: "dave@example.com", Age: 28},
    }
    
    result := db.Create(&users)
    if result.Error != nil {
        log.Fatalf("Error creating users: %v", result.Error)
    }
    
    for _, user := range users {
        fmt.Printf("Created user with ID: %d\n", user.ID)
    }
}

// ==================== READ ====================

// Find a record by primary key
func GetUserByID(db *gorm.DB, id uint) {
    var user User
    result := db.First(&user, id)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            fmt.Printf("User with ID %d not found\n", id)
        } else {
            log.Fatalf("Error retrieving user: %v", result.Error)
        }
        return
    }
    
    fmt.Printf("Found user: %+v\n", user)
}

// Get first record that matches condition
func GetUserByEmail(db *gorm.DB, email string) {
    var user User
    result := db.Where("email = ?", email).First(&user)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            fmt.Printf("User with email %s not found\n", email)
        } else {
            log.Fatalf("Error retrieving user: %v", result.Error)
        }
        return
    }
    
    fmt.Printf("Found user: %+v\n", user)
}

// Get all records
func GetAllUsers(db *gorm.DB) {
    var users []User
    result := db.Find(&users)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving users: %v", result.Error)
    }
    
    fmt.Printf("Found %d users\n", result.RowsAffected)
    for _, user := range users {
        fmt.Printf("User: %+v\n", user)
    }
}

// Get users with condition
func GetActiveUsers(db *gorm.DB) {
    var users []User
    result := db.Where("active = ?", true).Find(&users)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving users: %v", result.Error)
    }
    
    fmt.Printf("Found %d active users\n", result.RowsAffected)
}

// ==================== UPDATE ====================

// Update a single field
func UpdateUserName(db *gorm.DB, id uint, newName string) {
    result := db.Model(&User{}).Where("id = ?", id).Update("name", newName)
    
    if result.Error != nil {
        log.Fatalf("Error updating user: %v", result.Error)
    }
    
    fmt.Printf("Updated %d user(s)\n", result.RowsAffected)
}

// Update multiple fields
func UpdateUser(db *gorm.DB, id uint) {
    result := db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
        "name":   "Updated Name",
        "email":  "updated@example.com",
        "active": false,
    })
    
    if result.Error != nil {
        log.Fatalf("Error updating user: %v", result.Error)
    }
    
    fmt.Printf("Updated %d user(s)\n", result.RowsAffected)
}

// Update using struct (only non-zero fields)
func UpdateUserStruct(db *gorm.DB, id uint) {
    user := User{
        Name:  "New Name",
        Email: "new@example.com",
        // Age is not provided, so it won't be updated
    }
    
    result := db.Model(&User{ID: id}).Updates(user)
    
    if result.Error != nil {
        log.Fatalf("Error updating user: %v", result.Error)
    }
    
    fmt.Printf("Updated %d user(s)\n", result.RowsAffected)
}

// ==================== DELETE ====================

// Delete a record
func DeleteUser(db *gorm.DB, id uint) {
    result := db.Delete(&User{}, id)
    
    if result.Error != nil {
        log.Fatalf("Error deleting user: %v", result.Error)
    }
    
    fmt.Printf("Deleted %d user(s)\n", result.RowsAffected)
}

// Delete with condition
func DeleteInactiveUsers(db *gorm.DB) {
    result := db.Where("active = ?", false).Delete(&User{})
    
    if result.Error != nil {
        log.Fatalf("Error deleting users: %v", result.Error)
    }
    
    fmt.Printf("Deleted %d inactive user(s)\n", result.RowsAffected)
}
```

### Advanced Queries

GORM provides powerful query capabilities:

```go
// ==================== ADVANCED QUERIES ====================

// Select specific columns
func GetUserNames(db *gorm.DB) {
    type UserName struct {
        ID   uint
        Name string
    }
    
    var userNames []UserName
    result := db.Model(&User{}).Select("id", "name").Find(&userNames)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving user names: %v", result.Error)
    }
    
    for _, un := range userNames {
        fmt.Printf("User ID: %d, Name: %s\n", un.ID, un.Name)
    }
}

// Order results
func GetUsersOrderedByAge(db *gorm.DB) {
    var users []User
    result := db.Order("age desc").Find(&users)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving users: %v", result.Error)
    }
    
    for _, user := range users {
        fmt.Printf("User: %s, Age: %d\n", user.Name, user.Age)
    }
}

// Limit and offset (pagination)
func GetUsersPaginated(db *gorm.DB, page, pageSize int) {
    var users []User
    offset := (page - 1) * pageSize
    
    result := db.Offset(offset).Limit(pageSize).Find(&users)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving users: %v", result.Error)
    }
    
    fmt.Printf("Page %d (size %d): Found %d users\n", page, pageSize, len(users))
    for _, user := range users {
        fmt.Printf("User: %+v\n", user)
    }
}

// Group by and having
func GetUsersByAgeGroup(db *gorm.DB) {
    type AgeGroup struct {
        Age   int
        Count int
    }
    
    var results []AgeGroup
    result := db.Model(&User{}).
        Select("age, count(*) as count").
        Group("age").
        Having("count > ?", 1).
        Find(&results)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving age groups: %v", result.Error)
    }
    
    for _, ag := range results {
        fmt.Printf("Age: %d, Count: %d\n", ag.Age, ag.Count)
    }
}

// Joins
func GetUsersWithOrders(db *gorm.DB) {
    type UserOrder struct {
        UserID    uint
        UserName  string
        OrderID   uint
        OrderDate time.Time
        Total     float64
    }
    
    var userOrders []UserOrder
    result := db.Table("users").
        Select("users.id as user_id, users.name as user_name, orders.id as order_id, orders.order_date, orders.total").
        Joins("inner join orders on users.id = orders.user_id").
        Where("orders.total > ?", 100).
        Find(&userOrders)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving user orders: %v", result.Error)
    }
    
    for _, uo := range userOrders {
        fmt.Printf("User: %s, Order ID: %d, Total: $%.2f\n", uo.UserName, uo.OrderID, uo.Total)
    }
}

// Subqueries
func GetUsersWithHighPriceOrders(db *gorm.DB) {
    var users []User
    
    subQuery := db.Table("orders").
        Select("user_id").
        Where("total > ?", 1000).
        Group("user_id")
    
    result := db.Where("id IN (?)", subQuery).Find(&users)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving users: %v", result.Error)
    }
    
    fmt.Printf("Users with high-price orders: %+v\n", users)
}

// Raw SQL
func ExecuteRawSQL(db *gorm.DB) {
    var users []User
    result := db.Raw("SELECT * FROM users WHERE age > ? AND active = ?", 25, true).Scan(&users)
    
    if result.Error != nil {
        log.Fatalf("Error executing raw SQL: %v", result.Error)
    }
    
    fmt.Printf("Found %d users\n", result.RowsAffected)
}
```

### Relationships and Associations

GORM supports various relationship types between models:

```go
// ==================== MODEL RELATIONSHIPS ====================

// One-to-One relationship
type Profile struct {
    ID       uint   `gorm:"primaryKey"`
    Bio      string `gorm:"type:text"`
    UserID   uint   `gorm:"uniqueIndex"` // Foreign key
    User     User   // Belongs to User
}

type User struct {
    ID      uint     `gorm:"primaryKey"`
    Name    string   `gorm:"not null"`
    Email   string   `gorm:"uniqueIndex;not null"`
    Profile Profile  // Has one Profile
}

// One-to-Many relationship
type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"not null"`
    Content   string    `gorm:"type:text"`
    UserID    uint      `gorm:"index"` // Foreign key
    User      User      // Belongs to User
    CreatedAt time.Time
}

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Email string `gorm:"uniqueIndex;not null"`
    Posts []Post // Has many Posts
}

// Many-to-Many relationship
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Roles []Role `gorm:"many2many:user_roles;"` // Many-to-many relationship
}

type Role struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"uniqueIndex;not null"`
    Users []User `gorm:"many2many:user_roles;"` // Many-to-many relationship
}

// Working with associations
func CreateUserWithProfile(db *gorm.DB) {
    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Profile: Profile{
            Bio: "Software engineer with 5 years of experience",
        },
    }
    
    result := db.Create(&user) // Will create both user and profile
    
    if result.Error != nil {
        log.Fatalf("Error creating user with profile: %v", result.Error)
    }
    
    fmt.Printf("Created user with ID: %d and profile ID: %d\n", user.ID, user.Profile.ID)
}

// Retrieving associations (preloading)
func GetUserWithPosts(db *gorm.DB, userID uint) {
    var user User
    result := db.Preload("Posts").First(&user, userID)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving user with posts: %v", result.Error)
    }
    
    fmt.Printf("User: %s, Number of posts: %d\n", user.Name, len(user.Posts))
    for _, post := range user.Posts {
        fmt.Printf("Post ID: %d, Title: %s\n", post.ID, post.Title)
    }
}

// Nested preloading
func GetPostsWithUserAndComments(db *gorm.DB) {
    var posts []Post
    result := db.Preload("User").Preload("Comments.User").Find(&posts)
    
    if result.Error != nil {
        log.Fatalf("Error retrieving posts: %v", result.Error)
    }
    
    for _, post := range posts {
        fmt.Printf("Post: %s by %s\n", post.Title, post.User.Name)
        for _, comment := range post.Comments {
            fmt.Printf("  Comment by %s: %s\n", comment.User.Name, comment.Content)
        }
    }
}

// Creating with associations
func AddRolesToUser(db *gorm.DB, userID uint) {
    // Get user
    var user User
    db.First(&user, userID)
    
    // Create roles
    roles := []Role{
        {Name: "Admin"},
        {Name: "Editor"},
    }
    
    // Add roles to user
    result := db.Model(&user).Association("Roles").Append(&roles)
    
    if result != nil {
        log.Fatalf("Error adding roles to user: %v", result)
    }
    
    fmt.Printf("Added roles to user %s\n", user.Name)
}
```

### Transactions

GORM provides transaction support to ensure data integrity:

```go
// Simple transaction
func CreateUserAndProfile(db *gorm.DB) {
    // Start a transaction
    tx := db.Begin()
    
    // Check for errors in transaction creation
    if tx.Error != nil {
        log.Fatalf("Error starting transaction: %v", tx.Error)
    }
    
    // Perform operations in the transaction
    user := User{Name: "Transaction User", Email: "tx@example.com"}
    if err := tx.Create(&user).Error; err != nil {
        // Rollback the transaction on error
        tx.Rollback()
        log.Fatalf("Error creating user in transaction: %v", err)
    }
    
    profile := Profile{UserID: user.ID, Bio: "Created in a transaction"}
    if err := tx.Create(&profile).Error; err != nil {
        // Rollback the transaction on error
        tx.Rollback()
        log.Fatalf("Error creating profile in transaction: %v", err)
    }
    
    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        log.Fatalf("Error committing transaction: %v", err)
    }
    
    fmt.Println("Transaction completed successfully")
}

// Transaction with callback
func TransferFunds(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) {
    err := db.Transaction(func(tx *gorm.DB) error {
        // Get sender account
        var fromAccount Account
        if err := tx.Clauses(gorm.Clause{
            Name: "FOR UPDATE", // Lock the row
        }).First(&fromAccount, fromAccountID).Error; err != nil {
            return err
        }
        
        // Check balance
        if fromAccount.Balance < amount {
            return errors.New("insufficient funds")
        }
        
        // Get recipient account
        var toAccount Account
        if err := tx.Clauses(gorm.Clause{
            Name: "FOR UPDATE", // Lock the row
        }).First(&toAccount, toAccountID).Error; err != nil {
            return err
        }
        
        // Update sender balance
        if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance - amount).Error; err != nil {
            return err
        }
        
        // Update recipient balance
        if err := tx.Model(&toAccount).Update("balance", toAccount.Balance + amount).Error; err != nil {
            return err
        }
        
        // Create transaction record
        txRecord := Transaction{
            FromAccountID: fromAccountID,
            ToAccountID:   toAccountID,
            Amount:        amount,
        }
        if err := tx.Create(&txRecord).Error; err != nil {
            return err
        }
        
        return nil // Return nil to commit the transaction
    })
    
    if err != nil {
        log.Fatalf("Transaction failed: %v", err)
    }
    
    fmt.Printf("Successfully transferred $%.2f from account %d to account %d\n", 
        amount, fromAccountID, toAccountID)
}
```

### Hooks and Callbacks

GORM provides hooks to execute code before or after specific operations:

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    LastLogin time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

// BeforeCreate hook - hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    // Hash password (using a proper hashing function in real apps)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    
    return nil
}

// AfterCreate hook - perform action after user creation
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    // For example, create default settings for user
    settings := UserSettings{
        UserID:           u.ID,
        NotificationsOn:  true,
        DarkModeEnabled:  false,
    }
    
    return tx.Create(&settings).Error
}

// BeforeUpdate hook - validate data before update
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    // Validate email format
    if u.Email != "" {
        if !isValidEmail(u.Email) {
            return errors.New("invalid email format")
        }
    }
    
    return nil
}

// AfterFind hook - do something after record is retrieved
func (u *User) AfterFind(tx *gorm.DB) (err error) {
    // Clear sensitive info when retrieved
    u.Password = "[REDACTED]"
    
    return nil
}
```

### Database Migrations with golang-migrate

The golang-migrate tool provides a robust way to manage database schema changes:

```go
// Structure for the migrations directory:
// migrations/
// ├── 000001_create_users_table.up.sql
// ├── 000001_create_users_table.down.sql
// ├── 000002_add_user_roles.up.sql
// ├── 000002_add_user_roles.down.sql
// └── ...
```

Example migration files:

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    age TINYINT UNSIGNED,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;

-- 000002_add_user_roles.up.sql
CREATE TABLE roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_roles (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- 000002_add_user_roles.down.sql
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
```

Create a migration utility in your Go application:

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDatabase handles database migrations using golang-migrate
func MigrateDatabase() error {
    // Get database connection details
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    
    // Connect to MySQL
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", 
        dbUser, dbPassword, dbHost, dbPort, dbName)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("error opening database connection: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }
    
    // Create migration instance
driver, err := mysql.WithInstance(db, &mysql.Config{})
if err != nil {
    return fmt.Errorf("error creating migration driver: %w", err)
}
m, err := migrate.NewWithDatabaseInstance(
    "file://migrations", // Directory containing migration files
    "mysql",
driver,
)
if err != nil {
    return fmt.Errorf("error creating migration instance: %w", err)
}
// Perform migration
if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    return fmt.Errorf("error applying migrations: %w", err)
}
log.Println("Database migrations applied successfully.")
```
### Common ORM Challenges

1. **Performance Bottlenecks**
   - Inefficient queries leading to slow response times.
   - Over-reliance on ORM features causing unnecessary database load.
   - Handling large datasets with pagination and indexing.

2. **Data Integrity Issues**
   - Ensuring transactions maintain consistency.
   - Managing concurrent updates to avoid race conditions.
   - Implementing soft deletes versus hard deletes properly.

3. **Schema Migration and Versioning**
   - Handling database schema changes without breaking production.
   - Rolling back migrations safely.
   - Managing different database engines efficiently.

4. **Error Handling and Debugging**
   - Identifying and managing ORM-related errors.
   - Logging SQL queries for debugging and optimization.
   - Handling edge cases with missing relationships or invalid constraints.

### Best Practices for ORM

1. **Efficient Query Optimization**
   - Use eager loading (`Preload`) to avoid N+1 query problems.
   - Write raw SQL when performance is critical.
   - Use database indexing wisely for fast lookups.

2. **Data Integrity and Consistency**
   - Use transactions for batch operations.
   - Ensure foreign key constraints are enforced.
   - Validate data at both the application and database levels.

3. **Scalability Considerations**
   - Optimize connection pooling for high-concurrency systems.
   - Shard databases when necessary for distributed workloads.
   - Cache frequently accessed data using Redis or Memcached.

4. **Testing and Deployment**
   - Use mock databases for unit testing.
   - Automate migrations with version control.
   - Ensure rollback strategies are in place for failed deployments.

### Learning Challenges in ORM

1. **Understanding ORM vs. Raw SQL Trade-offs**
   - Knowing when to use ORM versus writing direct SQL queries.
   - Avoiding ORM abstractions that degrade performance.

2. **Handling Database Transactions Properly**
   - Using transactions efficiently without unnecessary locks.
   - Handling rollback scenarios correctly.

3. **Debugging and Profiling ORM Queries**
   - Logging slow queries for performance tuning.
   - Using database explain plans to understand query execution.

### Recommended Resources for ORM

1. **Official Documentation & Tutorials**
   - [GORM Documentation](https://gorm.io/docs/)
   - [SQL Performance Best Practices](https://use-the-index-luke.com/)

2. **Books & Courses**
   - "SQL Performance Explained" by Markus Winand.
   - "GORM for Go Developers" Udemy course.

3. **Open Source Examples**
   - [RealWorld Example App with GORM](https://github.com/gothinkster/golang-gin-realworld-example-app)

### Reflection Questions

1. How does GORM compare to other ORM frameworks in different languages?
2. What are the advantages and drawbacks of using an ORM over raw SQL?
3. How can you optimize database queries when working with an ORM?
4. What are the best practices for managing schema migrations in a production environment?
5. How do you handle database consistency when working with distributed transactions?
