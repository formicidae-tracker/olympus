package main

import (
	"sync"

	"github.com/formicidae-tracker/olympus/api"
)

func MayBeSubscribedTo(zone string, settings api.NotificationSettings) bool {
	return IsSubscribedTo(zone, api.AlarmLevel_EMERGENCY, settings)
}

func IsSubscribedTo(zone string, level api.AlarmLevel, settings api.NotificationSettings) bool {
	if zone == "services" {
		return settings.NotifyNonGraceful
	}

	if level == api.AlarmLevel_WARNING && settings.NotifyOnWarning == false {
		return false
	}

	return settings.SubscribedTo(zone)
}

type Notifier interface {
	Incoming() chan<- ZonedAlarmUpdate

	UpdatePushSubscription(*api.NotificationSettingsUpdate) error

	Loop()
}

type zoneRegistration struct {
	potentialEndpoints map[string]bool
}

type notifier struct {
	mx       sync.RWMutex
	incoming chan ZonedAlarmUpdate

	zones map[string]zoneRegistration

	subscriptions *PersistentMap[api.NotificationSettings]
}

func NewNotifier() Notifier {
	return &notifier{
		zones:         make(map[string]zoneRegistration),
		incoming:      make(chan ZonedAlarmUpdate, 100),
		subscriptions: NewPersistentMap[api.NotificationSettings]("push-notifications"),
	}
}

func (n *notifier) Incoming() chan<- ZonedAlarmUpdate {
	return n.incoming
}

func (n *notifier) UpdatePushSubscription(update *api.NotificationSettingsUpdate) error {
	n.mx.Lock()
	defer n.mx.Unlock()

	n.subscriptions.Map[update.Endpoint] = update.Settings
	n.subscriptions.SaveKey(update.Endpoint)

	for zone, reg := range n.zones {
		reg.potentialEndpoints[update.Endpoint] = MayBeSubscribedTo(zone, update.Settings)
	}

	return n.subscriptions.SaveKey(update.Endpoint)
}

func (n *notifier) Loop() {
	for {
		select {
		case update, ok := <-n.incoming:
			if ok == false {
				return
			}
			go n.handle(update)
		}
	}
}

func (n *notifier) getOrBuildRegistration(zone string) zoneRegistration {
	n.mx.Lock()
	defer n.mx.Unlock()

	reg, ok := n.zones[zone]
	if ok == false {
		reg = n.register(zone)
	}
	return reg
}

func (n *notifier) register(zone string) zoneRegistration {
	reg, ok := n.zones[zone]
	if ok == true {
		return reg
	}
	endpoints := make(map[string]bool)
	for endpoint, settings := range n.subscriptions.Map {
		endpoints[endpoint] = MayBeSubscribedTo(zone, settings)
	}

	reg = zoneRegistration{potentialEndpoints: endpoints}
	n.zones[zone] = reg

	return reg
}

func (n *notifier) handle(update ZonedAlarmUpdate) {
	reg := n.getOrBuildRegistration(update.Zone)
	n.mx.RLock()
	defer n.mx.RUnlock()

	for endpoint, maySend := range reg.potentialEndpoints {
		if maySend == true &&
			IsSubscribedTo(update.Zone,
				update.Update.Level,
				n.subscriptions.Map[endpoint]) {

			n.sendNotification(endpoint, update)
		}
	}
}

func (n *notifier) sendNotification(endpoint string, update ZonedAlarmUpdate) {
	//TODO: send the notification
}
