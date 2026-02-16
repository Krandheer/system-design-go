package main

import (
	"fmt"
)

// Notifier is the interface that all our notification types must satisfy.
// It defines the "what" (sending a notification), but not the "how".
type Notifier interface {
	Send(message string) error
}

// --- Concrete Implementations ---

// EmailNotifier is a concrete type that implements the Notifier interface.
type EmailNotifier struct{}

func (e EmailNotifier) Send(message string) error {
	fmt.Printf("Sending Email: %s\n", message)
	return nil
}

// SMSNotifier is another concrete type that implements the Notifier interface.
type SMSNotifier struct{}

func (s SMSNotifier) Send(message string) error {
	fmt.Printf("Sending SMS: %s\n", message)
	return nil
}

// --- The Factory ---

// GetNotifier is our factory function.
// It takes a string identifier and returns the appropriate concrete Notifier.
// The client code that calls this function doesn't know about EmailNotifier
// or SMSNotifier directly; it only knows it will get something that satisfies
// the Notifier interface.
func GetNotifier(notifierType string) (Notifier, error) {
	switch notifierType {
	case "email":
		return new(EmailNotifier), nil
	case "sms":
		return new(SMSNotifier), nil
	default:
		return nil, fmt.Errorf("unknown notifier type: %s", notifierType)
	}
}

// sendNotification is a helper function that demonstrates the client's perspective.
// It uses the factory to get a notifier and then uses it.
func sendNotification(notifierType, message string) {
	fmt.Printf("--- Attempting to send a '%s' notification ---\n", notifierType)

	notifier, err := GetNotifier(notifierType)
	if err != nil {
		fmt.Printf("Error creating notifier: %v\n", err)
		return
	}

	err = notifier.Send(message)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
	fmt.Println()
}

func main() {
	// The client code wants to send two notifications. It doesn't care about
	// the underlying implementation, only that it gets a "Notifier".

	sendNotification("email", "Hello, this is an email notification.")
	sendNotification("sms", "Hi, this is an SMS.")
	sendNotification("push", "This one will fail.")
}
