package main

import "fmt"

// Worker is a struct representing someone who can perform work.
type Worker struct {
	Name string
}

// Work is a method on the Worker struct.
func (w *Worker) Work() {
	fmt.Printf("%s is working.\n", w.Name)
}

// Manager is a struct representing a manager.
type Manager struct {
	Worker // This is struct embedding. It's an anonymous field.
	Team   []string
}

// Manage is a method specific to the Manager.
func (m *Manager) Manage() {
	fmt.Printf("%s is managing a team of %d people.\n", m.Name, len(m.Team))
}

func main() {
	// Create a Manager instance.
	m := Manager{
		Worker: Worker{Name: "Alice"}, // We initialize the embedded struct.
		Team:   []string{"Bob", "Charlie", "Dave"},
	}

	// Because we embedded Worker in Manager, the fields and methods of Worker
	// are "promoted" and can be called directly on the Manager instance.

	// We can access the Name field directly from the embedded Worker.
	fmt.Printf("Manager's name is %s\n", m.Name)

	// We can call the Work() method directly, as if it belonged to Manager.
	m.Work()

	// We can also call methods that are specific to the Manager.
	m.Manage()
}
