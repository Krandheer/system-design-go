# Asynchronous Architecture: Decoupling with Queues and Events

**Target Audience:** Engineers building systems that need to handle bursts of traffic without crashing, or who want to decouple complex microservices.

---

## 1. Message Queues (The Buffer)

### The Concept
Synchronous systems are brittle. If Service A calls Service B, and B is slow, A becomes slow. If B crashes, A crashes.
**Message Queues (MQ)** introduce a buffer.
*   **Producer:** Sends a message ("Do this task") to the Queue. Returns immediately.
*   **Queue:** Holds the message safely until a worker is free.
*   **Consumer:** Pulls the message and processes it.

**Benefits:**
1.  **Peak Load Handling:** If 1000 requests come in at once, they just sit in the queue. The consumers process them at a steady pace (e.g., 10/sec). The system doesn't crash.
2.  **Decoupling:** The Producer doesn't need to know who the Consumer is, or if it's even online.

### Implementation: Buffered Channels as a Queue
In Go, a buffered channel is a high-performance, in-memory message queue.

```go
package main

import (
    "fmt"
    "time"
)

type Job struct {
    ID      int
    Payload string
}

func worker(id int, queue <-chan Job) {
    for job := range queue {
        fmt.Printf("[Worker %d] Processing Job %d: %s\n", id, job.ID, job.Payload)
        time.Sleep(500 * time.Millisecond) // Simulate slow work
    }
}

func main() {
    // 1. Create a Queue (Buffered Channel)
    // Capacity = 100. If 101 jobs come, the producer blocks until a slot opens.
    queue := make(chan Job, 100)

    // 2. Start Consumers
    for i := 1; i <= 3; i++ {
        go worker(i, queue)
    }

    // 3. Start Producer (Burst)
    fmt.Println("--- Bursting 10 Jobs ---")
    for i := 1; i <= 10; i++ {
        queue <- Job{ID: i, Payload: "Resize Image"}
        fmt.Printf("[Producer] Enqueued Job %d\n", i)
    }
    
    // Producer is done instantly! It doesn't wait for workers.
    fmt.Println("--- Producer Done ---")

    close(queue)
    time.Sleep(3 * time.Second) // Wait for workers to finish
}
```

---

## 2. Event-Driven Architecture (EDA)

### The Concept
In a Monolith, `UserService` might call `EmailService.SendWelcome()`. This is **Coupling**. If you want to add `AnalyticsService`, you have to modify `UserService`.

In **EDA**, services broadcast **Events**: "User Created".
Other services **Subscribe** to events they care about.

*   **Publisher:** "Something happened!" (Fire and Forget)
*   **Subscriber:** "I care about that!" (Reacts)

**Benefits:**
1.  **Extensibility:** Add new features (subscribers) without touching existing code.
2.  **Resilience:** If the Email service is down, the User service still succeeds. The event waits in the bus.

### Implementation: In-Memory Event Bus

```go
package main

import (
    "fmt"
    "sync"
)

// Event Types
const (
    UserCreated = "UserCreated"
    OrderPlaced = "OrderPlaced"
)

type Event struct {
    Type string
    Data string
}

// The Bus
type EventBus struct {
    subscribers map[string][]func(Event)
    mu          sync.RWMutex
}

func (eb *EventBus) Subscribe(topic string, handler func(Event)) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    eb.subscribers[topic] = append(eb.subscribers[topic], handler)
}

func (eb *EventBus) Publish(e Event) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()
    
    if handlers, found := eb.subscribers[e.Type]; found {
        for _, h := range handlers {
            // Async execution! Don't block the publisher.
            go h(e)
        }
    }
}

// Services
func EmailService(e Event) {
    fmt.Printf("[Email] Sending welcome email to %s\n", e.Data)
}

func AnalyticsService(e Event) {
    fmt.Printf("[Analytics] Tracking sign-up for %s\n", e.Data)
}

func main() {
    bus := &EventBus{subscribers: make(map[string][]func(Event))}

    // Wiring
    bus.Subscribe(UserCreated, EmailService)
    bus.Subscribe(UserCreated, AnalyticsService)

    // Execution
    fmt.Println("--- Registering User ---")
    // UserService just publishes. It doesn't know about Email or Analytics.
    bus.Publish(Event{Type: UserCreated, Data: "alice@example.com"})

    // Wait for async handlers
    // In real code, use WaitGroup or Channels
    fmt.Scanln() 
}
```

---

## 3. Pub/Sub vs. Message Queues

It is crucial to distinguish these two patterns.

| Feature | Message Queue (Point-to-Point) | Pub/Sub (Fan-Out) |
| :--- | :--- | :--- |
| **Goal** | Distribute work. | Broadcast information. |
| **Receivers** | **One** consumer gets the message. | **All** subscribers get the message. |
| **Example** | "Resize this image". Only one worker should do it. | "User Signed Up". Email, Analytics, and Marketing all need to know. |
| **Technology** | RabbitMQ, AWS SQS. | Kafka, Google SNS, Redis Pub/Sub. |

### Implementation: Simulating Both

```go
package main

import "fmt"

func main() {
    // 1. Message Queue (Channels)
    // Only one worker will receive "Task 1".
    queue := make(chan string)
    go func() { fmt.Println("Worker A got:", <-queue) }()
    go func() { fmt.Println("Worker B got:", <-queue) }()
    
    queue <- "Task 1" // Only A OR B will print this. Not both.

    // 2. Pub/Sub (Slice of Channels)
    // Both subscribers will receive "Event 1".
    subA := make(chan string)
    subB := make(chan string)
    
    go func() { fmt.Println("Sub A got:", <-subA) }()
    go func() { fmt.Println("Sub B got:", <-subB) }()

    msg := "Event 1"
    subA <- msg
    subB <- msg
}
```
