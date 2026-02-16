# The Tao of Go: Mastering OOP Principles in a Non-OOP Language

**Target Audience:** Software Engineers transitioning from Java/Python/C++ to Go who want to write *idiomatic* code. This guide explains not just *how*, but *why* Go does things differently, defining core OOP concepts from first principles so you don't need to Google them.

---

## 1. The Philosophy: Data and Behavior are Separate

In classic Object-Oriented Programming (OOP), a **class** is a blueprint that bundles **data (attributes/fields)** and **behavior (methods)** that acts on that data into a single unit.

**Go takes a different approach.**
*   **Structs** are for Data(attributes/fields). They are purely state.
*   **Functions** are for Behavior.
*   **Methods** are just functions with a specific "receiver" (context).

This separation forces you to think about your data structures first (memory layout) and then attach behavior later.

### Implementation: The Employee "Class"

Instead of a class, we define a struct and attach methods to it.

```go
package main

import "fmt"

// 1. DATA: Pure state. No logic here.
// This is just a memory layout.
type Employee struct {
    Name   string
    Title  string
    Salary float64
}

// 2. BEHAVIOR: Attached via the Receiver (e *Employee).
// The receiver appears between `func` and the function name.
// It tells Go: "This function belongs to the Employee type."
//
// We use a Pointer Receiver (*Employee) because we want to modify the original struct.
// If we used a Value Receiver (e Employee), Go would create a COPY of the employee,
// we would give the raise to the copy, and the original employee would remain unchanged.
func (e *Employee) GiveRaise(percent float64) {
    e.Salary = e.Salary * (1 + (percent / 100.0))
}

func (e *Employee) Display() {
    fmt.Printf("Name: %s | Title: %s | Salary: $%.2f\n", e.Name, e.Title, e.Salary)
}

func main() {
    // No "new" keyword. Just struct initialization.
    emp := Employee{
        Name:   "Alice",
        Title:  "Engineer",
        Salary: 100000,
    }

    emp.GiveRaise(10) // Modifies the struct in place
    emp.Display()
}
```

---

## 2. Polymorphism: Interfaces are Implicit

**Definition:** Polymorphism is the ability of different types to be treated as the same type. For example, treating a `Dog` and a `Cat` both as an `Animal`.

In Java, you must explicitly say `class Rectangle implements Shape`. This creates a rigid dependency. The `Rectangle` file *must* import the `Shape` file.

In Go, **Interfaces are satisfied implicitly**. If you can do the job, you get the job. You don't need to sign a contract beforehand.

### Deep Dive: "Define Interfaces Where You Use Them"

This is a critical Go idiom.
*   **Java:** The library defines the interface (`List`), and you implement it (`ArrayList`).
*   **Go:** The *consumer* defines the interface.

**Example:**
Imagine you are writing a function that saves a user.
You don't need a `Database` object. You just need "something that can save".

```go
// defined in YOUR code (the consumer)
type Saver interface {
    Save(user User) error
}

// any type that has implement save method can be passed as saver here.
func CreateUser(s Saver, u User) {
    s.Save(u)
}
```

Now, you can pass a `PostgresDB`, a `RedisCache`, or a `MockSaver` (for testing) into `CreateUser`. As long as they have a `Save` method, it works. The `PostgresDB` type doesn't even need to know your `Saver` interface exists!

### Implementation: The Shape Interface

```go
package main

import (
    "fmt"
    "math"
)

// The Contract: Anyone who has an Area() method is a Shape.
type Shape interface {
    Area() float64
}

// Concrete Type 1
type Rectangle struct {
    Width, Height float64
}

// Implicitly satisfies Shape because it has the Area() method.
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Concrete Type 2
type Circle struct {
    Radius float64
}

// Implicitly satisfies Shape because it has the Area() method.
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Polymorphic Function: Accepts ANY Shape.
// It doesn't know about Rectangles or Circles. It only knows "Area()".
func PrintArea(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
}

func main() {
    r := Rectangle{Width: 10, Height: 5}
    c := Circle{Radius: 7}

    // Both work!
    PrintArea(r)
    PrintArea(c)
}
```

---

## 3. Composition over Inheritance

**Definition:**
*   **Inheritance (Is-a):** A `Car` *is a* `Vehicle`. It inherits all properties of a vehicle. This leads to deep, fragile hierarchies (e.g., `Banana` inherits from `Fruit` which inherits from `Food` which inherits from `Object`). If you change `Food`, you break `Banana`.
*   **Composition (Has-a):** A `Car` *has an* `Engine`. It is built by assembling smaller, independent parts.

Go does **not** have inheritance. There is no `extends` keyword.
Instead, Go uses **Embedding** (Syntactic sugar for Composition).

When you embed a struct, its fields and methods are "promoted" to the outer struct. It *looks* like inheritance (you can call the methods directly), but it is actually composition.

### Implementation: Manager "Has-a" Worker

```go
package main

import "fmt"

type Worker struct {
    Name string
}

func (w *Worker) Work() {
    fmt.Printf("%s is working hard.\n", w.Name)
}

type Manager struct {
    // Anonymous Embedding.
    // We don't give it a name like 'workerField Worker'.
    // We just say 'Worker'.
    // This tells Go: "Manager has a Worker, and Manager can do everything a Worker can do."
    Worker 
    TeamSize int
}

func (m *Manager) Manage() {
    fmt.Printf("%s is managing %d people.\n", m.Name, m.TeamSize)
}

func main() {
    m := Manager{
        Worker:   Worker{Name: "Bob"},
        TeamSize: 5,
    }

    // We can call Work() directly on Manager!
    // Go automatically forwards the call to m.Worker.Work()
    m.Work()   
    m.Manage()
}
```

---

## 4. Encapsulation: The Capitalization Rule

**Definition:** Encapsulation is the practice of bundling data with methods that operate on that data, and restricting direct access to some of an object's components (hiding the internal state).

Go keeps it simple. No `public`, `private`, `protected` keywords.

*   **Capitalized (e.g., `User`):** Exported (Public). Visible to other packages.
*   **Lowercase (e.g., `user`):** Unexported (Private). Visible **only** within the same package.

This applies to everything: structs, fields, functions, methods, and constants.

### Implementation: The Protected Salary

```go
// Package employee
package employee

import "errors"

type Employee struct {
    Name  string // Public: Anyone can read/write this
    Title string // Public
    
    // Private! Cannot be accessed directly from main.go
    // This forces users to use our Getter/Setter methods.
    salary float64 
}

// Constructor (Factory)
// Since 'salary' is private, we need a public function to create the object.
func New(name, title string, salary float64) *Employee {
    return &Employee{
        Name:   name,
        Title:  title,
        salary: salary, // We can access it here because we are in the 'employee' package
    }
}

// Getter: Controlled Read Access
func (e *Employee) GetSalary() float64 {
    return e.salary
}

// Setter: Controlled Write Access with Validation logic
func (e *Employee) SetSalary(amount float64) error {
    if amount < 0 {
        return errors.New("salary cannot be negative")
    }
    e.salary = amount
    return nil
}
```

---

## 5. Error Handling: Errors are Values

Go does not have Exceptions (`try/catch`).
*   **Exceptions:** Implicit control flow. You throw an exception here, and it might be caught 10 functions up the stack. It's hard to trace.
*   **Go Errors:** Explicit values. An error is just a value, like an integer or a string. You must handle it immediately, or explicitly return it to the caller.

This makes Go code verbose but incredibly **robust** and **readable**. You always know exactly where an error can happen.

### Implementation: Custom Errors and Wrapping

```go
package main

import (
    "errors"
    "fmt"
)

// Custom Error Type
// By implementing the Error() string method, this struct satisfies the 'error' interface.
type NotFoundError struct {
    ID string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("resource %s not found", e.ID)
}

func FindUser(id string) (string, error) {
    if id != "123" {
        // Return the custom error
        return "", &NotFoundError{ID: id}
    }
    return "Alice", nil
}

func main() {
    user, err := FindUser("999")
    if err != nil {
        // Check if it's a specific type of error using errors.As
        // This is like `catch (NotFoundError e)` in Java.
        var notFound *NotFoundError
        if errors.As(err, &notFound) {
            fmt.Printf("Specific handling: User %s was missing.\n", notFound.ID)
        } else {
            fmt.Println("Generic error:", err)
        }
        return
    }
    fmt.Println("Found:", user)
}
```

---

**Summary:**
Go achieves OOP goals through simplicity.
1.  **Encapsulation** -> Packages & Capitalization. (https://gemini.google.com/app/b555c2263a48d0ac)
2.  **Abstraction** -> Interfaces.
3.  **Polymorphism** -> Interfaces (Implicit).
4.  **Inheritance** -> Replaced by Composition (Embedding).
