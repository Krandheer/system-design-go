package main

import "fmt"

// --- The Component Interface ---
// Pizza is the interface that both our base component and our decorators will implement.
type Pizza interface {
	GetPrice() int
	GetDescription() string
}

// --- The Concrete Component ---
// PlainPizza is our basic, undecorated component.
type PlainPizza struct{}

func (p *PlainPizza) GetPrice() int {
	return 10 // Base price of a plain pizza
}

func (p *PlainPizza) GetDescription() string {
	return "Plain Pizza"
}

// --- The Concrete Decorators ---

// CheeseTopping is a decorator.
type CheeseTopping struct {
	// It wraps a component that conforms to the Pizza interface.
	pizza Pizza
}

func (c *CheeseTopping) GetPrice() int {
	// It adds its own price to the price of the wrapped component.
	wrappedPrice := c.pizza.GetPrice()
	return wrappedPrice + 3 // Price of cheese
}

func (c *CheeseTopping) GetDescription() string {
	// It adds its own description to the description of the wrapped component.
	wrappedDesc := c.pizza.GetDescription()
	return wrappedDesc + ", with Cheese"
}

// TomatoTopping is another decorator.
type TomatoTopping struct {
	pizza Pizza
}

func (t *TomatoTopping) GetPrice() int {
	wrappedPrice := t.pizza.GetPrice()
	return wrappedPrice + 2 // Price of tomato
}

func (t *TomatoTopping) GetDescription() string {
	wrappedDesc := t.pizza.GetDescription()
	return wrappedDesc + ", with Tomato"
}

func main() {
	// Start with a base component.
	pizza := &PlainPizza{}
	printPizzaDetails(pizza)

	// Now, let's decorate it.
	// 1. Wrap the plain pizza with a cheese topping.
	pizzaWithCheese := &CheeseTopping{pizza: pizza}
	printPizzaDetails(pizzaWithCheese)
	
	// 2. Wrap the already-decorated pizza with another topping.
	// This shows how decorators can be stacked.
	pizzaWithCheeseAndTomato := &TomatoTopping{pizza: pizzaWithCheese}
	printPizzaDetails(pizzaWithCheeseAndTomato)
}

func printPizzaDetails(p Pizza) {
	fmt.Printf("Description: %s\n", p.GetDescription())
	fmt.Printf("Price: $%d\n---\n", p.GetPrice())
}
