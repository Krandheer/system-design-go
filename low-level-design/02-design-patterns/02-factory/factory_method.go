package main

//factory method deals with creating one product at a time.
// abstract factory deals with creating multiple products that must work together.
// we have not implemented abstract factory in this example.
type NotifierFactory interface {
	CreateNotifier() Notifier
}

type EmailNotifierFactory struct{}

func (f *EmailNotifierFactory) CreateNotifier() Notifier {
	return &EmailNotifier{}
}

type SMSNotifierFactory struct{}

func (f *SMSNotifierFactory) CreateNotifier() Notifier {
	return &SMSNotifier{}
}