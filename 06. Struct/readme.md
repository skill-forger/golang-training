# Module 06: Structs in Go

## Overview
Structs are one of Go's most powerful features, allowing you to create custom data types by grouping related data together. They form the foundation of object-oriented programming in Go and enable you to model complex real-world entities in your code. While Go isn't traditionally considered an object-oriented language, structs with methods provide many of the same capabilities in a simpler, more direct approach.

## Learning Objectives
By the end of this module, you will:
- Understand how to define and use structs to model complex data
- Create methods that operate on struct data
- Implement struct composition to build complex types
- Use struct tags for metadata and serialization
- Apply best practices for designing and organizing struct-based code
- Implement practical applications using structs

## What are Structs?

Structs are user-defined types that group together variables of different data types under a single name. Think of them as blueprints for creating data objects that represent real-world entities or concepts.

### Basic Struct Definition and Instantiation

```go
// basic_struct.go
package main

import "fmt"

// Person defines a basic struct with three fields
type Person struct {
    Name    string  // Person's full name
    Age     int     // Person's age in years
    Email   string  // Contact email address
}

func main() {
    // Method 1: Creating a struct with named fields (recommended)
    alice := Person{
        Name:  "Alice Johnson",
        Age:   30,
        Email: "alice@example.com",
    }
    fmt.Println("Alice:", alice)
    
    // Method 2: Creating a struct with positional fields (order matters)
    // This method is less readable and not recommended for structs with many fields
    bob := Person{"Bob Smith", 25, "bob@example.com"}
    fmt.Println("Bob:", bob)
    
    // Method 3: Creating an empty struct (zero-value initialization)
    var charlie Person  // All fields set to zero values (empty string, 0, empty string)
    fmt.Println("Charlie (before):", charlie)
    
    // Assigning values to fields after creation
    charlie.Name = "Charlie Brown"
    charlie.Age = 22
    charlie.Email = "charlie@example.com"
    fmt.Println("Charlie (after):", charlie)
    
    // Accessing individual struct fields
    fmt.Printf("%s is %d years old and can be reached at %s\n", 
               alice.Name, alice.Age, alice.Email)
}
```

### Struct Field Visibility

Go controls access to struct fields through capitalization:

```go
// visibility.go
package main

import (
    "fmt"
    "example/user"  // Imaginary package
)

// User defines a struct with both exported and unexported fields
type User struct {
    Username string  // Exported (visible outside the package)
    password string  // Unexported (only visible within the package)
    Email    string  // Exported
    age      int     // Unexported
}

func main() {
    u := User{
        Username: "gopher",
        password: "secret123",  // Only visible within this package
        Email:    "gopher@example.com",
        age:      4,            // Only visible within this package
    }
    
    // Within the same package, we can access all fields
    fmt.Println("Username:", u.Username)
    fmt.Println("Password:", u.password)  // Works inside the same package
    
    // Importing code would only be able to access Username and Email
    // externalUser := user.New()
    // fmt.Println(externalUser.Username)  // Works
    // fmt.Println(externalUser.password)  // Compiler error
}
```

## Struct Methods: Adding Behavior

While structs define data, methods define behavior. Go lets you attach methods to struct types, creating an elegant way to encapsulate both data and operations.

### Value Receivers vs Pointer Receivers

```go
// struct_methods.go
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

// Area is a method with a value receiver
// It doesn't modify the original Rectangle
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Scale is a method with a pointer receiver
// It modifies the original Rectangle
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10.0, Height: 5.0}
    
    // Using a method with a value receiver
    area := rect.Area()
    fmt.Printf("Rectangle: %.2f x %.2f\n", rect.Width, rect.Height)
    fmt.Printf("Area: %.2f\n", area)
    
    // Using a method with a pointer receiver
    rect.Scale(2.0)
    fmt.Printf("After scaling: %.2f x %.2f\n", rect.Width, rect.Height)
    fmt.Printf("New area: %.2f\n", rect.Area())
}
```

### When to Use Value vs Pointer Receivers

- **Use value receivers when:**
  - The method doesn't modify the receiver
  - The struct is small and cheap to copy
  - You want immutability (the original struct is not affected)

- **Use pointer receivers when:**
  - The method needs to modify the receiver
  - The struct is large and would be expensive to copy
  - Consistency is important (all methods use the same receiver type)

## Struct Composition: Building Complex Types

Go favors composition over inheritance. Instead of creating complex inheritance hierarchies, you can embed one struct inside another to reuse fields and methods.

### Basic Composition with Embedding

```go
// composition.go
package main

import "fmt"

// Address stores location information
type Address struct {
    Street  string
    City    string
    State   string
    ZipCode string
}

// Person contains basic information about a person
type Person struct {
    Name    string
    Age     int
    Email   string
}

// Employee combines Person and Address through embedding
type Employee struct {
    Person          // Embedded struct (all fields accessible directly)
    Address         // Another embedded struct
    Position string
    Salary   float64
}

func main() {
    // Creating an employee with embedded types
    emp := Employee{
        Person: Person{
            Name:  "Alice Johnson",
            Age:   30,
            Email: "alice@company.com",
        },
        Address: Address{
            Street:  "123 Main St",
            City:    "Anytown",
            State:   "CA",
            ZipCode: "12345",
        },
        Position: "Software Engineer",
        Salary:   90000.0,
    }
    
    // Accessing fields directly from embedded structs
    fmt.Printf("Employee: %s, %s\n", emp.Name, emp.Position)
    fmt.Printf("Contact: %s\n", emp.Email)
    fmt.Printf("Location: %s, %s, %s %s\n", 
               emp.Street, emp.City, emp.State, emp.ZipCode)
    
    // You can also access through the embedded field name
    fmt.Printf("Age: %d (through Person.Age)\n", emp.Person.Age)
}
```

### Composition vs Inheritance

Composition offers several advantages over traditional inheritance:

1. **Flexibility**: Mix and match behavior without complex hierarchies
2. **Clarity**: The source of each field and method is explicit
3. **No Diamond Problem**: Avoids ambiguity when inheriting from multiple sources
4. **Simplicity**: Easier to understand and maintain

### Method Overriding with Composition

```go
// method_override.go
package main

import "fmt"

type Animal struct {
    Species string
}

func (a Animal) MakeSound() string {
    return "Some generic animal sound"
}

type Dog struct {
    Animal
    Breed string
}

// Override the MakeSound method for Dog
func (d Dog) MakeSound() string {
    return "Woof!"
}

func main() {
    // Create a generic animal
    animal := Animal{Species: "Unknown"}
    fmt.Printf("Animal says: %s\n", animal.MakeSound())
    
    // Create a dog
    dog := Dog{
        Animal: Animal{Species: "Canine"},
        Breed:  "Golden Retriever",
    }
    fmt.Printf("Dog says: %s\n", dog.MakeSound())
    
    // We can still access the original method
    fmt.Printf("Original sound: %s\n", dog.Animal.MakeSound())
}
```

## Struct Tags: Metadata and Serialization

Struct tags provide metadata about struct fields, commonly used for tasks like serialization, validation, and database mapping.

### JSON Serialization with Struct Tags

```go
// json_tags.go
package main

import (
    "encoding/json"
    "fmt"
)

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description,omitempty"` // Skip if empty
    Price       float64 `json:"price"`
    SKU         string  `json:"-"`                     // Ignore this field
    InStock     bool    `json:"in_stock"`
}

func main() {
    // Create a product
    product := Product{
        ID:          1001,
        Name:        "Mechanical Keyboard",
        Description: "", // This will be omitted in JSON
        Price:       79.99,
        SKU:         "KB-MEC-001", // This will be ignored in JSON
        InStock:     true,
    }
    
    // Convert to JSON
    jsonData, err := json.MarshalIndent(product, "", "  ")
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return
    }
    
    // Print the JSON
    fmt.Println("JSON Output:")
    fmt.Println(string(jsonData))
    
    // JSON output will be:
    // {
    //   "id": 1001,
    //   "name": "Mechanical Keyboard",
    //   "price": 79.99,
    //   "in_stock": true
    // }
    
    // Unmarshaling JSON back to a struct
    jsonString := `{
        "id": 1002,
        "name": "Wireless Mouse",
        "description": "Ergonomic design with long battery life",
        "price": 45.50,
        "in_stock": false
    }`
    
    var newProduct Product
    err = json.Unmarshal([]byte(jsonString), &newProduct)
    if err != nil {
        fmt.Println("Error unmarshaling JSON:", err)
        return
    }
    
    fmt.Println("\nUnmarshaled product:")
    fmt.Printf("ID: %d\n", newProduct.ID)
    fmt.Printf("Name: %s\n", newProduct.Name)
    fmt.Printf("Description: %s\n", newProduct.Description)
    fmt.Printf("Price: $%.2f\n", newProduct.Price)
    fmt.Printf("SKU: %s (not from JSON)\n", newProduct.SKU)
    fmt.Printf("In Stock: %t\n", newProduct.InStock)
}
```

### Common Tag Formats

- **JSON**: `json:"fieldname,options"`
- **XML**: `xml:"fieldname,options"`
- **YAML**: `yaml:"fieldname,options"`
- **Form**: `form:"fieldname"`
- **Validate**: `validate:"required,min=1,max=100"`
- **GORM (ORM)**: `gorm:"column:fieldname;type:varchar(100);unique_index"`

## Advanced Struct Techniques

### Anonymous Structs for Temporary Use

```go
// anonymous_struct.go
package main

import "fmt"

func main() {
    // Anonymous struct defined and initialized inline
    point := struct {
        X, Y int
    }{
        X: 10,
        Y: 20,
    }
    
    fmt.Printf("Point: (%d, %d)\n", point.X, point.Y)
    
    // Anonymous structs are useful for one-off data structures
    config := struct {
        Hostname string
        Port     int
        Debug    bool
    }{
        Hostname: "localhost",
        Port:     8080,
        Debug:    true,
    }
    
    fmt.Printf("Server config: %s:%d (debug: %t)\n", 
               config.Hostname, config.Port, config.Debug)
}
```

### Struct Equality and Comparison

```go
// struct_comparison.go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    p1 := Person{Name: "Alice", Age: 30}
    p2 := Person{Name: "Alice", Age: 30}
    p3 := Person{Name: "Bob", Age: 25}
    
    // Direct comparison works if all fields are comparable
    fmt.Println("p1 == p2:", p1 == p2)  // true
    fmt.Println("p1 == p3:", p1 == p3)  // false
    
    // For structs with fields that aren't comparable (maps, slices),
    // use reflect.DeepEqual
    type ComplexPerson struct {
        Name     string
        Age      int
        Hobbies  []string           // Slices aren't directly comparable
        Metadata map[string]string  // Maps aren't directly comparable
    }
    
    cp1 := ComplexPerson{
        Name:    "Charlie",
        Age:     35,
        Hobbies: []string{"Reading", "Hiking"},
        Metadata: map[string]string{
            "department": "Engineering",
            "level":      "Senior",
        },
    }
    
    cp2 := ComplexPerson{
        Name:    "Charlie",
        Age:     35,
        Hobbies: []string{"Reading", "Hiking"},
        Metadata: map[string]string{
            "department": "Engineering",
            "level":      "Senior",
        },
    }
    
    // Can't use == with complex structs that contain non-comparable types
    // fmt.Println(cp1 == cp2)  // This would cause a compile error
    
    // Instead use reflect.DeepEqual
    fmt.Println("cp1 DeepEqual cp2:", reflect.DeepEqual(cp1, cp2))  // true
}
```

### Constructor Functions for Structs

```go
// constructor.go
package main

import "fmt"

type Database struct {
    connectionString string
    maxConnections   int
    timeout          int
    isConnected      bool
}

// NewDatabase is a constructor function that ensures proper initialization
func NewDatabase(connString string) *Database {
    // Provide sensible defaults and validation
    if connString == "" {
        connString = "localhost:5432"
    }
    
    return &Database{
        connectionString: connString,
        maxConnections:   100,  // Default value
        timeout:          30,   // Default timeout in seconds
        isConnected:      false,
    }
}

// Connect method attempts to connect to the database
func (db *Database) Connect() error {
    // Simulating connection logic
    fmt.Printf("Connecting to %s...\n", db.connectionString)
    db.isConnected = true
    return nil
}

func main() {
    // Using the constructor instead of direct initialization
    db := NewDatabase("postgres://user:pass@remotehost:5432/mydb")
    
    fmt.Printf("Database configuration:\n")
    fmt.Printf("- Connection string: %s\n", db.connectionString)
    fmt.Printf("- Max connections: %d\n", db.maxConnections)
    fmt.Printf("- Timeout: %d seconds\n", db.timeout)
    
    db.Connect()
    fmt.Printf("Connection status: %t\n", db.isConnected)
}
```

## Best Practices for Using Structs

1. **Design for Clarity**
   - Keep structs focused on a single responsibility
   - Use meaningful field and method names
   - Document complex or non-obvious fields

2. **Choose Receivers Appropriately**
   - Use pointer receivers for methods that modify state
   - Use value receivers for immutable operations
   - Be consistent within a struct's method set

3. **Leverage Composition**
   - Prefer composition over complex type hierarchies
   - Use embedding to reuse code without inheritance
   - Keep embedded types orthogonal (separate concerns)

4. **Encapsulation**
   - Use unexported fields to hide implementation details
   - Provide methods or exported fields for controlled access
   - Create constructor functions for complex initialization

5. **Optimizations**
   - Order struct fields to minimize memory padding (largest to smallest)
   - Use pointers for large structs to avoid copying
   - Consider memory implications in struct design

## Practice Exercises

### Exercise 1: Employee Management System

Create a program that manages employees in a company. Implement the following:

- `Employee` struct with fields for name, position, salary, etc.
- `Department` struct that contains employees
- Methods for hiring, promoting, and transferring employees
- Implement sorting employees by different criteria (name, salary, etc.)

### Exercise 2: Shape Hierarchy

Create a geometry calculation system with:

- A `Shape` interface with methods for calculating area and perimeter
- Structs for different shapes (Circle, Rectangle, Triangle)
- Methods to implement the Shape interface for each struct
- Functions that can operate on any shape

### Exercise 3: E-Commerce Cart System

Implement a shopping cart system:

- `Product` struct with price, name, SKU
- `CartItem` struct that references a product and quantity
- `ShoppingCart` struct that manages a collection of cart items
- Methods for adding, removing, updating quantities
- Calculate subtotals, taxes, and final prices

## Summary

In this module, you've learned:
- How to define and use structs to model complex data
- Adding behavior to structs with methods
- Implementing composition with embedded structs
- Using struct tags for serialization and metadata
- Advanced struct techniques for comparison and construction
- Applying best practices for struct design
- Building real-world applications using structs

Structs are at the heart of Go's approach to programming with types. By mastering structs, you gain the ability to model almost any kind of data and build sophisticated, maintainable programs.

## Additional Resources

- [A Tour of Go: Structs](https://tour.golang.org/moretypes/2)
- [Effective Go: Structs](https://golang.org/doc/effective_go.html#structs)
- [Go by Example: Structs](https://gobyexample.com/structs)
- [Go by Example: Methods](https://gobyexample.com/methods)
- [Go by Example: Embedding](https://gobyexample.com/struct-embedding)
- [Practical Go: Real World Advice](https://dave.cheney.net/practical-go/presentations/qcon-china.html)
- [Package encoding/json](https://golang.org/pkg/encoding/json/) - For struct tags and serialization