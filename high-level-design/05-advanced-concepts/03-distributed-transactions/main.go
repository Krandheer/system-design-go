package main

import (
	"errors"
	"fmt"
)

// Step represents a single transaction in the Saga.
type Step struct {
	Name        string
	Execute     func() error
	Compensate  func()
}

// SagaOrchestrator manages the execution flow.
type SagaOrchestrator struct {
	steps []Step
}

func (s *SagaOrchestrator) AddStep(step Step) {
	s.steps = append(s.steps, step)
}

func (s *SagaOrchestrator) Execute() error {
	fmt.Println("--- Starting Saga Transaction ---")
	
	// Track executed steps so we know what to rollback.
	var executedSteps []Step

	for _, step := range s.steps {
		fmt.Printf("Executing: %s... ", step.Name)
		err := step.Execute()
		
		if err != nil {
			fmt.Println("FAILED!")
			fmt.Printf("Error: %v\n", err)
			s.rollback(executedSteps)
			return err
		}
		
		fmt.Println("Success.")
		executedSteps = append(executedSteps, step)
	}

	fmt.Println("--- Saga Completed Successfully ---")
	return nil
}

func (s *SagaOrchestrator) rollback(executedSteps []Step) {
	fmt.Println("\n--- Initiating Rollback (Compensation) ---")
	// Iterate backwards through executed steps
	for i := len(executedSteps) - 1; i >= 0; i-- {
		step := executedSteps[i]
		fmt.Printf("Compensating: %s... ", step.Name)
		step.Compensate()
		fmt.Println("Done.")
	}
	fmt.Println("--- Rollback Completed ---")
}

func main() {
	saga := &SagaOrchestrator{}

	// Step 1: Create Order
	saga.AddStep(Step{
		Name: "Create Order",
		Execute: func() error {
			// Logic to insert order into DB
			return nil
		},
		Compensate: func() {
			// Logic to update order status to 'CANCELLED'
			fmt.Print("(Order status set to CANCELLED)")
		},
	})

	// Step 2: Reserve Inventory
	saga.AddStep(Step{
		Name: "Reserve Inventory",
		Execute: func() error {
			// Logic to decrement stock
			return nil
		},
		Compensate: func() {
			// Logic to increment stock back
			fmt.Print("(Stock incremented back)")
		},
	})

	// Step 3: Process Payment (Simulate FAILURE)
	saga.AddStep(Step{
		Name: "Process Payment",
		Execute: func() error {
			// Simulate insufficient funds
			return errors.New("insufficient funds")
		},
		Compensate: func() {
			// Logic to refund money (not needed here since Execute failed)
			fmt.Print("(Refund issued)")
		},
	})

	// Run the Saga
	saga.Execute()
}
