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

// TimerPlugin delay the execution time by specific duration
type TimerPlugin struct{}

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

	manager.RegisterPlugin(TimerPlugin{})

	fmt.Println("\nPlugins after adding new one:")
	for _, name := range manager.ListPlugins() {
		plugin, _ := manager.GetPlugin(name)
		fmt.Printf("- %s (v%s)\n", name, plugin.Version())
	}
}
