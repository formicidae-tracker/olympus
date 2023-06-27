package olympus

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/formicidae-tracker/olympus/pkg/api"
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
	Outgoing() <-chan NotificationFor

	RegisterPushSubscription(*webpush.Subscription) error
	UpdatePushSubscription(*api.NotificationSettingsUpdate) error

	Loop()
}

type zoneRegistration struct {
	potentialEndpoints map[string]bool
}

type NotificationSubscription struct {
	Push     *webpush.Subscription
	Settings api.NotificationSettings
}

type notifier struct {
	mx                   sync.RWMutex
	incoming             chan ZonedAlarmUpdate
	outgoingNotification chan NotificationFor

	zones map[string]zoneRegistration

	subscriptions *PersistentMap[*NotificationSubscription]

	wg sync.WaitGroup

	outgoing map[string]chan<- ZonedAlarmUpdate

	batchPeriod time.Duration
	log         *log.Logger
}

func NewNotifier(batchPeriod time.Duration) Notifier {
	res := &notifier{
		zones:                make(map[string]zoneRegistration),
		incoming:             make(chan ZonedAlarmUpdate, 100),
		subscriptions:        NewPersistentMap[*NotificationSubscription]("push-notifications"),
		outgoing:             make(map[string]chan<- ZonedAlarmUpdate),
		outgoingNotification: make(chan NotificationFor, 100),
		batchPeriod:          batchPeriod,
		log:                  log.New(os.Stderr, "[notifications]: ", log.LstdFlags),
	}
	for _, sub := range res.subscriptions.Map {
		res.ensureOutgoing(sub.Push)
	}
	return res
}

func (n *notifier) Incoming() chan<- ZonedAlarmUpdate {
	return n.incoming
}

func (n *notifier) Outgoing() <-chan NotificationFor {
	return n.outgoingNotification
}

func (n *notifier) RegisterPushSubscription(s *webpush.Subscription) error {
	if len(s.Keys.Auth) == 0 || len(s.Keys.P256dh) == 0 {
		return errors.New("invalid push subscription: missing keys")
	}

	n.mx.Lock()
	defer n.mx.Unlock()

	n.subscriptions.Map[s.Endpoint] = &NotificationSubscription{
		Push: s,
	}

	n.ensureOutgoing(s)

	n.log.Printf("new push subscription to %s", s.Endpoint)

	return n.subscriptions.SaveKey(s.Endpoint)
}

func (n *notifier) UpdatePushSubscription(update *api.NotificationSettingsUpdate) error {
	n.mx.Lock()
	defer n.mx.Unlock()

	subscription, ok := n.subscriptions.Map[update.Endpoint]
	if ok == false {
		return UnknownEndpointError
	}

	subscription.Settings = update.Settings

	for zone, reg := range n.zones {
		reg.potentialEndpoints[update.Endpoint] = MayBeSubscribedTo(zone, update.Settings)
	}

	return n.subscriptions.SaveKey(update.Endpoint)
}

func (n *notifier) Loop() {
	defer func() {

		for _, ch := range n.outgoing {
			close(ch)
		}
		n.wg.Wait()

		close(n.outgoingNotification)
	}()

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

func (n *notifier) Register(zone string) zoneRegistration {
	n.mx.Lock()
	defer n.mx.Unlock()

	reg, ok := n.zones[zone]
	if ok == true {
		return reg
	}
	endpoints := make(map[string]bool)
	for endpoint, sub := range n.subscriptions.Map {
		endpoints[endpoint] = MayBeSubscribedTo(zone, sub.Settings)
	}

	reg = zoneRegistration{potentialEndpoints: endpoints}
	n.zones[zone] = reg

	return reg
}

func (n *notifier) getOrRegister(lock sync.Locker, zone string) zoneRegistration {
	reg, ok := n.zones[zone]
	if ok == true {
		return reg
	}
	defer lock.Lock()
	lock.Unlock()
	return n.Register(zone)
}

func (n *notifier) handle(update ZonedAlarmUpdate) {
	n.mx.RLock()
	defer n.mx.RUnlock()
	reg := n.getOrRegister(n.mx.RLocker(), update.Zone)

	for endpoint, maySend := range reg.potentialEndpoints {
		if maySend == true &&
			IsSubscribedTo(update.Zone,
				update.Update.Level,
				n.subscriptions.Map[endpoint].Settings) {
			n.outgoing[endpoint] <- update
		}
	}
}

func (n *notifier) ensureOutgoing(sub *webpush.Subscription) {
	_, ok := n.outgoing[sub.Endpoint]
	if ok == true {
		return
	}

	unfiltered := make(chan ZonedAlarmUpdate)
	filtered := make(chan []ZonedAlarmUpdate)

	n.wg.Add(2)
	go func() {
		defer n.wg.Done()
		BatchAlarmUpdate(n.batchPeriod)(filtered, unfiltered)
	}()

	go func(sub *webpush.Subscription) {
		defer n.wg.Done()
		for u := range filtered {
			n.outgoingNotification <- NotificationFor{Subscription: sub, Updates: u}
		}
	}(sub)

	n.outgoing[sub.Endpoint] = unfiltered
}
