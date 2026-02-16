package main

import (
	"fmt"
	"sync"
	"time"
)

// --- Domain Events ---

// EventType is a unique name for an event.
type EventType string

const (
	UserCreated EventType = "UserCreated"
	OrderPlaced EventType = "OrderPlaced"
)

// Event represents something that happened in the past.
type Event struct {
	Type      EventType
	Data      interface{} // Arbitrary payload
	Timestamp time.Time
}

// --- The Event Bus (Infrastructure) ---

// EventHandler is a function that processes an event.
type EventHandler func(event Event)

// EventBus coordinates the publishing and subscribing of events.
type EventBus struct {
	subscribers map[EventType][]EventHandler
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]EventHandler),
	}
}

// Subscribe allows a service to listen for a specific event type.
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Publish sends an event to all subscribers of that type.
// In a real system, this would likely be asynchronous (using channels or a queue).
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.subscribers[event.Type]
	if !ok {
		return
	}

	fmt.Printf("\n[EventBus] Publishing event: %s\n", event.Type)
	for _, handler := range handlers {
		// Launch each handler in a separate goroutine to simulate async processing
		go handler(event)
	}
}

// --- Services (Subscribers) ---

// EmailService sends emails when users are created.
func EmailService(event Event) {
	username := event.Data.(string)
	fmt.Printf("   -> [EmailService] Sending welcome email to %s...\n", username)
}

// AnalyticsService tracks stats when users are created.
func AnalyticsService(event Event) {
	username := event.Data.(string)
	fmt.Printf("   -> [AnalyticsService] Incrementing daily sign-up counter for %s.\n", username)
}

// MarketingService adds users to a mailing list.
func MarketingService(event Event) {
	username := event.Data.(string)
	fmt.Printf("   -> [MarketingService] Adding %s to newsletter.\n", username)
}

func main() {
	// 1. Initialize the Event Bus
	bus := NewEventBus()

	// 2. Register Subscribers (The "Wiring")
	// Notice how we are wiring unrelated services together via the bus.
	bus.Subscribe(UserCreated, EmailService)
	bus.Subscribe(UserCreated, AnalyticsService)
	bus.Subscribe(UserCreated, MarketingService)

	// 3. Simulate a "Command" (User Registration)
	// The UserService does its job (creates the user) and then just says "I'm done".
	// It doesn't know about Email, Analytics, or Marketing.
	fmt.Println("--- User Registration Flow ---")
	newUser := "Alice"
	fmt.Printf("UserService: Created user '%s' in DB.\n", newUser)

	// 4. Publish Event
	event := Event{
		Type:      UserCreated,
		Data:      newUser,
		Timestamp: time.Now(),
	}
	bus.Publish(event)

	// Give the async handlers a moment to finish (since we used `go handler()`)
	time.Sleep(100 * time.Millisecond)
}
