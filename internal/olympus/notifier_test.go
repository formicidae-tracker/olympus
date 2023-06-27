package olympus

import (
	"path"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/formicidae-tracker/olympus/pkg/api"
	. "gopkg.in/check.v1"
)

type NotifierSuite struct {
	datapath string

	notifier Notifier
}

var _ = Suite(&NotifierSuite{})

func (s *NotifierSuite) SetUpSuite(c *C) {
	s.datapath = _datapath
}

func (s *NotifierSuite) TearDownSuite(c *C) {
	_datapath = s.datapath
}

func (s *NotifierSuite) SetUpTest(c *C) {
	_datapath = c.MkDir()
	s.notifier = NewNotifier(0)
}

func (s *NotifierSuite) TestClosingOnIncoming(c *C) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.notifier.Loop()
	}()

	s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "a",
	})

	s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "b",
	})

	close(s.notifier.Incoming())

	grace := 5 * time.Millisecond
	select {
	case <-time.After(grace):
		c.Errorf("notifier did not close after %s", grace)
	case <-done:
	}

	select {
	case _, ok := <-s.notifier.Outgoing():
		c.Check(ok, Equals, false)
	default:
		c.Errorf("notifier did not close its outgoing channel")
	}
}

type alarmData struct {
	Endpoint string
	ID       string
	Level    api.AlarmLevel
}

func (d alarmData) ToAlarmUpdate() ZonedAlarmUpdate {
	zone, ID := path.Split(d.ID)
	zone = strings.TrimSuffix(zone, "/")
	return ZonedAlarmUpdate{
		Zone: zone,
		Update: &api.AlarmUpdate{
			Identification: ID,
			Level:          d.Level,
		},
	}
}

func (d alarmData) ToNotification() NotificationFor {
	return NotificationFor{
		Subscription: &webpush.Subscription{Endpoint: d.Endpoint},
		Updates:      []ZonedAlarmUpdate{d.ToAlarmUpdate()},
	}
}

func (s *NotifierSuite) TestSendOnlyToSubscribed(c *C) {
	go func() {
		s.notifier.Loop()
	}()
	c.Check(s.notifier.RegisterPushSubscription(&webpush.Subscription{Endpoint: "a",
		Keys: webpush.Keys{Auth: "a", P256dh: "a"}}), IsNil)
	c.Check(s.notifier.RegisterPushSubscription(&webpush.Subscription{Endpoint: "b",
		Keys: webpush.Keys{Auth: "a", P256dh: "a"}}), IsNil)
	c.Check(s.notifier.RegisterPushSubscription(&webpush.Subscription{Endpoint: "c",
		Keys: webpush.Keys{Auth: "a", P256dh: "a"}}), IsNil)
	c.Check(s.notifier.RegisterPushSubscription(&webpush.Subscription{Endpoint: "d",
		Keys: webpush.Keys{Auth: "a", P256dh: "a"}}), IsNil)

	c.Check(s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{Endpoint: "a"}), IsNil)
	c.Check(s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "b",
		Settings: api.NotificationSettings{
			SubscribeToAll: true,
		},
	}), IsNil)
	c.Check(s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "c",
		Settings: api.NotificationSettings{
			NotifyOnWarning: true,
			Subscriptions:   []string{"foo"},
		},
	}), IsNil)
	c.Check(s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "d",
		Settings: api.NotificationSettings{
			NotifyNonGraceful: true,
		},
	}), IsNil)

	alarms := []alarmData{
		{"", "foo/critical", api.AlarmLevel_EMERGENCY},
		{"", "foo/warning", api.AlarmLevel_WARNING},
		{"", "bar/critical", api.AlarmLevel_EMERGENCY},
		{"", "bar/warning", api.AlarmLevel_WARNING},
		{"", "services/nongraceful", api.AlarmLevel_EMERGENCY},
	}

	expectedList := []alarmData{
		{"b", "foo/critical", api.AlarmLevel_EMERGENCY},
		{"c", "foo/critical", api.AlarmLevel_EMERGENCY},
		{"c", "foo/warning", api.AlarmLevel_WARNING},
		{"b", "bar/critical", api.AlarmLevel_EMERGENCY},
		{"d", "services/nongraceful", api.AlarmLevel_EMERGENCY},
	}
	expected := make(map[string]bool)
	for _, e := range expectedList {
		ID := path.Join(e.Endpoint, e.ID)
		expected[ID] = true
	}

	go func() {
		for _, a := range alarms {
			s.notifier.Incoming() <- a.ToAlarmUpdate()
		}
		time.Sleep(5 * time.Millisecond)
		close(s.notifier.Incoming())
	}()

	for r := range s.notifier.Outgoing() {
		ID := path.Join(r.Subscription.Endpoint, r.Updates[0].ID())
		c.Check(expected[ID], Equals, true, Commentf("for %s", ID))
		delete(expected, ID)
	}

	for e := range expected {
		c.Errorf("Missing %s", e)
	}

}

func (s *NotifierSuite) TestSubscriptionPersistence(c *C) {
	c.Check(s.notifier.RegisterPushSubscription(&webpush.Subscription{Endpoint: "a",
		Keys: webpush.Keys{Auth: "a", P256dh: "a"}}), IsNil)
	c.Check(s.notifier.UpdatePushSubscription(&api.NotificationSettingsUpdate{
		Endpoint: "a",
		Settings: api.NotificationSettings{
			SubscribeToAll: true,
		},
	}), IsNil)

	notifier := NewNotifier(0)

	go func() {
		notifier.Loop()
	}()

	alarms := []alarmData{
		{"", "foo/critical", api.AlarmLevel_EMERGENCY},
		{"", "bar/critical", api.AlarmLevel_EMERGENCY},
		{"", "services/nongraceful", api.AlarmLevel_EMERGENCY},
	}

	expectedList := []alarmData{
		{"a", "foo/critical", api.AlarmLevel_EMERGENCY},
		{"a", "bar/critical", api.AlarmLevel_EMERGENCY},
	}

	expected := make(map[string]bool)
	for _, e := range expectedList {
		ID := path.Join(e.Endpoint, e.ID)
		expected[ID] = true
	}

	go func() {
		for _, a := range alarms {
			notifier.Incoming() <- a.ToAlarmUpdate()
		}
		time.Sleep(5 * time.Millisecond)
		close(notifier.Incoming())
	}()

	for r := range notifier.Outgoing() {
		ID := path.Join(r.Subscription.Endpoint, r.Updates[0].ID())
		c.Check(expected[ID], Equals, true, Commentf("for %s", ID))
		delete(expected, ID)
	}

	for e := range expected {
		c.Errorf("Missing %s", e)
	}

}
