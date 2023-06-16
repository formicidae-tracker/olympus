package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	log                         *log.Logger
	ctx                         context.Context
}

func NewNotificationSender() (NotificationSender, error) {
	private := os.Getenv("OLYMPUS_VAPID_PRIVATE")
	public := os.Getenv("OLYMPUS_VAPID_PUBLIC")
	subscriber := os.Getenv("OLYMPUS_PUSH_SUBSCRIBER")

	if len(private) == 0 || len(public) == 0 || len(subscriber) == 0 {

		return discardNotification{}, errors.New("missing OLYMPUS_VAPID_PUBLIC,OLYMPUS_VAPID_PRIVATE or OLYMPUS_PUSH_SUBSCRIBER")
	}

	res := &webpushSender{
		subscriber: subscriber,
		public:     public,
		private:    private,
	}

	if len(os.Getenv("OLYMPUS_DEBUG_WEBPUSH")) > 0 {
		res.log = log.New(os.Stderr, "[webpush]: ", log.LstdFlags)
	} else {
		res.log = log.New(io.Discard, "", 0)
	}

	return res, nil
}

func (l discardNotification) Send(NotificationFor) error {
	return nil
}

func (s *webpushSender) Send(n NotificationFor) error {
	if len(n.Updates) == 0 {
		return nil
	}

	s.log.Printf("sending to %s: %s", n.Subscription.Endpoint, n.Updates)

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

	response, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("Sending push notification, response: %s, error: %w",
			string(response), err)
	}
	return nil
}
