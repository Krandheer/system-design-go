package main
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