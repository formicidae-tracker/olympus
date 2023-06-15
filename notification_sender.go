package main

import "errors"

type NotificationFor struct {
	Endpoint string
	Updates  []ZonedAlarmUpdate
}

type NotificationSender interface {
	Send(NotificationFor) error
}

type discardNotification struct{}

type logNotification struct {
	Logs []NotificationFor
}

type firebaseSender struct{}

func NewNotificationLogger() *logNotification {
	return &logNotification{}
}

func (l *logNotification) Send(n NotificationFor) error {
	l.Logs = append(l.Logs, n)
	return nil
}

func (l *discardNotification) Send(NotificationFor) error {
	return nil
}

func (l *firebaseSender) Send(NotificationFor) error {
	return errors.New("Not Yet Implemented")
}
