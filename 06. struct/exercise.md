## Practical Exercises

### Exercise 1: Building a Library Management System

Design and implement a library management system using structs and methods. This exercise will help you understand how to model real-world entities as structs and implement operations as methods.

Your system should include:
1. A `Book` struct with fields for ID, title, author, publication year, and availability status
2. A `Member` struct with fields for ID, name, email, join date, books borrowed, and borrowing limit
3. A `BorrowRecord` struct to track book loans, including borrow date and due date
4. A `Library` struct that manages books, members, and borrowing records
5. Methods to:
   - Add books and members to the library
   - Allow members to borrow books with appropriate validation
   - Process book returns
   - Display library status
6. Error handling for various scenarios (book not found, unavailable books, etc.)
7. A demonstration in the `main` function showing the complete workflow

```go
// library.go
package main

import (
    "fmt"
    "time"
)

// Book represents a book in the library
type Book struct {
    ID            string
    Title         string
    Author        string
    PublishedYear int
    Available     bool
}

// Member represents a library member
type Member struct {
    ID        string
    Name      string
    Email     string
    JoinedOn  time.Time
    BooksOut  int
    MaxBooks  int
}

// BorrowRecord tracks a book being borrowed
type BorrowRecord struct {
    BookID      string
    MemberID    string
    BorrowedOn  time.Time
    DueDate     time.Time
    ReturnedOn  *time.Time // Pointer because it might be nil (not returned yet)
}

// Library manages the book collection and members
type Library struct {
    Name     string
    Books    map[string]*Book
    Members  map[string]*Member
    Borrows  []BorrowRecord
}

// NewLibrary creates a new library instance
func NewLibrary(name string) *Library {
    return &Library{
        Name:     name,
        Books:    make(map[string]*Book),
        Members:  make(map[string]*Member),
        Borrows:  []BorrowRecord{},
    }
}

// AddBook adds a book to the library
func (l *Library) AddBook(book Book) {
    l.Books[book.ID] = &book
}

// AddMember adds a member to the library
func (l *Library) AddMember(member Member) {
    l.Members[member.ID] = &member
}

// BorrowBook allows a member to borrow a book
func (l *Library) BorrowBook(bookID, memberID string) error {
    // Find the book
    book, found := l.Books[bookID]
    if !found {
        return fmt.Errorf("book not found")
    }
    
    // Check if the book is available
    if !book.Available {
        return fmt.Errorf("book is not available")
    }
    
    // Find the member
    member, found := l.Members[memberID]
    if !found {
        return fmt.Errorf("member not found")
    }
    
    // Check if the member can borrow more books
    if member.BooksOut >= member.MaxBooks {
        return fmt.Errorf("member has reached maximum number of books")
    }
    
    // Create a borrow record
    now := time.Now()
    borrowRecord := BorrowRecord{
        BookID:     bookID,
        MemberID:   memberID,
        BorrowedOn: now,
        DueDate:    now.AddDate(0, 0, 14), // Due in 14 days
    }
    
    // Update book and member
    book.Available = false
    member.BooksOut++
    
    // Add the record
    l.Borrows = append(l.Borrows, borrowRecord)
    
    return nil
}

// ReturnBook processes a book return
func (l *Library) ReturnBook(bookID, memberID string) error {
    // Find the book
    book, found := l.Books[bookID]
    if !found {
        return fmt.Errorf("book not found")
    }
    
    // Find the member
    member, found := l.Members[memberID]
    if !found {
        return fmt.Errorf("member not found")
    }
    
    // Find the borrow record
    recordIndex := -1
    for i, record := range l.Borrows {
        if record.BookID == bookID && record.MemberID == memberID && record.ReturnedOn == nil {
            recordIndex = i
            break
        }
    }
    
    if recordIndex == -1 {
        return fmt.Errorf("no active borrow record found")
    }
    
    // Update the record
    now := time.Now()
    l.Borrows[recordIndex].ReturnedOn = &now
    
    // Update book and member
    book.Available = true
    member.BooksOut--
    
    return nil
}

func main() {
    // Create a new library
    library := NewLibrary("Community Library")
    
    // Add books
    library.AddBook(Book{
        ID:            "B001",
        Title:         "The Go Programming Language",
        Author:        "Alan A. A. Donovan & Brian W. Kernighan",
        PublishedYear: 2015,
        Available:     true,
    })
    
    library.AddBook(Book{
        ID:            "B002",
        Title:         "Go in Action",
        Author:        "William Kennedy",
        PublishedYear: 2016,
        Available:     true,
    })
    
    // Add members
    library.AddMember(Member{
        ID:        "M001",
        Name:      "John Doe",
        Email:     "john@example.com",
        JoinedOn:  time.Now(),
        BooksOut:  0,
        MaxBooks:  3,
    })
    
    library.AddMember(Member{
        ID:        "M002",
        Name:      "Jane Smith",
        Email:     "jane@example.com",
        JoinedOn:  time.Now(),
        BooksOut:  0,
        MaxBooks:  5,
    })
    
    // Borrow a book
    err := library.BorrowBook("B001", "M001")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        fmt.Println("Book B001 borrowed by member M001")
    }
    
    // Display library status
    fmt.Println("\nLibrary Status:")
    fmt.Printf("Name: %s\n", library.Name)
    fmt.Printf("Books: %d\n", len(library.Books))
    fmt.Printf("Members: %d\n", len(library.Members))
    fmt.Printf("Active Borrows: %d\n", len(library.Borrows))
    
    // Return the book
    err = library.ReturnBook("B001", "M001")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        fmt.Println("\nBook B001 returned by member M001")
    }
}
```

### Exercise 2: Employee Management System

Create a comprehensive employee management system that models an organization's structure. This exercise will demonstrate how structs can represent complex relationships and operations.

Your implementation should include:
1. An `Address` struct for storing location information
2. An `Employee` struct with personal details, employment information, and a nested Address
3. A `Company` struct for managing employees and department organization
4. Methods for:
   - Adding employees to the company
   - Updating employee salaries
   - Transferring employees between departments
   - Marking employees as inactive (terminated)
   - Generating department statistics
5. Helper methods for employees (e.g., `FullName()`, `YearsOfService()`)
6. A demonstration showing typical HR operations

```go
// employee_system.go
package main

import (
    "fmt"
    "time"
)

// Address represents a physical address
type Address struct {
    Street     string
    City       string
    State      string
    PostalCode string
    Country    string
}

// Employee defines the base employee structure
type Employee struct {
    ID            string
    FirstName     string
    LastName      string
    Email         string
    HireDate      time.Time
    Address       Address
    Position      string
    Salary        float64
    ManagerID     string
    Department    string
    IsActive      bool
}

// FullName returns the employee's full name
func (e Employee) FullName() string {
    return e.FirstName + " " + e.LastName
}

// YearsOfService calculates the years an employee has worked
func (e Employee) YearsOfService() float64 {
    now := time.Now()
    duration := now.Sub(e.HireDate)
    return duration.Hours() / 24 / 365.25
}

// Company contains all employees and departments
type Company struct {
    Name         string
    Employees    map[string]*Employee
    Departments  map[string][]string // Department name -> slice of employee IDs
}

// NewCompany creates a new company
func NewCompany(name string) *Company {
    return &Company{
        Name:        name,
        Employees:   make(map[string]*Employee),
        Departments: make(map[string][]string),
    }
}

// AddEmployee adds a new employee to the company
func (c *Company) AddEmployee(e Employee) error {
    // Check if employee already exists
    if _, exists := c.Employees[e.ID]; exists {
        return fmt.Errorf("employee with ID %s already exists", e.ID)
    }
    
    // Add employee to the map
    c.Employees[e.ID] = &e
    
    // Add employee to the department
    if e.Department != "" {
        c.Departments[e.Department] = append(c.Departments[e.Department], e.ID)
    }
    
    return nil
}

// UpdateSalary updates an employee's salary
func (c *Company) UpdateSalary(employeeID string, newSalary float64) error {
    employee, exists := c.Employees[employeeID]
    if !exists {
        return fmt.Errorf("employee with ID %s not found", employeeID)
    }
    
    employee.Salary = newSalary
    return nil
}

// TransferEmployee moves an employee to a new department
func (c *Company) TransferEmployee(employeeID, newDepartment string) error {
    employee, exists := c.Employees[employeeID]
    if !exists {
        return fmt.Errorf("employee with ID %s not found", employeeID)
    }
    
    // Remove from old department
    oldDepartment := employee.Department
    if oldDepartment != "" {
        for i, id := range c.Departments[oldDepartment] {
            if id == employeeID {
                // Remove by replacing with last element and truncating
                c.Departments[oldDepartment][i] = c.Departments[oldDepartment][len(c.Departments[oldDepartment])-1]
                c.Departments[oldDepartment] = c.Departments[oldDepartment][:len(c.Departments[oldDepartment])-1]
                break
            }
        }
    }
    
    // Add to new department
    employee.Department = newDepartment
    c.Departments[newDepartment] = append(c.Departments[newDepartment], employeeID)
    
    return nil
}

// TerminateEmployee marks an employee as inactive
func (c *Company) TerminateEmployee(employeeID string) error {
    employee, exists := c.Employees[employeeID]
    if !exists {
        return fmt.Errorf("employee with ID %s not found", employeeID)
    }
    
    employee.IsActive = false
    return nil
}

// GetDepartmentStats returns statistics about a department
func (c *Company) GetDepartmentStats(department string) (count int, avgSalary float64, avgService float64) {
    employeeIDs, exists := c.Departments[department]
    if !exists || len(employeeIDs) == 0 {
        return 0, 0, 0
    }
    
    activeCount := 0
    totalSalary := 0.0
    totalService := 0.0
    
    for _, id := range employeeIDs {
        employee := c.Employees[id]
        if employee.IsActive {
            activeCount++
            totalSalary += employee.Salary
            totalService += employee.YearsOfService()
        }
    }
    
    if activeCount == 0 {
        return 0, 0, 0
    }
    
    return activeCount, totalSalary / float64(activeCount), totalService / float64(activeCount)
}

func main() {
    // Create a new company
    company := NewCompany("Tech Innovations Inc.")
    
    // Add employees
    employees := []Employee{
        {
            ID:        "E001",
            FirstName: "John",
            LastName:  "Doe",
            Email:     "john.doe@example.com",
            HireDate:  time.Date(2018, 5, 15, 0, 0, 0, 0, time.UTC),
            Address: Address{
                Street:     "123 Main St",
                City:       "Boston",
                State:      "MA",
                PostalCode: "02108",
                Country:    "USA",
            },
            Position:   "Software Engineer",
            Salary:     95000,
            Department: "Engineering",
            IsActive:   true,
        },
        {
            ID:        "E002",
            FirstName: "Jane",
            LastName:  "Smith",
            Email:     "jane.smith@example.com",
            HireDate:  time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC),
            Address: Address{
                Street:     "456 Park Ave",
                City:       "New York",
                State:      "NY",
                PostalCode: "10022",
                Country:    "USA",
            },
            Position:   "Marketing Specialist",
            Salary:     85000,
            Department: "Marketing",
            IsActive:   true,
        },
        {
            ID:        "E003",
            FirstName: "Robert",
            LastName:  "Johnson",
            Email:     "robert.johnson@example.com",
            HireDate:  time.Date(2017, 1, 20, 0, 0, 0, 0, time.UTC),
            Address: Address{
                Street:     "789 Oak St",
                City:       "San Francisco",
                State:      "CA",
                PostalCode: "94107",
                Country:    "USA",
            },
            Position:   "Senior Software Engineer",
            Salary:     120000,
            Department: "Engineering",
            IsActive:   true,
        },
    }
    
    for _, employee := range employees {
        err := company.AddEmployee(employee)
        if err != nil {
            fmt.Printf("Error adding employee: %s\n", err)
        }
    }
    
    // Display company information
    fmt.Printf("Company: %s\n", company.Name)
    fmt.Printf("Total Employees: %d\n\n", len(company.Employees))
    
    // Display departments
    for dept, empIDs := range company.Departments {
        count, avgSalary, avgService := company.GetDepartmentStats(dept)
        fmt.Printf("Department: %s\n", dept)
        fmt.Printf("  Active Employees: %d\n", count)
        fmt.Printf("  Average Salary: $%.2f\n", avgSalary)
        fmt.Printf("  Average Years of Service: %.1f\n\n", avgService)
    }
    
    // Perform some operations
    fmt.Println("Operations:")
    
    // Give a raise
    err := company.UpdateSalary("E001", 100000)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        employee := company.Employees["E001"]
        fmt.Printf("- %s received a raise to $%.2f\n", employee.FullName(), employee.Salary)
    }
    
    // Transfer an employee
    err = company.TransferEmployee("E002", "Sales")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        employee := company.Employees["E002"]
        fmt.Printf("- %s transferred to %s department\n", employee.FullName(), employee.Department)
    }
    
    // Terminate an employee
    err = company.TerminateEmployee("E003")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        employee := company.Employees["E003"]
        fmt.Printf("- %s is no longer active\n", employee.FullName())
    }
}
```

### Exercise 3: Product Inventory System

Develop an inventory management system for tracking products, stock levels, and transactions. This exercise will show how to use structs to model a business system with complex operations.

Your system should include:
1. A `Product` struct with detailed product information (SKU, name, description, pricing, stock levels)
2. A `Transaction` struct that records inventory changes (purchases, sales, adjustments)
3. An `Inventory` struct that manages products and their transaction history
4. Methods to:
   - Add new products to inventory
   - Record product purchases (stock increases)
   - Record product sales (stock decreases) with validation
   - Adjust stock levels (e.g., after inventory count)
   - Generate reports (low stock products, inventory value)
5. Helper methods for products (e.g., calculating profit margins, checking reorder needs)
6. A demonstration that includes various inventory operations and reporting

```go
// inventory_system.go
package main

import (
    "fmt"
    "time"
)

// Product represents an item in the inventory
type Product struct {
    SKU         string
    Name        string
    Description string
    Category    string
    Price       float64
    Cost        float64
    StockLevel  int
    ReorderLevel int
    Supplier    string
    DateAdded   time.Time
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
    ID          string
    ProductSKU  string
    Type        string // "purchase", "sale", "adjustment"
    Quantity    int
    Date        time.Time
    Reference   string // invoice or order number
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
            SKU:         "LAPTOP001",
            Name:        "Pro Laptop 15\"",
            Description: "High-performance laptop with 16GB RAM",
            Category:    "Electronics",
            Price:       1299.99,
            Cost:        950.00,
            StockLevel:  0,
            ReorderLevel: 5,
            Supplier:    "TechSuppliers Inc.",
        },
        {
            SKU:         "PHONE001",
            Name:        "Smartphone X",
            Description: "Latest model smartphone",
            Category:    "Electronics",
            Price:       799.99,
            Cost:        550.00,
            StockLevel:  0,
            ReorderLevel: 10,
            Supplier:    "MobileTech Inc.",
        },
        {
            SKU:         "CHAIR001",
            Name:        "Ergonomic Office Chair",
            Description: "Adjustable office chair with lumbar support",
            Category:    "Furniture",
            Price:       249.99,
            Cost:        125.00,
            StockLevel:  0,
            ReorderLevel: 3,
            Supplier:    "Office Furnishings Co.",
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
```
