package main

import (
	"fmt"
	"log"

	"krandheer.github.com/low-level-design/01-thinking-in-go/04-encapsulation-by-packages/employee"
)

func main() {
	// We must use the constructor function New to create an employee,
	// because we cannot set the 'salary' field directly from this package.
	emp, err := employee.New("Alice", "CEO", 300000)
	if err != nil {
		log.Fatal(err)
	}

	// We can access the exported fields directly.
	fmt.Printf("Employee: %s\n", emp.Name)
	fmt.Printf("Title: %s\n", emp.Title)

	// --- THIS LINE WOULD CAUSE A COMPILE ERROR ---
	// fmt.Printf("Salary: %f\n", emp.salary)
	// Error: emp.salary undefined (cannot refer to unexported field or method salary)
	// ---------------------------------------------

	// We must use the exported "getter" method to read the salary.
	fmt.Printf("Salary (via GetSalary): $%.2f\n", emp.GetSalary())

	// Let's try to set a new, valid salary using the "setter" method.
	fmt.Println("\nUpdating salary to $350,000...")
	err = emp.SetSalary(350000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New Salary (via GetSalary): $%.2f\n", emp.GetSalary())

	// Now, let's try to set an invalid salary.
	fmt.Println("\nAttempting to set a negative salary...")
	err = emp.SetSalary(-50000)
	if err != nil {
		fmt.Printf("Caught expected error: %v\n", err)
	}

	fmt.Printf("Salary remains unchanged: $%.2f\n", emp.GetSalary())
}
