# Module 08: Interfaces in Go

## Table of Contents

<ol>
    <li><a href="#objectives">Objectives</a></li>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#the-essence-of-interfaces">The Essence of Interfaces</a></li>
    <li><a href="#interface-declaration-and-implementation">Interface Declaration and Implementation</a></li>
    <li><a href="#empty-interface">Empty Interface</a></li>
    <li><a href="#interface-composition">Interface Composition</a></li>
    <li><a href="#type-assertions-and-type-switches">Type Assertions and Type Switches</a></li>
    <li><a href="#common-patterns">Common Patterns</a></li>
    <li><a href="#common-mistakes">Common Mistakes</a></li>
    <li><a href="#best-practices">Best Practices</a></li>
    <li><a href="#practice-exercises">Practice Exercises</a></li>
</ol>

## Objectives

By the end of this module, you will:
- Understand the concept and purpose of interfaces in Go
- Master interface declaration and implementation
- Learn how to use the empty interface for type-agnostic functions
- Apply interface composition to build complex behaviors
- Implement type assertions and type switches
- Recognize common interface design patterns and best practices
- Avoid common pitfalls when working with interfaces

## Overview

Interfaces are a cornerstone of Go's type system, 
providing a powerful way to express abstraction and polymorphism without the complexity of inheritance-based systems. 
Unlike many object-oriented languages, Go's interfaces are implemented implicitly, 
focusing on what types can do rather than what they are. This approach leads to more flexible, 
decoupled code that's easier to test and maintain.

## The Essence of Interfaces

Interfaces in Go define behavior by specifying a set of method signatures. 
Unlike many languages, Go interfaces are implemented implicitly - 
there's no explicit declaration of intent to implement an interface.

### Interface as a Contract

Think of an interface as a contract that defines what a type can do, not what it is:

```go
// interface_basics.go
package main

import "fmt"

// Sounder defines a behavior - making a sound
type Sounder interface {
	MakeSound() string
}

// Types that implement Sounder:

// Dog implicitly implements Sounder
type Dog struct {
	Name  string
	Breed string
}

func (d Dog) MakeSound() string {
	return "Woof!"
}

// Cat implicitly implements Sounder
type Cat struct {
	Name  string
	Color string
}

func (c Cat) MakeSound() string {
	return "Meow!"
}

// A function that can work with ANY type that fulfills the Sounder contract
func AnimalConcert(animals []Sounder) {
	fmt.Println("Animal concert begins:")
	for _, animal := range animals {
		fmt.Println(animal.MakeSound())
	}
	fmt.Println("Concert ends")
}

func main() {
	// Create different concrete types
	dog := Dog{Name: "Buddy", Breed: "Golden Retriever"}
	cat := Cat{Name: "Whiskers", Color: "Orange"}

	// Both can be treated as Sounders
	animals := []Sounder{dog, cat}

	// Pass them to a function expecting the interface
	AnimalConcert(animals)
}

```

### Key Characteristics of Go Interfaces

1. **Implicit Implementation**: Types implement interfaces automatically by implementing all required methods
2. **Method-Based**: Only methods define interfaces, not fields
3. **Decoupling**: Interfaces separate what something does from how it does it
4. **Runtime Verification**: Interface compliance is checked at runtime, not compile time

## Interface Declaration and Implementation

### Defining an Interface

```go
// interface_definition.go
package main

import (
    "fmt"
    "math"
)

// Shape defines a geometric shape interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Circle implements the Shape interface
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Rectangle implements the Shape interface
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2*r.Width + 2*r.Height
}

// PrintShapeInfo works with any Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    c := Circle{Radius: 5}
    r := Rectangle{Width: 3, Height: 4}
    
    fmt.Println("Circle:")
    PrintShapeInfo(c)
    
    fmt.Println("\nRectangle:")
    PrintShapeInfo(r)
}
```

### Implementation Rules
- A type implements an interface by implementing all its methods
- No explicit declaration of intent is needed
- A type can implement multiple interfaces simultaneously
- Method signatures must match exactly (including parameter and return types)

## Empty Interface

The empty interface (`interface{}`) is a special case that's implemented by all types, making it Go's approach to generics (prior to Go 1.18).

```go
// empty_interface.go
package main

import (
    "fmt"
    "strings"
)

// PrintAny can accept any type of parameter
func PrintAny(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

// ProcessValue demonstrates type switching
func ProcessValue(v interface{}) {
    // Type switch allows different handling based on the concrete type
    switch val := v.(type) {
    case int:
        fmt.Println("Integer:", val*2)
    case string:
        fmt.Println("String:", strings.ToUpper(val))
    case []byte:
        fmt.Println("Byte slice of length:", len(val))
    case bool:
        if val {
            fmt.Println("Boolean: true")
        } else {
            fmt.Println("Boolean: false")
        }
    default:
        fmt.Println("Unknown type")
    }
}

func main() {
    // The empty interface can hold values of any type
    var x interface{}
    
    x = 42
    PrintAny(x)
    
    x = "Hello, interfaces"
    PrintAny(x)
    
    x = true
    PrintAny(x)
    
    x = []int{1, 2, 3}
    PrintAny(x)
    
    // Process different types with the same function
    fmt.Println("\nProcessing different types:")
    ProcessValue(100)
    ProcessValue("go programming")
    ProcessValue([]byte{65, 66, 67})
    ProcessValue(true)
    ProcessValue(struct{ name string }{"Custom struct"})
}
```

### Key Points About the Empty Interface
- It has no methods, so all types implement it
- Useful for functions that need to accept any type
- Requires type assertions or type switches to access the underlying value
- Similar to `Object` in Java or `any` in TypeScript
- Use with caution as it sacrifices type safety

## Interface Composition

Interfaces in Go can be composed by embedding other interfaces, creating more complex behavior contracts.

```go
// interface_composition.go
package main

import "fmt"

// Basic interfaces with single responsibilities
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Composed interfaces
type ReadWriter interface {
    Reader
    Writer
}

type ReadCloser interface {
    Reader
    Closer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// SimpleFile implements all three basic interfaces
type SimpleFile struct {
    data      []byte
    isOpen    bool
    readPos   int
    writePos  int
}

func NewSimpleFile() *SimpleFile {
    return &SimpleFile{
        data:   make([]byte, 0, 1024),
        isOpen: true,
    }
}

// Implement Read method
func (f *SimpleFile) Read(p []byte) (n int, err error) {
    if !f.isOpen {
        return 0, fmt.Errorf("file is closed")
    }
    
    if f.readPos >= len(f.data) {
        return 0, fmt.Errorf("EOF")
    }
    
    n = copy(p, f.data[f.readPos:])
    f.readPos += n
    return n, nil
}

// Implement Write method
func (f *SimpleFile) Write(p []byte) (n int, err error) {
    if !f.isOpen {
        return 0, fmt.Errorf("file is closed")
    }
    
    // Ensure capacity
    if f.writePos+len(p) > cap(f.data) {
        // Grow the slice
        newData := make([]byte, len(f.data), (cap(f.data)+len(p))*2)
        copy(newData, f.data)
        f.data = newData
    }
    
    // Extend if needed
    if f.writePos+len(p) > len(f.data) {
        f.data = f.data[:f.writePos+len(p)]
    }
    
    n = copy(f.data[f.writePos:], p)
    f.writePos += n
    return n, nil
}

// Implement Close method
func (f *SimpleFile) Close() error {
    if !f.isOpen {
        return fmt.Errorf("file already closed")
    }
    
    f.isOpen = false
    return nil
}

// Functions that work with different interface combinations

func CopyData(r Reader, w Writer) error {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf)
        if err != nil {
            if err.Error() == "EOF" {
                return nil
            }
            return err
        }
        
        _, err = w.Write(buf[:n])
        if err != nil {
            return err
        }
    }
}

func main() {
    file1 := NewSimpleFile()
    file2 := NewSimpleFile()
    
    // Write some data to file1
    message := []byte("Interface composition in Go is powerful!")
    _, err := file1.Write(message)
    if err != nil {
        fmt.Println("Error writing:", err)
        return
    }
    
    // Reset read position
    file1.readPos = 0
    
    // Copy data from file1 to file2 using interfaces
    err = CopyData(file1, file2)
    if err != nil {
        fmt.Println("Error copying:", err)
        return
    }
    
    // Reset read position and read from file2
    file2.readPos = 0
    readBuf := make([]byte, 1024)
    n, err := file2.Read(readBuf)
    if err != nil && err.Error() != "EOF" {
        fmt.Println("Error reading:", err)
        return
    }
    
    fmt.Println("Read from file2:", string(readBuf[:n]))
    
    // Use the file as a ReadWriteCloser
    var rwc ReadWriteCloser = file1
    rwc.Close()
    
    // Try to write after closing
    _, err = file1.Write([]byte("More data"))
    fmt.Println("Expected error after close:", err)
}
```

### Benefits of Interface Composition
- Builds complex behaviors from simpler ones
- Promotes the Single Responsibility Principle
- Creates flexible API contracts
- Minimizes interface pollution with large interfaces

## Type Assertions and Type Switches

Type assertions and switches allow you to safely extract the concrete type from an interface value.

```go
// type_assertions.go
package main

import "fmt"

func main() {
    // Start with an interface value
    var i interface{} = "hello"
    
    // Type assertion with explicit check (safe)
    s, ok := i.(string)
    if ok {
        fmt.Println("String value:", s)
    } else {
        fmt.Println("Not a string")
    }
    
    // Type assertion without check (unsafe)
    // Will panic if the type doesn't match
    // s = i.(string)
    
    // Try an incorrect assertion
    n, ok := i.(int)
    if ok {
        fmt.Println("Integer value:", n)
    } else {
        fmt.Println("Not an integer")
    }
    
    // Type switch for multiple type checks
    checkType(i)
    checkType(42)
    checkType(3.14)
    checkType([]string{"a", "b", "c"})
    checkType(map[string]int{"one": 1})
}

func checkType(v interface{}) {
    fmt.Printf("Checking value: %v\n", v)
    
    switch x := v.(type) {
    case nil:
        fmt.Println("Type: nil")
    
    case int:
        fmt.Println("Type: int, Value squared:", x*x)
    
    case float64:
        fmt.Println("Type: float64, Value doubled:", x*2)
    
    case string:
        fmt.Println("Type: string, Length:", len(x))
    
    case []string:
        fmt.Println("Type: []string, Elements:", len(x))
        for i, s := range x {
            fmt.Printf("  %d: %s\n", i, s)
        }
    
    case map[string]int:
        fmt.Println("Type: map[string]int, Keys:", len(x))
        for k, v := range x {
            fmt.Printf("  %s: %d\n", k, v)
        }
    
    default:
        fmt.Printf("Type: %T (unhandled specific type)\n", x)
    }
    
    fmt.Println()
}
```

### Type Assertion Best Practices
1. **Always use the two-return form** (`value, ok := x.(Type)`) to avoid panics
2. **Consider type switches** for multiple type possibilities
3. **Check for interface implementation** rather than concrete types when possible
4. **Be specific about the types** you expect and handle

## Common Patterns

### The io.Reader and io.Writer Pattern
One of the most powerful patterns in the Go standard library is the `io.Reader` and `io.Writer` interfaces:

```go
// reader_writer_pattern.go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

func main() {
    // Various types that implement io.Reader
    sources := []io.Reader{
        strings.NewReader("string reader source"),
        bytes.NewReader([]byte("bytes reader source")),
        bytes.NewBuffer([]byte("bytes buffer source")),
    }
    
    // Process all readers the same way
    for i, source := range sources {
        fmt.Printf("Source #%d:\n", i+1)
        ProcessReader(source)
        fmt.Println()
    }
    
    // Various types that implement io.Writer
    var stringBuffer strings.Builder
    bytesBuffer := bytes.NewBuffer([]byte{})
    
    // Write to multiple destinations
    destinations := []io.Writer{
        &stringBuffer,
        bytesBuffer,
    }
    
    for i, dest := range destinations {
        fmt.Printf("Destination #%d:\n", i+1)
        WriteData(dest, fmt.Sprintf("Hello, destination #%d!", i+1))
    }
    
    // Check the results
    fmt.Println("\nString builder result:", stringBuffer.String())
    fmt.Println("Bytes buffer result:", bytesBuffer.String())
}

func ProcessReader(r io.Reader) {
    // Read up to 1024 bytes
    buf := make([]byte, 1024)
    n, err := r.Read(buf)
    
    if err != nil && err != io.EOF {
        fmt.Println("Error reading:", err)
        return
    }
    
    fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
}

func WriteData(w io.Writer, data string) {
    n, err := w.Write([]byte(data))
    if err != nil {
        fmt.Println("Error writing:", err)
        return
    }
    
    fmt.Printf("Wrote %d bytes\n", n)
}
```

### The Stringer Interface
The `fmt` package uses the `Stringer` interface for custom string representations:

```go
// stringer_pattern.go
package main

import "fmt"

// Define a custom type
type Person struct {
    FirstName string
    LastName  string
    Age       int
}

// Implement the Stringer interface
func (p Person) String() string {
    return fmt.Sprintf("%s %s (%d years)", p.FirstName, p.LastName, p.Age)
}

func main() {
    // Create a Person
    person := Person{
        FirstName: "John",
        LastName:  "Doe",
        Age:       30,
    }
    
    // Thanks to the String() method, fmt.Println will use our custom format
    fmt.Println("Person:", person)
    
    // Other ways fmt uses the Stringer interface
    fmt.Printf("Person with %%v: %v\n", person)
    fmt.Printf("Person with %%s: %s\n", person)
    
    // Without Stringer, we'd see the default struct representation
    type SimpleStruct struct {
        Name string
        ID   int
    }
    
    simple := SimpleStruct{Name: "Test", ID: 123}
    fmt.Println("Simple struct (no Stringer):", simple)
}
```

### The Error Interface
The `error` interface is perhaps the most ubiquitous interface in Go:

```go
// error_interface.go
package main

import (
    "fmt"
    "time"
)

// Custom error type
type TimeoutError struct {
    Operation string
    Timeout   time.Duration
}

// Implement the error interface
func (e *TimeoutError) Error() string {
    return fmt.Sprintf("operation %s timed out after %v", e.Operation, e.Timeout)
}

// Function that returns the custom error
func fetchData(timeout time.Duration) ([]byte, error) {
    if timeout < 100*time.Millisecond {
        return nil, &TimeoutError{
            Operation: "fetchData",
            Timeout:   timeout,
        }
    }
    
    // Simulating successful data fetch
    return []byte("some data"), nil
}

func main() {
    // Test with timeout too short
    data, err := fetchData(50 * time.Millisecond)
    if err != nil {
        fmt.Println("Error occurred:", err)
        
        // Type assertion to access custom error fields
        if timeoutErr, ok := err.(*TimeoutError); ok {
            fmt.Printf("Timeout error in operation '%s' with timeout %v\n",
                timeoutErr.Operation, timeoutErr.Timeout)
        }
    } else {
        fmt.Println("Data fetched:", string(data))
    }
    
    // Test with sufficient timeout
    data, err = fetchData(200 * time.Millisecond)
    if err != nil {
        fmt.Println("Error occurred:", err)
    } else {
        fmt.Println("Data fetched:", string(data))
    }
}
```

## Common Mistakes

### Interface Nil vs Nil Comparison Trap

```go
// nil_interface.go
package main

import "fmt"

type MyError struct {
    Msg string
}

func (e *MyError) Error() string {
    if e == nil {
        return "nil error"
    }
    return e.Msg
}

func mayReturnError(fail bool) error {
    if fail {
        return &MyError{Msg: "something failed"}
    }
    
    // This returns a nil *MyError, not a nil error interface
    return (*MyError)(nil)
}

func main() {
    // This won't work as expected
    err := mayReturnError(false)
    
    // The interface is NOT nil because it contains a type description
    if err != nil {
        fmt.Println("Error is not nil, even though the value is nil!")
        fmt.Printf("Error type: %T, value: %v\n", err, err)
    } else {
        fmt.Println("Error is nil")
    }
    
    // To fix, return a true nil:
    var nilError error
    fmt.Println("True nil error:", nilError == nil) // This is true
}
```

### Interface Pollution

```go
// interface_pollution.go
package main

import "fmt"

// Example of interface pollution: too many unnecessary interfaces

// BAD: Creating interfaces for every type
type UserGetter interface {
    GetUser(id string) (User, error)
}

type UserSaver interface {
    SaveUser(user User) error
}

type UserDeleter interface {
    DeleteUser(id string) error
}

type UserService interface {
    UserGetter
    UserSaver
    UserDeleter
}

// GOOD: Only create interfaces when you need abstraction
type User struct {
    ID   string
    Name string
}

// Concrete implementation
type UserRepository struct {
    // Storage fields
}

func (r *UserRepository) GetUser(id string) (User, error) {
    // Implementation
    return User{ID: id, Name: "Example"}, nil
}

func (r *UserRepository) SaveUser(user User) error {
    // Implementation
    return nil
}

func (r *UserRepository) DeleteUser(id string) error {
    // Implementation
    return nil
}

// Only create interfaces at the point of consumption, if needed:

type UserCache struct {
    // Cache fields
    repo *UserRepository
}

// This function only needs the GetUser functionality
func (c *UserCache) GetUserWithCache(id string, getter interface{ GetUser(string) (User, error) }) (User, error) {
    // Check cache first...
    
    // If not in cache, use the provided getter
    user, err := getter.GetUser(id)
    if err != nil {
        return User{}, err
    }
    
    // Cache the result...
    
    return user, nil
}

func main() {
    repo := &UserRepository{}
    cache := &UserCache{repo: repo}
    
    // Use the repo through the minimal interface
    user, err := cache.GetUserWithCache("123", repo)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("User found:", user)
}
```

## Best Practices

### The Interface Size Principle

**Go Proverb**: "The bigger the interface, the weaker the abstraction."

Small, focused interfaces provide better abstraction:

```go
// interface_size.go
package main

import "fmt"

// BAD: Large interface with many methods
type FileManager interface {
    Open(name string) error
    Close() error
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Seek(offset int64, whence int) (int64, error)
    Stat() (FileInfo, error)
    Truncate(size int64) error
    // ... and many more methods
}

// GOOD: Small, focused interfaces
type Opener interface {
    Open(name string) error
}

type Closer interface {
    Close() error
}

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Function using small interfaces
func CopyData(r Reader, w Writer) error {
    data := make([]byte, 1024)
    n, err := r.Read(data)
    if err != nil {
        return err
    }
    
    _, err = w.Write(data[:n])
    return err
}
```

### Accept Interfaces, Return Concrete Types
```go
// interfaces_vs_concrete.go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

// Good: Accept interface
func CountLetters(r io.Reader) (int, error) {
    buf := make([]byte, 1024)
    count := 0
    
    for {
        n, err := r.Read(buf)
        for i := 0; i < n; i++ {
            if (buf[i] >= 'A' && buf[i] <= 'Z') || (buf[i] >= 'a' && buf[i] <= 'z') {
                count++
            }
        }
        
        if err == io.EOF {
            break
        }
        if err != nil {
            return 0, err
        }
    }
    
    return count, nil
}

// Good: Return concrete type
func CreateLetterCounter() *bytes.Buffer {
    return bytes.NewBuffer(nil)
}

func main() {
    // We can pass any reader to CountLetters
    stringReader := strings.NewReader("Hello, Go interfaces are great!")
    count, err := CountLetters(stringReader)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Letter count:", count)
    
    // When we get a concrete type, we know exactly what we're working with
    buffer := CreateLetterCounter()
    buffer.WriteString("More text to analyze")
    
    // We can still pass it to functions accepting interfaces
    count, _ = CountLetters(buffer)
    fmt.Println("Letter count from buffer:", count)
}
```

### Test with Interfaces, Not Implementations
Interfaces make testing easier by allowing mock implementations:

```go
// testing_with_interfaces.go
package main

import "fmt"

// Define the interface
type DataStore interface {
    Save(key string, value string) error
    Load(key string) (string, error)
}

// Real implementation
type DatabaseStore struct {
    // In a real app, this would have database connection details
}

func (db *DatabaseStore) Save(key string, value string) error {
    // In a real app, this would save to a database
    fmt.Printf("Saving %s=%s to database\n", key, value)
    return nil
}

func (db *DatabaseStore) Load(key string) (string, error) {
    // In a real app, this would load from a database
    fmt.Printf("Loading %s from database\n", key)
    return "database_value", nil
}

// Mock implementation for testing
type MockStore struct {
    Data map[string]string
}

func NewMockStore() *MockStore {
    return &MockStore{
        Data: make(map[string]string),
    }
}

func (m *MockStore) Save(key string, value string) error {
    m.Data[key] = value
    return nil
}

func (m *MockStore) Load(key string) (string, error) {
    value, exists := m.Data[key]
    if !exists {
        return "", fmt.Errorf("key %s not found", key)
    }
    return value, nil
}

// Business logic using the interface
type UserService struct {
    store DataStore
}

func NewUserService(store DataStore) *UserService {
    return &UserService{store: store}
}

func (s *UserService) SavePreference(userID string, preference string) error {
    key := fmt.Sprintf("user:%s:pref", userID)
    return s.store.Save(key, preference)
}

func (s *UserService) GetPreference(userID string) (string, error) {
    key := fmt.Sprintf("user:%s:pref", userID)
    return s.store.Load(key)
}

func main() {
    // In production, we'd use the real database
    // realDB := &DatabaseStore{}
    // service := NewUserService(realDB)
    
    // For this example, we'll use the mock
    mockStore := NewMockStore()
    service := NewUserService(mockStore)
    
    // Use the service
    err := service.SavePreference("user123", "dark_mode")
    if err != nil {
        fmt.Println("Error saving preference:", err)
        return
    }
    
    pref, err := service.GetPreference("user123")
    if err != nil {
        fmt.Println("Error getting preference:", err)
        return
    }
    
    fmt.Println("User preference:", pref)
    
    // In a real test, we could verify the mock directly
    fmt.Println("Mock store contents:", mockStore.Data)
}
```

## Practice Exercises

### Exercise 1: Event System with Interface-Based Pub/Sub
Build an event management system using interfaces to implement a publisher-subscriber pattern.
This exercise demonstrates how interfaces enable flexible and loosely coupled communication between components.

Your implementation should include:
1. An `Event` interface that defines methods to access event information:
    - `Type()` to get the event category
    - `Data()` to retrieve event data
    - `Timestamp()` to get when the event occurred
2. A concrete implementation of the `Event` interface (`BaseEvent`)
3. An `EventHandler` interface for components that process events
4. A function type that implements the `EventHandler` interface for convenient usage
5. An `EventBus` that manages:
    - Subscriptions to different event types
    - Event publishing to appropriate handlers
    - Concurrency safety using mutexes
6. A demonstration that shows:
    - Subscribing to specific event types
    - Publishing different kinds of events
    - Handling events with type assertions
    - Processing events asynchronously

This exercise illustrates how interfaces can create flexible,
extensible systems where components interact without tight coupling.

### Exercise 2: Shape Calculator with Interface Hierarchy
Create a shape calculation system that uses interfaces to handle different geometric shapes uniformly.
This exercise shows how interfaces enable polymorphism in Go.

Your implementation should include:
1. A base `Shape` interface that requires methods for:
    - Calculating area
    - Calculating perimeter
    - Getting shape name
2. Multiple concrete shape implementations (Circle, Rectangle, Triangle)
3. A more specialized `ThreeDimensionalShape` interface that extends the base interface with volume calculation
4. Implementations of 3D shapes (Sphere, Cube)
5. A shape processor that can:
    - Handle any shape type through the interface
    - Sort shapes by area
    - Filter shapes by type
    - Generate reports on shape properties
6. A demonstration showing how the same functions can process different shape types uniformly

### Exercise 3: Plugin System with Interfaces
Create a plugin system that allows dynamically loading and using modules through a common interface.
This exercise shows how interfaces enable extensible architectures.

Your implementation should include:
1. A `Plugin` interface that defines common methods:
    - `Name()` to identify the plugin
    - `Execute()` to run the plugin's main functionality
    - `Version()` to return the plugin version
2. Several plugin implementations with different behaviors
3. A `PluginManager` that can:
    - Register and unregister plugins
    - Find plugins by name or capability
    - Execute plugins on demand
4. A demonstration showing how new functionality can be added to the system without changing existing code
