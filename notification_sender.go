package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

type NotificationFor struct {
	Subscription *webpush.Subscription
	Updates      []ZonedAlarmUpdate
}

type NotificationSender interface {
	Send(NotificationFor) error
}

type logNotification struct {
	Logs []NotificationFor
}

func NewNotificationLogger() *logNotification {
	return &logNotification{}
}

func (l *logNotification) Send(n NotificationFor) error {
	l.Logs = append(l.Logs, n)
	return nil
}

type discardNotification struct{}

type webpushSender struct {
	public, private, subscriber string
	ctx                         context.Context
}

func NewNotificationSender() NotificationSender {
	private := os.Getenv("OLYMPUS_VAPID_PRIVATE")
	public := os.Getenv("OLYMPUS_VAPID_PUBLIC")
	subscriber := os.Getenv("OLYMPUS_PUSH_SUBSCRIBER")

	if len(private) == 0 || len(public) == 0 || len(subscriber) == 0 {
		return discardNotification{}
	}

	return &webpushSender{
		subscriber: subscriber,
		public:     public,
		private:    private,
	}
}

func (l discardNotification) Send(NotificationFor) error {
	return nil
}

func (s *webpushSender) Send(n NotificationFor) error {
	if len(n.Updates) == 0 {
		return nil
	}

	resp, err := webpush.SendNotification([]byte(n.Updates[0].Update.Description),
		n.Subscription,
		&webpush.Options{
			Subscriber:      s.subscriber,
			Topic:           n.Updates[0].ID(),
			VAPIDPublicKey:  s.public,
			VAPIDPrivateKey: s.private,
			TTL:             60,
		})

	defer resp.Body.Close()

	if err != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Sending push notification, response: %s, error: %w",
			string(response), err)

	}
	return nil
}
