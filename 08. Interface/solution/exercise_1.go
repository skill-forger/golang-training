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
