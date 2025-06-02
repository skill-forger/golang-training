## Practical Exercises

### Exercise 1: Event System with Interface-Based Pub/Sub

Build an event management system using interfaces to implement a publisher-subscriber pattern. This exercise demonstrates how interfaces enable flexible and loosely coupled communication between components.

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

This exercise illustrates how interfaces can create flexible, extensible systems where components interact without tight coupling.

```go
// event_system.go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Define the Event interface
type Event interface {
    Type() string
    Data() interface{}
    Timestamp() time.Time
}

// Concrete implementation of Event
type BaseEvent struct {
    EventType string
    EventData interface{}
    EventTime time.Time
}

func (e BaseEvent) Type() string {
    return e.EventType
}

func (e BaseEvent) Data() interface{} {
    return e.EventData
}

func (e BaseEvent) Timestamp() time.Time {
    return e.EventTime
}

// Define the Handler interface
type EventHandler interface {
    Handle(event Event)
}

// Handler function type for convenience
type EventHandlerFunc func(Event)

// Make EventHandlerFunc implement EventHandler
func (f EventHandlerFunc) Handle(event Event) {
    f(event)
}

// EventBus manages event subscriptions and publishing
type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

// Subscribe registers a handler for a specific event type
func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeFunc is a convenience method for function-based handlers
func (b *EventBus) SubscribeFunc(eventType string, handlerFunc func(Event)) {
    b.Subscribe(eventType, EventHandlerFunc(handlerFunc))
}

// Publish sends an event to all registered handlers
func (b *EventBus) Publish(event Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    
    // Find handlers for this event type
    handlers, exists := b.handlers[event.Type()]
    if !exists {
        return
    }
    
    // Notify all handlers
    for _, handler := range handlers {
        handler.Handle(event)
    }
}

// Example usage
func main() {
    // Create the event bus
    bus := NewEventBus()
    
    // Subscribe to user.created events
    bus.SubscribeFunc("user.created", func(event Event) {
        userData, ok := event.Data().(map[string]string)
        if !ok {
            fmt.Println("Invalid user data format")
            return
        }
        
        fmt.Printf("User created at %v: %s (%s)\n", 
            event.Timestamp().Format("15:04:05"),
            userData["name"], 
            userData["email"])
    })
    
    // Subscribe to payment.received events
    bus.SubscribeFunc("payment.received", func(event Event) {
        amount, ok := event.Data().(float64)
        if !ok {
            fmt.Println("Invalid payment data format")
            return
        }
        
        fmt.Printf("Payment received at %v: $%.2f\n", 
            event.Timestamp().Format("15:04:05"),
            amount)
    })
    
    // Structured logger for all events
    bus.SubscribeFunc("*", func(event Event) {
        fmt.Printf("[LOG] %s event at %v with data: %v\n",
            event.Type(),
            event.Timestamp().Format("15:04:05"),
            event.Data())
    })
    
    // Publish some events
    bus.Publish(BaseEvent{
        EventType: "user.created",
        EventData: map[string]string{
            "name":  "John Doe",
            "email": "john@example.com",
        },
        EventTime: time.Now(),
    })
    
    time.Sleep(1 * time.Second)
    
    bus.Publish(BaseEvent{
        EventType: "payment.received",
        EventData: 125.50,
        EventTime: time.Now(),
    })
    
    time.Sleep(1 * time.Second)
    
    bus.Publish(BaseEvent{
        EventType: "user.updated",
        EventData: map[string]string{
            "id":    "123",
            "name":  "John Updated",
            "email": "john.updated@example.com",
        },
        EventTime: time.Now(),
    })
}
```

### Exercise 2: Shape Calculator with Interface Hierarchy

Create a shape calculation system that uses interfaces to handle different geometric shapes uniformly. This exercise shows how interfaces enable polymorphism in Go.

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

```go
// shape_calculator.go
package main

import (
    "fmt"
    "math"
    "sort"
)

// Shape interface defines methods all shapes must implement
type Shape interface {
    Area() float64
    Perimeter() float64
    Name() string
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

func (c Circle) Name() string {
    return "Circle"
}

// Rectangle implements the Shape interface
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
    return "Rectangle"
}

// Triangle implements the Shape interface
type Triangle struct {
    SideA float64
    SideB float64
    SideC float64
}

func (t Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

func (t Triangle) Area() float64 {
    // Heron's formula
    s := t.Perimeter() / 2
    return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t Triangle) Name() string {
    return "Triangle"
}

// ThreeDimensionalShape extends the Shape interface
type ThreeDimensionalShape interface {
    Shape
    Volume() float64
}

// Sphere implements ThreeDimensionalShape
type Sphere struct {
    Radius float64
}

func (s Sphere) Area() float64 {
    return 4 * math.Pi * s.Radius * s.Radius
}

func (s Sphere) Perimeter() float64 {
    return 2 * math.Pi * s.Radius // Great circle
}

func (s Sphere) Volume() float64 {
    return (4.0 / 3.0) * math.Pi * math.Pow(s.Radius, 3)
}

func (s Sphere) Name() string {
    return "Sphere"
}

// Cube implements ThreeDimensionalShape
type Cube struct {
    Side float64
}

func (c Cube) Area() float64 {
    return 6 * c.Side * c.Side
}

func (c Cube) Perimeter() float64 {
    return 12 * c.Side
}

func (c Cube) Volume() float64 {
    return math.Pow(c.Side, 3)
}

func (c Cube) Name() string {
    return "Cube"
}

// ShapeProcessor provides utility functions for working with shapes
type ShapeProcessor struct{}

// SortByArea sorts shapes by their area
func (sp ShapeProcessor) SortByArea(shapes []Shape) {
    sort.Slice(shapes, func(i, j int) bool {
        return shapes[i].Area() < shapes[j].Area()
    })
}

// PrintShapeInfo displays information about a shape
func (sp ShapeProcessor) PrintShapeInfo(shape Shape) {
    fmt.Printf("%s:\n", shape.Name())
    fmt.Printf("  Area: %.2f\n", shape.Area())
    fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())
    
    // Check if it's also a 3D shape
    if threeDShape, ok := shape.(ThreeDimensionalShape); ok {
        fmt.Printf("  Volume: %.2f\n", threeDShape.Volume())
    }
}

// FilterByType returns shapes of a specific type
func (sp ShapeProcessor) FilterByType(shapes []Shape, typeName string) []Shape {
    var result []Shape
    for _, shape := range shapes {
        if shape.Name() == typeName {
            result = append(result, shape)
        }
    }
    return result
}

func main() {
    // Create various shapes
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 4, Height: 6},
        Triangle{SideA: 3, SideB: 4, SideC: 5},
        Sphere{Radius: 3},
        Cube{Side: 4},
    }
    
    processor := ShapeProcessor{}
    
    // Print information for each shape
    fmt.Println("All Shapes:")
    for _, shape := range shapes {
        processor.PrintShapeInfo(shape)
        fmt.Println()
    }
    
    // Sort shapes by area
    processor.SortByArea(shapes)
    fmt.Println("Shapes sorted by area:")
    for _, shape := range shapes {
        fmt.Printf("%s: %.2f\n", shape.Name(), shape.Area())
    }
    
    // Filter 3D shapes
    var threeDShapes []ThreeDimensionalShape
    for _, shape := range shapes {
        if threeDShape, ok := shape.(ThreeDimensionalShape); ok {
            threeDShapes = append(threeDShapes, threeDShape)
        }
    }
    
    fmt.Println("\nThree-dimensional shapes:")
    for _, shape := range threeDShapes {
        fmt.Printf("%s - Volume: %.2f\n", shape.Name(), shape.Volume())
    }
}
```

### Exercise 3: Plugin System with Interfaces

Create a plugin system that allows dynamically loading and using modules through a common interface. This exercise shows how interfaces enable extensible architectures.

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

```go
// plugin_system.go
package main

import (
    "fmt"
    "time"
)

// Plugin interface defines the required methods for all plugins
type Plugin interface {
    Name() string
    Execute(data map[string]interface{}) (interface{}, error)
    Version() string
}

// LoggerPlugin implements a simple logging plugin
type LoggerPlugin struct {
    logLevel string
}

func (p LoggerPlugin) Name() string {
    return "Logger"
}

func (p LoggerPlugin) Execute(data map[string]interface{}) (interface{}, error) {
    message, ok := data["message"].(string)
    if !ok {
        return nil, fmt.Errorf("message is required and must be a string")
    }
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, p.logLevel, message)
    
    fmt.Println(logEntry)
    return logEntry, nil
}

func (p LoggerPlugin) Version() string {
    return "1.0.0"
}

// CalculatorPlugin implements basic math operations
type CalculatorPlugin struct{}

func (p CalculatorPlugin) Name() string {
    return "Calculator"
}

func (p CalculatorPlugin) Execute(data map[string]interface{}) (interface{}, error) {
    operation, ok := data["operation"].(string)
    if !ok {
        return nil, fmt.Errorf("operation is required and must be a string")
    }
    
    a, aOK := data["a"].(float64)
    b, bOK := data["b"].(float64)
    
    if !aOK || !bOK {
        return nil, fmt.Errorf("a and b are required and must be numbers")
    }
    
    switch operation {
    case "add":
        return a + b, nil
    case "subtract":
        return a - b, nil
    case "multiply":
        return a * b, nil
    case "divide":
        if b == 0 {
            return nil, fmt.Errorf("division by zero")
        }
        return a / b, nil
    default:
        return nil, fmt.Errorf("unsupported operation: %s", operation)
    }
}

func (p CalculatorPlugin) Version() string {
    return "1.0.0"
}

// FormatterPlugin formats different data types
type FormatterPlugin struct{}

func (p FormatterPlugin) Name() string {
    return "Formatter"
}

func (p FormatterPlugin) Execute(data map[string]interface{}) (interface{}, error) {
    format, ok := data["format"].(string)
    if !ok {
        return nil, fmt.Errorf("format is required and must be a string")
    }
    
    value := data["value"]
    
    switch format {
    case "uppercase":
        str, ok := value.(string)
        if !ok {
            return nil, fmt.Errorf("value must be a string for uppercase format")
        }
        return fmt.Sprintf("%s", str), nil
    case "json":
        // In a real implementation, this would convert to JSON
        return fmt.Sprintf("%v", value), nil
    case "date":
        timeVal, ok := value.(time.Time)
        if !ok {
            return nil, fmt.Errorf("value must be a time.Time for date format")
        }
        return timeVal.Format("2006-01-02"), nil
    default:
        return nil, fmt.Errorf("unsupported format: %s", format)
    }
}

func (p FormatterPlugin) Version() string {
    return "1.0.0"
}

// PluginManager handles registration and execution of plugins
type PluginManager struct {
    plugins map[string]Plugin
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]Plugin),
    }
}

// RegisterPlugin adds a plugin to the manager
func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
    pm.plugins[plugin.Name()] = plugin
}

// UnregisterPlugin removes a plugin from the manager
func (pm *PluginManager) UnregisterPlugin(name string) {
    delete(pm.plugins, name)
}

// GetPlugin retrieves a plugin by name
func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
    plugin, exists := pm.plugins[name]
    return plugin, exists
}

// ExecutePlugin runs a plugin by name with the provided data
func (pm *PluginManager) ExecutePlugin(name string, data map[string]interface{}) (interface{}, error) {
    plugin, exists := pm.plugins[name]
    if !exists {
        return nil, fmt.Errorf("plugin '%s' not found", name)
    }
    
    return plugin.Execute(data)
}

// ListPlugins returns the names of all registered plugins
func (pm *PluginManager) ListPlugins() []string {
    var names []string
    for name := range pm.plugins {
        names = append(names, name)
    }
    return names
}

func main() {
    // Create a plugin manager
    manager := NewPluginManager()
    
    // Register plugins
    manager.RegisterPlugin(LoggerPlugin{logLevel: "INFO"})
    manager.RegisterPlugin(CalculatorPlugin{})
    manager.RegisterPlugin(FormatterPlugin{})
    
    // List available plugins
    fmt.Println("Available plugins:")
    for _, name := range manager.ListPlugins() {
        plugin, _ := manager.GetPlugin(name)
        fmt.Printf("- %s (v%s)\n", name, plugin.Version())
    }
    
    fmt.Println("\nExecuting plugins:")
    
    // Execute logger plugin
    result, err := manager.ExecutePlugin("Logger", map[string]interface{}{
        "message": "This is a test log message",
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Logger result: %v\n", result)
    }
    
    // Execute calculator plugin
    result, err = manager.ExecutePlugin("Calculator", map[string]interface{}{
        "operation": "add",
        "a":         10.5,
        "b":         5.2,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Calculator result: %.1f\n", result)
    }
    
    // Handle an error case
    result, err = manager.ExecutePlugin("Calculator", map[string]interface{}{
        "operation": "divide",
        "a":         10.0,
        "b":         0.0,
    })
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
    }
    
    // Add a new plugin at runtime
    type TimerPlugin struct{}
    
    manager.RegisterPlugin(TimerPlugin{})
    
    fmt.Println("\nPlugins after adding new one:")
    for _, name := range manager.ListPlugins() {
        plugin, _ := manager.GetPlugin(name)
        fmt.Printf("- %s (v%s)\n", name, plugin.Version())
    }
}

// Implementation of the TimerPlugin methods would be here in a real example
func (p TimerPlugin) Name() string {
    return "Timer"
}

func (p TimerPlugin) Execute(data map[string]interface{}) (interface{}, error) {
    duration, ok := data["duration"].(int)
    if !ok {
        return nil, fmt.Errorf("duration is required and must be an integer")
    }
    
    start := time.Now()
    time.Sleep(time.Duration(duration) * time.Millisecond)
    elapsed := time.Since(start)
    
    return elapsed.Milliseconds(), nil
}

func (p TimerPlugin) Version() string {
    return "1.0.0"
}
```
