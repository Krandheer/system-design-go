# Go Design Patterns: The Idiomatic Way

**Target Audience:** Engineers who know the "Gang of Four" patterns but struggle to implement them cleanly in Go without classes.

This guide covers the most essential design patterns, adapted for Go's unique features like interfaces, channels, and goroutines.

---

## Part 1: Creational Patterns (Object Creation)

### 1. Singleton (Thread-Safe)
**Problem:** Ensure a struct has only one instance and provide a global point of access to it.
**Go Solution:** Use `sync.Once`. It guarantees that initialization code runs exactly once, even if 100 goroutines call it simultaneously.

```go
package main

import (
    "fmt"
    "sync"
)

type Database struct {
    URL string
}

var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    // 'once.Do' is the magic. It uses atomic counters and mutexes internally.
    // It is much safer and faster than writing your own "if instance == nil" check.
    once.Do(func() {
        fmt.Println("Initializing Database...")
        instance = &Database{URL: "postgres://localhost:5432"}
    })
    return instance
}

func main() {
    // Even if we call this multiple times, "Initializing..." prints only once.
    db1 := GetDatabase()
    db2 := GetDatabase()
    fmt.Println(db1 == db2) // true
}
```

### 2. Factory Method
**Problem:** You want to create objects without specifying the exact class of object that will be created.
**Go Solution:** A simple function that returns an Interface.

```go
package main

import "fmt"

// The Interface
type Notifier interface {
    Send(msg string)
}

// Concrete Type 1
type Email struct{}
func (e Email) Send(msg string) { fmt.Println("Sending Email:", msg) }

// Concrete Type 2
type SMS struct{}
func (s SMS) Send(msg string) { fmt.Println("Sending SMS:", msg) }

// The Factory
func NewNotifier(method string) (Notifier, error) {
    switch method {
    case "email":
        return Email{}, nil
    case "sms":
        return SMS{}, nil
    default:
        return nil, fmt.Errorf("unknown method")
    }
}

func main() {
    n, _ := NewNotifier("sms")
    n.Send("Hello Factory")
}
```

### 3. Builder
**Problem:** Constructing a complex object with many optional parameters.
**Go Solution:** A struct with methods that return the builder itself (Fluent Interface).

```go
package main

import "fmt"

type Server struct {
    Host string
    Port int
    TLS  bool
}

type ServerBuilder struct {
    server Server
}

func NewBuilder(host string) *ServerBuilder {
    return &ServerBuilder{server: Server{Host: host, Port: 80}} // Default port
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
    b.server.Port = port
    return b
}

func (b *ServerBuilder) WithTLS() *ServerBuilder {
    b.server.TLS = true
    return b
}

func (b *ServerBuilder) Build() Server {
    return b.server
}

func main() {
    // Clean, readable construction
    s := NewBuilder("localhost").WithPort(8080).WithTLS().Build()
    fmt.Printf("%+v\n", s)
}
```

---

## Part 2: Structural Patterns (Composition)

### 4. Adapter
**Problem:** You have an existing class (Adaptee) with an incompatible interface that you need to use.
**Go Solution:** Create a struct that embeds or holds the Adaptee and implements the Target interface.

```go
package main

import "fmt"

// Target Interface
type PaymentProcessor interface {
    Pay(amount float64)
}

// Adaptee (Third-party library, incompatible)
type Stripe struct{}
func (s *Stripe) ChargeCents(cents int) {
    fmt.Println("Stripe charged cents:", cents)
}

// Adapter
type StripeAdapter struct {
    stripe *Stripe
}

func (a *StripeAdapter) Pay(amount float64) {
    // Convert dollars to cents
    cents := int(amount * 100)
    a.stripe.ChargeCents(cents)
}

func main() {
    // Client code only knows about PaymentProcessor
    var p PaymentProcessor = &StripeAdapter{stripe: &Stripe{}}
    p.Pay(10.50)
}
```

### 5. Decorator
**Problem:** Add behavior to an individual object dynamically without affecting other objects.
**Go Solution:** A struct that wraps an Interface and implements that same Interface.

```go
package main

import "fmt"

// Component
type Pizza interface {
    GetPrice() int
}

// Concrete Component
type VeggiePizza struct{}
func (p VeggiePizza) GetPrice() int { return 15 }

// Decorator
type CheeseTopping struct {
    pizza Pizza
}

func (c CheeseTopping) GetPrice() int {
    return c.pizza.GetPrice() + 5 // Add $5 for cheese
}

func main() {
    pizza := VeggiePizza{}
    fmt.Println("Base:", pizza.GetPrice())

    // Wrap it
    pizzaWithCheese := CheeseTopping{pizza: pizza}
    fmt.Println("With Cheese:", pizzaWithCheese.GetPrice())
}
```

### 6. Facade
**Problem:** Provide a simplified interface to a complex subsystem.
**Go Solution:** A "Wrapper" struct that orchestrates calls to other structs.

```go
package main

import "fmt"

type CPU struct{}
func (c CPU) Freeze() { fmt.Println("CPU Freeze") }
func (c CPU) Jump(addr int) { fmt.Println("CPU Jump to", addr) }

type Memory struct{}
func (m Memory) Load(addr int, data string) { fmt.Println("Memory Load", data) }

// Facade
type Computer struct {
    cpu CPU
    mem Memory
}

func (c Computer) Start() {
    c.cpu.Freeze()
    c.mem.Load(0, "BOOT_LOADER")
    c.cpu.Jump(0)
}

func main() {
    comp := Computer{}
    comp.Start() // Client doesn't need to know about CPU or Memory details
}
```

---

## Part 3: Behavioral Patterns (Communication)

### 7. Strategy
**Problem:** Define a family of algorithms, encapsulate each one, and make them interchangeable.
**Go Solution:** An Interface for the algorithm, and a Context struct that holds the Interface.

```go
package main

import "fmt"

// Strategy Interface
type RouteStrategy interface {
    BuildRoute(a, b string)
}

// Concrete Strategies
type Walk struct{}
func (w Walk) BuildRoute(a, b string) { fmt.Println("Walking route from", a, "to", b) }

type Drive struct{}
func (d Drive) BuildRoute(a, b string) { fmt.Println("Driving route from", a, "to", b) }

// Context
type Navigator struct {
    strategy RouteStrategy
}

func (n *Navigator) SetStrategy(s RouteStrategy) {
    n.strategy = s
}

func (n *Navigator) Navigate(a, b string) {
    n.strategy.BuildRoute(a, b)
}

func main() {
    nav := Navigator{}
    
    nav.SetStrategy(Walk{})
    nav.Navigate("Home", "Park")

    nav.SetStrategy(Drive{})
    nav.Navigate("Home", "Work")
}
```

### 8. Observer
**Problem:** Define a one-to-many dependency so that when one object changes state, all its dependents are notified.
**Go Solution:** A list of Interfaces.

```go
package main

import "fmt"

type Observer interface {
    Update(msg string)
}

type Subscriber struct {
    ID string
}
func (s *Subscriber) Update(msg string) {
    fmt.Printf("[%s] Received: %s\n", s.ID, msg)
}

type Publisher struct {
    subscribers []Observer
}

func (p *Publisher) Subscribe(o Observer) {
    p.subscribers = append(p.subscribers, o)
}

func (p *Publisher) Notify(msg string) {
    for _, sub := range p.subscribers {
        sub.Update(msg)
    }
}

func main() {
    pub := Publisher{}
    pub.Subscribe(&Subscriber{ID: "Alice"})
    pub.Subscribe(&Subscriber{ID: "Bob"})
    
    pub.Notify("New Video Uploaded!")
}
```

### 9. Chain of Responsibility
**Problem:** Pass a request along a chain of handlers.
**Go Solution:** A linked list of Interfaces.

```go
package main

import "fmt"

type Handler interface {
    SetNext(h Handler)
    Handle(request string)
}

type BaseHandler struct {
    next Handler
}

func (h *BaseHandler) SetNext(next Handler) {
    h.next = next
}

type AuthHandler struct { BaseHandler }
func (h *AuthHandler) Handle(req string) {
    if req == "admin" {
        fmt.Println("Auth Passed")
        if h.next != nil { h.next.Handle(req) }
    } else {
        fmt.Println("Auth Failed")
    }
}

type LogHandler struct { BaseHandler }
func (h *LogHandler) Handle(req string) {
    fmt.Println("Request Logged")
    if h.next != nil { h.next.Handle(req) }
}

func main() {
    auth := &AuthHandler{}
    logger := &LogHandler{}
    
    auth.SetNext(logger) // Build chain: Auth -> Logger
    
    auth.Handle("admin")
}
```
