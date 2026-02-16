// Package employee defines a type and methods for employees.
package employee

import "fmt"

// Employee defines the data for an employee.
// Note the capitalization of the fields.
type Employee struct {
	Name  string // Exported field
	Title string // Exported field

	// salary is unexported because it starts with a lowercase letter.
	// It can only be accessed by code within the 'employee' package.
	salary float64
}

// New is a "constructor" function. It's the idiomatic way to create
// an instance of a struct that has unexported fields. Since it's
// part of the 'employee' package, it is allowed to access the
// unexported 'salary' field.
func New(name, title string, salary float64) (*Employee, error) {
	if salary < 0 {
		return nil, fmt.Errorf("salary cannot be negative")
	}
	return &Employee{
		Name:   name,
		Title:  title,
		salary: salary,
	}, nil
}

// SetSalary is a "setter" method. It allows code from other packages
// to modify the unexported 'salary' field, but it enforces a business
// rule (salary can't be negative).
func (e *Employee) SetSalary(salary float64) error {
	if salary < 0 {
		return fmt.Errorf("salary cannot be negative")
	}
	e.salary = salary
	return nil
}

// GetSalary is a "getter" method. It allows code from other packages
// to read the value of the unexported 'salary' field.
func (e *Employee) GetSalary() float64 {
	return e.salary
}
