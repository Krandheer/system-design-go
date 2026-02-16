package main

import (
	"fmt"
	"sync"

	"krandheer.github.com/low-level-design/02-design-patterns/01-singleton/settings"
)

func main() {
	// We'll use a WaitGroup to make the main function wait for all
	// goroutines to finish their work.
	var wg sync.WaitGroup
	wg.Add(5) // We're launching 5 goroutines.

	fmt.Println("Attempting to get settings instance from 5 concurrent goroutines.")

	for i := range 5 {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d: getting instance...\n", id)

			// Each goroutine calls GetInstance().
			// Despite these concurrent calls, the "Initializing..." message
			// from the settings package will only ever be printed once.
			s := settings.GetInstance()

			// We can prove they all get the exact same instance by printing
			// its memory address and one of its values.
			fmt.Printf("Goroutine %d: Instance received. Address: %p, Port: %d\n", id, s, s.GetPort())
		}(i)
	}

	// Wait for all goroutines to complete.
	wg.Wait()

	fmt.Println("\n--- All goroutines finished ---")

	// Let's get the instance one more time from the main function.
	// The initialization will not happen again.
	fmt.Println("Getting instance from main function...")
	finalSettings := settings.GetInstance()
	fmt.Printf("Main function: Instance received. Address: %p, Port: %d\n", finalSettings, finalSettings.GetPort())

	fmt.Println("\nNotice that the 'Initializing...' message appeared only once,")
	fmt.Println("and all goroutines received a pointer to the same memory address.")
}
