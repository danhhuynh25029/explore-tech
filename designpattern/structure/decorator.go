package main

import "fmt"

type Notification interface {
	Send(msg string)
}

type TelegramNotification struct {
}

func (t *TelegramNotification) Send(msg string) {
	fmt.Println("telegram notification send to " + msg)
}

type SMSNotification struct {
}

func (s *SMSNotification) Send(msg string) {
	fmt.Println("sms notification send to " + msg)
}

type NotificationDecorator struct {
	core         *NotificationDecorator
	notification Notification
}

func (n NotificationDecorator) Send(msg string) {
	n.notification.Send(msg)

	if n.core != nil {
		n.core.Send(msg)
	}
}

func (n NotificationDecorator) Decorate(notification Notification) NotificationDecorator {
	return NotificationDecorator{
		core:         &n,
		notification: notification,
	}
}

func NewNotificationDecorator(n Notification) NotificationDecorator {
	return NotificationDecorator{notification: n}
}

type Service struct {
	notification Notification
}

func (s Service) SendNotification(msg string) {
	s.notification.Send(msg)
}

func main() {
	noti := NewNotificationDecorator(&TelegramNotification{}).Decorate(&SMSNotification{})

	service := Service{
		notification: noti,
	}

	service.SendNotification("Hello World")
}
