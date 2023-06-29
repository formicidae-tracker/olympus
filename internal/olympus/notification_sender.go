package olympus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/formicidae-tracker/olympus/pkg/tm"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	log                         *logrus.Entry
	ctx                         context.Context
}

type WebPushAction struct {
	Action string `json:"action,omitempty"`
	Title  string `json:"title,omitempty"`
}

type WebPushTargetAction struct {
	Operation string `json:"operation,omitempty"`
	URL       string `json:"url,omitempty"`
}

type WebPushData struct {
	OnActionClick map[string]WebPushTargetAction `json:"onActionClick,omitempty"`
}

type WebPushNotification struct {
	Title   string          `json:"title,omitempty"` // Mandatory
	Body    string          `json:"body,omitempty"`  // Mandatory
	Actions []WebPushAction `json:"actions,omitempty"`
	Data    WebPushData     `json:"data,omitempty"`
	Icon    string          `json:"icon,omitempty"`    // Yes we would like something : https://web.dev/push-notifications-display-a-notification/#icon
	Image   string          `json:"image,omitempty"`   // Preferably not see https://web.dev/push-notifications-display-a-notification/#image
	Badge   string          `json:"badge,omitempty"`   // Yes we would like something : https://web.dev/push-notifications-display-a-notification/#badge
	Vibrate []int           `json:"vibrate,omitempty"` // Seems uneffective
	Sound   string          `json:"sound.omitempty"`   // Doesn't work
}

func NewWebPushNotification(updates []ZonedAlarmUpdate) (WebPushNotification, error) {
	switch len(updates) {
	case 0:
		return WebPushNotification{}, errors.New("no updates")
	case 1:
		return NewSingleWebPushNotification(updates[0]), nil
	default:
		return NewMultiWebPushNotification(updates), nil
	}
}

func buildURL(zone string) string {
	if zone == "services" {
		return "/logs"
	}
	splits := strings.Split(zone, ".")
	if len(splits) < 2 {
		return "/"
	}
	return fmt.Sprintf("/host/%s/zone/%s", splits[0], splits[1])
}

func getIconURL(level api.AlarmLevel) string {
	return "/assets/fort.svg"
}

func getBadgeURL(level api.AlarmLevel) string {
	return "/assets/badge.png"
}

func NewSingleWebPushNotification(update ZonedAlarmUpdate) WebPushNotification {
	return WebPushNotification{
		Title: fmt.Sprintf("One %s on %s",
			cases.Title(language.English).String(update.Update.Level.String()),
			update.Zone),
		Body: update.Update.Description,
		Data: WebPushData{
			OnActionClick: map[string]WebPushTargetAction{
				"default": {
					Operation: "navigateLastFocusedOrOpen",
					URL:       buildURL(update.Zone),
				},
			},
		},
		Badge:   getBadgeURL(update.Update.Level),
		Icon:    getIconURL(update.Update.Level),
		Vibrate: []int{10, 20, 50, 100, 50, 20, 10},
	}
}

func NewMultiWebPushNotification(updates []ZonedAlarmUpdate) WebPushNotification {
	zones, emergencies, warnings := collectInfos(updates)
	level := api.AlarmLevel_EMERGENCY
	if emergencies == 0 {
		level = api.AlarmLevel_WARNING
	}
	actions, data := buildMultiActions(zones)
	return WebPushNotification{
		Title:   buildMultiTitle(emergencies, warnings),
		Body:    buildMultiBody(zones),
		Actions: actions,
		Data:    data,
		Badge:   getBadgeURL(level),
		Icon:    getIconURL(level),
		Vibrate: []int{10, 20, 50, 100, 50, 20, 10},
	}
}

func buildMultiActions(zones []string) ([]WebPushAction, WebPushData) {
	if len(zones) == 1 {
		return nil, WebPushData{
			OnActionClick: map[string]WebPushTargetAction{
				"default": {
					Operation: "navigateLastFocusedOrOpen",
					URL:       buildURL(zones[0]),
				},
			},
		}
	}
	if len(zones) > 2 {
		return nil, WebPushData{
			OnActionClick: map[string]WebPushTargetAction{
				"default": {
					Operation: "navigateLastFocusedOrOpen",
					URL:       "/",
				},
			},
		}
	}

	actions := make([]WebPushAction, 0, 2)
	data := WebPushData{
		OnActionClick: map[string]WebPushTargetAction{
			"default": {
				Operation: "navigateLastFocusedOrOpen",
				URL:       "/",
			},
		},
	}

	for _, z := range zones {
		actions = append(actions, WebPushAction{
			Action: z,
			Title:  fmt.Sprintf("Open *%s*", z),
		})
		data.OnActionClick[z] = WebPushTargetAction{
			Operation: "navigateLastFocusedOrOpen",
			URL:       buildURL(z),
		}
	}

	return actions, data
}

func collectInfos(updates []ZonedAlarmUpdate) (zones []string, emergencies int, warnings int) {
	zonesSet := map[string]bool{}
	warnings = 0
	emergencies = 0
	for _, u := range updates {
		if u.Update.Level == api.AlarmLevel_WARNING {
			warnings++
		} else {
			emergencies++
		}
		zonesSet[u.Zone] = true
	}

	zones = make([]string, 0, len(zonesSet))
	for z := range zonesSet {
		zones = append(zones, z)
	}
	sort.Strings(zones)

	return
}

func buildMultiTitle(emergencies, warnings int) string {
	if emergencies == 0 {
		return fmt.Sprintf("%d New Warnings", warnings)
	}
	if emergencies == 1 {
		if warnings == 1 {
			return "1 New Emergency and 1 New Warning"
		}
		return fmt.Sprintf("1 New Emergency and %d New Warnings", warnings)
	}
	if warnings == 0 {
		return fmt.Sprintf("%d New Emergencies", emergencies)
	}

	if warnings == 1 {
		return fmt.Sprintf("%d New Emergencies and 1 New Warning", emergencies)
	}

	return fmt.Sprintf("%d New Emergencies and %d New Warnings",
		emergencies, warnings)
}

func buildMultiBody(zones []string) string {

	if len(zones) == 1 {
		return fmt.Sprintf("%s has new alarms", zones[0])
	}
	if len(zones) == 2 {
		return fmt.Sprintf("%s and %s have alarms", zones[0], zones[1])
	}
	if len(zones) == 3 {
		return fmt.Sprintf("%s, %s and one another have alarms", zones[0], zones[1])
	}
	return fmt.Sprintf("%s, %s and %d others have alarms", zones[0], zones[1], len(zones)-2)
}

func BuildNotificationPayload(updates []ZonedAlarmUpdate) ([]byte, error) {
	notification, err := NewWebPushNotification(updates)
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]WebPushNotification{"notification": notification})
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

	logger := logrus.New()

	if len(os.Getenv("OLYMPUS_DEBUG_WEBPUSH")) == 0 {
		logger.SetOutput(io.Discard)
	}

	res.log = tm.NewLogger("webpush")

	return res, nil
}

func (l discardNotification) Send(NotificationFor) error {
	return nil
}

func (s *webpushSender) Send(n NotificationFor) (err error) {
	defer func() {
		entry := s.log.WithFields(logrus.Fields{
			"endpoint": n.Subscription.Endpoint,
			"updates":  n.Updates,
		})

		if err != nil {
			entry.WithField("error", err).Errorf("could not send notification")
		} else {
			entry.Debugf("sent")
		}
	}()

	if len(n.Updates) == 0 {
		return nil
	}

	payload, err := BuildNotificationPayload(n.Updates)
	if err != nil {
		return err
	}

	topic := n.Updates[0].ID()
	if len(n.Updates) > 1 {
		topic = "multiple"
	}

	resp, err := webpush.SendNotification(payload,
		n.Subscription,
		&webpush.Options{
			Subscriber:      s.subscriber,
			Topic:           topic,
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
