package main

import "fmt"

// Employee is a struct that holds data about an employee.
// This is analogous to the attributes of a class in Python/Java.
// Notice it only contains data, no behavior.
type Employee struct {
	Name   string
	Title  string
	Salary float64
}

// GiveRaise is a method associated with the Employee struct.
// The '(e *Employee)' part is called the "receiver". It specifies which
// struct this method "belongs" to.
// This is how Go attaches behavior to data types.
// We use a pointer receiver (*Employee) because we want to modify the
// original Employee struct. If we didn't use a pointer, the method
// would operate on a copy of the Employee, and the salary change
// would be lost.
func (e *Employee) GiveRaise(percent float64) {
	e.Salary = e.Salary * (1 + (percent / 100.0))
}

// Display is another method on Employee to print its details.
// This method doesn't need to modify the Employee, so it can use a
// value receiver (e Employee), but it's idiomatic in Go to keep
// all methods for a given type on pointer receivers for consistency.
func (e *Employee) Display() {
	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Title: %s\n", e.Title)
	fmt.Printf("Salary: $%.2f\n", e.Salary)
}

func main() {
	// In Go, we don't have a formal "constructor".
	// We create an instance by literally constructing the struct.
	emp := Employee{
		Name:   "Alice",
		Title:  "Software Engineer",
		Salary: 100000,
	}

	fmt.Println("Initial State:")
	emp.Display()
	fmt.Println("---")

	// We call the method on the instance, just like in OOP.
	emp.GiveRaise(10) // Give a 10% raise

	fmt.Println("After 10% Raise:")
	emp.Display()
}
