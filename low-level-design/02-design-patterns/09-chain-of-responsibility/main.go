package main

import "fmt"

// --- The Handler Interface ---
type Approver interface {
	SetNext(approver Approver)
	ProcessRequest(amount float64)
}

// --- Concrete Handlers ---

// Manager is a concrete handler.
type Manager struct {
	next Approver
}

func (m *Manager) SetNext(approver Approver) {
	m.next = approver
}

func (m *Manager) ProcessRequest(amount float64) {
	// If the manager can approve it, they do.
	if amount <= 500 {
		fmt.Printf("Manager approved expense of $%.2f\n", amount)
		return
	}
	// Otherwise, if there's a next handler in the chain, pass it on.
	if m.next != nil {
		fmt.Println("Manager cannot approve, passing to Director.")
		m.next.ProcessRequest(amount)
	}
}

// Director is another concrete handler.
type Director struct {
	next Approver
}

func (d *Director) SetNext(approver Approver) {
	d.next = approver
}

func (d *Director) ProcessRequest(amount float64) {
	if amount <= 5000 {
		fmt.Printf("Director approved expense of $%.2f\n", amount)
		return
	}
	if d.next != nil {
		fmt.Println("Director cannot approve, passing to Vice President.")
		d.next.ProcessRequest(amount)
	}
}

// VicePresident is the final handler in our chain.
type VicePresident struct {
	next Approver // In our case, this will be nil.
}

func (vp *VicePresident) SetNext(approver Approver) {
	vp.next = approver
}

func (vp *VicePresident) ProcessRequest(amount float64) {
	// The VP can approve any amount.
	fmt.Printf("Vice President approved expense of $%.2f\n", amount)
}

func main() {
	// Create the individual handlers.
	manager := &Manager{}
	director := &Director{}
	vp := &VicePresident{}

	// --- Build the chain of responsibility ---
	// The client sets up the chain. manager -> director -> vp
	manager.SetNext(director)
	director.SetNext(vp)
	
	// The client sends requests to the *start* of the chain (the manager).
	
	fmt.Println("Processing expense of $300...")
	manager.ProcessRequest(300) // Handled by Manager
	fmt.Println()

	fmt.Println("Processing expense of $2500...")
	manager.ProcessRequest(2500) // Passed to and handled by Director
	fmt.Println()

	fmt.Println("Processing expense of $10000...")
	manager.ProcessRequest(10000) // Passed to Director, then handled by VP
	fmt.Println()
}
