package main

import "fmt"

// --- The Interfaces ---

// Observer defines the interface for objects that should be notified of changes.
type Observer interface {
	Update(news string)
	ID() string // An ID to help with deregistration
}

// Subject defines the interface for the object that holds the state and notifies observers.
type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	NotifyAll()
}

// --- The Concrete Subject ---

// NewsAgency is our concrete subject. It maintains a list of observers
// and sends them news updates.
type NewsAgency struct {
	observers []Observer
	news      string
}

func (na *NewsAgency) Register(observer Observer) {
	fmt.Printf("Registering observer %s\n", observer.ID())
	na.observers = append(na.observers, observer)
}

func (na *NewsAgency) Deregister(observer Observer) {
	fmt.Printf("Deregistering observer %s\n", observer.ID())
	for i, obs := range na.observers {
		if obs.ID() == observer.ID() {
			// Remove the observer from the slice
			na.observers = append(na.observers[:i], na.observers[i+1:]...)
			return
		}
	}
}

func (na *NewsAgency) NotifyAll() {
	fmt.Println("\n--- News Agency: Notifying all observers of new article ---")
	for _, observer := range na.observers {
		observer.Update(na.news)
	}
}

// SetNews is a method to update the state of the agency. When the news
// changes, it notifies all registered observers.
func (na *NewsAgency) SetNews(news string) {
	na.news = news
	na.NotifyAll()
}

// --- The Concrete Observer ---

// NewsSubscriber is a concrete observer.
type NewsSubscriber struct {
	subscriberID string
}

func (ns *NewsSubscriber) Update(news string) {
	fmt.Printf("Subscriber %s received news: '%s'\n", ns.subscriberID, news)
}

func (ns *NewsSubscriber) ID() string {
	return ns.subscriberID
}

func main() {
	// Create the subject.
	agency := &NewsAgency{}

	// Create some observers.
	sub1 := &NewsSubscriber{subscriberID: "Reader A"}
	sub2 := &NewsSubscriber{subscriberID: "Reader B"}
	sub3 := &NewsSubscriber{subscriberID: "Reader C"}

	// Register the observers with the subject.
	agency.Register(sub1)
	agency.Register(sub2)
	agency.Register(sub3)

	// The subject's state changes. It publishes new content.
	// All registered observers are automatically notified.
	agency.SetNews("Go 1.29 Released with Exciting New Features!")

	// One observer unsubscribes.
	agency.Deregister(sub2)

	// The subject's state changes again.
	// Only the remaining observers are notified.
	agency.SetNews("Design Patterns Prove Useful in Modern Software Engineering!")
}
