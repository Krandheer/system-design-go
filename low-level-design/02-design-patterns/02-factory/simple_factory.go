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
