package main

import "fmt"

// --- The Strategy Interface ---
// RouteStrategy defines the interface for our family of algorithms.
// Any type that can build a route will satisfy this interface.
type RouteStrategy interface {
	BuildRoute(start, end string)
}

// --- Concrete Strategies ---

// WalkingStrategy is a concrete implementation of RouteStrategy.
type WalkingStrategy struct{}

func (s *WalkingStrategy) BuildRoute(start, end string) {
	fmt.Printf("Building a walking route from %s to %s.\n", start, end)
	fmt.Println("Route: Take the scenic path through the park.")
}

// DrivingStrategy is another concrete implementation.
type DrivingStrategy struct{}

func (s *DrivingStrategy) BuildRoute(start, end string) {
	fmt.Printf("Building a driving route from %s to %s.\n", start, end)
	fmt.Println("Route: Take the I-5 highway.")
}

// BikingStrategy is a third implementation.
type BikingStrategy struct{}

func (s *BikingStrategy) BuildRoute(start, end string) {
	fmt.Printf("Building a biking route from %s to %s.\n", start, end)
	fmt.Println("Route: Use the dedicated bike lane on Main St.")
}

// --- The Context ---
// The Navigator is the context. It is configured with a strategy object.
type Navigator struct {
	strategy RouteStrategy
}

// SetStrategy allows the client to change the strategy at runtime.
func (n *Navigator) SetStrategy(strategy RouteStrategy) {
	n.strategy = strategy
}

// Navigate executes the algorithm defined by the current strategy.
// The Navigator itself doesn't know the details of the algorithm;
// it just delegates the work to the strategy object.
func (n *Navigator) Navigate(start, end string) {
	fmt.Println("--- Navigator starting route calculation ---")
	if n.strategy == nil {
		fmt.Println("Error: No strategy has been set.")
		return
	}
	n.strategy.BuildRoute(start, end)
	fmt.Println("--- Route calculation finished ---")
	fmt.Println()
}

func main() {
	// Create the context object.
	navigator := &Navigator{}

	startPoint := "Home"
	endPoint := "Work"

	// Let's go for a walk first.
	// We configure the navigator with the walking strategy.
	navigator.SetStrategy(&WalkingStrategy{})
	navigator.Navigate(startPoint, endPoint)

	// It's raining! Let's switch to the driving strategy.
	// The client can change the behavior of the navigator at runtime.
	fmt.Println("It started raining! Changing strategy to driving...")
	navigator.SetStrategy(&DrivingStrategy{})
	navigator.Navigate(startPoint, endPoint)
	
	// The weather is nice again, let's bike.
	fmt.Println("The weather is nice again! Changing strategy to biking...")
	navigator.SetStrategy(&BikingStrategy{})
	navigator.Navigate(startPoint, endPoint)
}
