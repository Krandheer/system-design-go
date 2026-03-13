package main

import "fmt"

// sendNotification is a helper function that demonstrates the client's perspective.
// It uses the factory to get a notifier and then uses it.
func sendNotificationSimpleFactory(notifierType, message string) {
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

func sendNotificationFactoryMethod(factory NotifierFactory, message string) {
	notifier := factory.CreateNotifier()
	err := notifier.Send(message)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
	fmt.Println()
}

func main() {
	// The client code wants to send two notifications. It doesn't care about
	// the underlying implementation, only that it gets a "Notifier".

	fmt.Println("--- Simple Factory ---")
	sendNotificationSimpleFactory("email", "Hello, this is an email notification.")
	sendNotificationSimpleFactory("sms", "Hi, this is an SMS.")
	sendNotificationSimpleFactory("push", "This one will fail.")

	fmt.Println("--- Factory Method ---")
	emailFactory := &EmailNotifierFactory{}
	smsFactory := &SMSNotifierFactory{}

	sendNotificationFactoryMethod(emailFactory, "Hello, this is an email notification.")
	sendNotificationFactoryMethod(smsFactory, "Hi, this is an SMS.")
}
