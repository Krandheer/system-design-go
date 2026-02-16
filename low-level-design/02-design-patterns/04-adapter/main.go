package main

import "fmt"

// --- The Target Interface ---
// This is the interface our application's client code expects to work with.
type PaymentProcessor interface {
	ProcessPayment(amount float32)
}

// --- An Existing, Compatible Type ---
// PaypalClient is an existing type in our system that already conforms to the target interface.
type PaypalClient struct{}

func (p *PaypalClient) ProcessPayment(amount float32) {
	fmt.Printf("Processing payment of $%.2f through PayPal.\n", amount)
}

// --- The Adaptee (The Incompatible Type) ---
// StripeClient is a new, third-party type we want to integrate.
// Its method signature is different from what our system expects.
type StripeClient struct{}

// Charge is the incompatible method. It takes the amount in cents (an int).
func (s *StripeClient) Charge(amountInCents int) {
	fmt.Printf("Charging %d cents through Stripe.\n", amountInCents)
}

// --- The Adapter ---
// StripeAdapter wraps the StripeClient (the adaptee) and implements
// the PaymentProcessor interface (the target).
type StripeAdapter struct {
	// The adapter holds a reference to the adaptee.
	stripeClient *StripeClient
}

// ProcessPayment is the implementation of the target interface.
// This is where the translation happens.
func (sa *StripeAdapter) ProcessPayment(amount float32) {
	// 1. We get the call in a format our system understands (e.g., $25.50).
	// 2. We translate it to the format the adaptee expects (e.g., 2550 cents).
	amountInCents := int(amount * 100)
	// 3. We call the adaptee's method with the translated data.
	sa.stripeClient.Charge(amountInCents)
}

// process represents the client code in our application.
// It only knows about the PaymentProcessor interface.
func process(p PaymentProcessor, amount float32) {
	fmt.Println("--- Kicking off new payment process ---")
	p.ProcessPayment(amount)
	fmt.Println()
}

func main() {
	// Our client code can process payments without knowing the specifics.
	
	// Using the existing, compatible client is straightforward.
	paypal := &PaypalClient{}
	process(paypal, 25.50)

	// To use the new Stripe client, we wrap it in our adapter first.
	stripeAdaptee := &StripeClient{}
	stripeAdapter := &StripeAdapter{stripeClient: stripeAdaptee}

	// Now we can pass the adapter to our client code, which treats it just like
	// any other PaymentProcessor. The client is completely unaware of the
	// underlying StripeClient and its different method signature.
	process(stripeAdapter, 12.75)
}

