package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/formicidae-tracker/olympus/api"
	"github.com/gorilla/mux"
)

var UnknownEndpointError = errors.New("unknown PushSubscription endpoint")

type UnexpectedStreamServerError struct {
	Got, Expected string
}

func (e UnexpectedStreamServerError) Error() string {
	return fmt.Sprintf("unexpected stream server %s. expected: %s", e.Got, e.Expected)
}

type ZoneNotFoundError string

func (z ZoneNotFoundError) Error() string {
	return fmt.Sprintf("olympus: unknown zone '%s'", string(z))
}

type NoClimateRunningError string

func (z NoClimateRunningError) Error() string {
	return fmt.Sprintf("olympus: no climate running in zone '%s'", string(z))
}

type NoTrackingRunningError string

func (z NoTrackingRunningError) Error() string {
	return fmt.Sprintf("olympus: no tracking running in zone '%s'", string(z))
}

type AlreadyExistError string

func (a AlreadyExistError) Error() string {
	return string(a) + " is already registered"
}

type ClosedOlympusServerError struct{}

func (e ClosedOlympusServerError) Error() string {
	return "closed olympus server"
}

type multipleError []error

func appendError(e multipleError, errors ...error) multipleError {
	for _, err := range errors {
		if err == nil {
			continue
		}
		e = append(e, err)
	}
	return e
}

func (m multipleError) Error() string {
	if len(m) == 0 {
		return ""
	}
	if len(m) == 1 {
		return m[0].Error()
	}

	res := "multiple errors:"
	for _, e := range m {
		res += "\n" + e.Error()
	}
	return res
}

type Olympus struct {
	api.UnimplementedOlympusServer

	mx             sync.RWMutex
	subscriptionWg sync.WaitGroup
	notificationWg sync.WaitGroup

	log         *log.Logger
	csrfHandler *CSRFHandler

	cancelSubscription  context.CancelFunc
	subscriptionContext context.Context
	subscriptions       map[string]*subscription

	serviceLogger ServiceLogger

	unfilteredAlarms   chan ZonedAlarmUpdate
	notifier           Notifier
	notificationSender NotificationSender
	serverPublicKey    string
	serverSecret       []byte

	hostname string
}

type GrpcSubscription[T any] struct {
	zone        string
	object      T
	alarmLogger AlarmLogger
	updates     chan<- ZonedAlarmUpdate
}

func (s GrpcSubscription[T]) NotifyAlarms(updates []*api.AlarmUpdate) {
	for _, u := range updates {
		s.updates <- ZonedAlarmUpdate{Zone: s.zone, Update: u}
	}
}

type subscription struct {
	host, name  string
	climate     *GrpcSubscription[ClimateLogger]
	tracking    *GrpcSubscription[TrackingLogger]
	alarmLogger AlarmLogger
}

func NewOlympus() (*Olympus, error) {
	ctx, cancel := context.WithCancel(context.Background())

	minimumOn := 1 * time.Minute
	batchPeriod := 5 * time.Minute

	debugWebpush := os.Getenv("OLYMPUS_DEBUG_WEBPUSH")
	if len(debugWebpush) > 0 {
		minimumOn = 1 * time.Second

		var err error
		batchPeriod, err = time.ParseDuration(debugWebpush)
		if err != nil {
			batchPeriod = 5 * time.Second
		}
	}

	res := &Olympus{
		log:                 log.New(os.Stderr, "[olympus]: ", log.LstdFlags),
		subscriptionContext: ctx,
		cancelSubscription:  cancel,
		subscriptions:       make(map[string]*subscription),
		serviceLogger:       NewServiceLogger(),
		unfilteredAlarms:    make(chan ZonedAlarmUpdate, 100),
		notifier:            NewNotifier(batchPeriod),
		serverPublicKey:     os.Getenv("OLYMPUS_VAPID_PUBLIC"),
	}
	var err error

	res.buildCSRFHandler()

	res.hostname, err = os.Hostname()
	if err != nil {
		return nil, err
	}

	res.notificationSender, err = NewNotificationSender()
	if err != nil {
		res.log.Printf("push notifications will be disabled: %s", err)
	}

	res.notificationWg.Add(3)
	go func() {
		defer res.notificationWg.Done()
		FilterAlarmUpdates(minimumOn)(res.notifier.Incoming(), res.unfilteredAlarms)
	}()

	go func() {
		defer res.notificationWg.Done()
		res.notifier.Loop()
	}()

	go func() {
		defer res.notificationWg.Done()
		for n := range res.notifier.Outgoing() {
			if err := res.notificationSender.Send(n); err != nil {
				res.log.Printf("could not send notificaton: %s", err)
			}
		}
	}()

	return res, nil
}

func getOlympusSecret() ([]byte, error) {
	secret64 := os.Getenv("OLYMPUS_SECRET")
	if len(secret64) == 0 {
		return nil, errors.New("OLYMPUS_SECRET environment variable is not set")
	}
	secret, err := base64.URLEncoding.DecodeString(secret64)
	if err != nil {
		return nil, errors.New("could not decode OLYMPUS_SECRET")
	}
	return secret, nil
}

func (o *Olympus) buildCSRFHandler() {
	secret, err := getOlympusSecret()
	if err != nil {
		o.log.Printf("could not get olympus secret: %s", err)
		return
	}
	o.csrfHandler, err = NewCSRFHandler(secret)
	if err != nil {
		o.log.Printf("could not set CSRF Handler: %s", err)
	}
}

func (o *Olympus) Close() (err error) {
	o.mx.Lock()
	defer o.mx.Unlock()

	o.cancelSubscription()

	// wait for all subscription to terminate to avoid panic on alarms
	// channels. There is almost a slight race condition where an
	// attempt to register could happen, but since the context for
	// gRPC subscription is done, it should not read and process
	// messages.
	o.mx.Unlock()
	o.subscriptionWg.Wait()
	o.mx.Lock()

	o.subscriptions = nil

	defer func() {
		rerr := recover()
		if rerr == nil {
			return
		}
		err = fmt.Errorf("%s", rerr)
	}()

	close(o.unfilteredAlarms)
	o.notificationWg.Wait()

	return nil
}

func (o *Olympus) GetServiceLogs() []api.ServiceLog {
	return o.serviceLogger.Logs()
}

func (o *Olympus) ZoneIsRegistered(host, zone string) bool {
	o.mx.RLock()
	defer o.mx.RUnlock()
	if o.subscriptions == nil {
		return false
	}
	_, ok := o.subscriptions[ZoneIdentifier(host, zone)]
	return ok
}

func (o *Olympus) GetZones() []api.ZoneReportSummary {
	o.mx.RLock()
	defer o.mx.RUnlock()
	if o.subscriptions == nil {
		return []api.ZoneReportSummary{}
	}

	res := make([]api.ZoneReportSummary, 0, len(o.subscriptions))

	for _, s := range o.subscriptions {
		sum := api.ZoneReportSummary{
			Host: s.host,
			Name: s.name,
		}
		if s.climate != nil {
			sum.Climate = s.climate.object.GetClimateReport()
		}

		if s.tracking != nil {
			sum.Tracking = s.tracking.object.TrackingInfo()
		}

		if s.alarmLogger != nil {
			sum.ActiveWarnings, sum.ActiveEmergencies = s.alarmLogger.ActiveAlarmsCount()
		}

		res = append(res, sum)
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Host == res[j].Host {
			return res[i].Name < res[j].Name
		}
		return res[i].Host < res[j].Host
	})

	return res
}

func (o *Olympus) getClimateLogger(host, zone string) (ClimateLogger, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()
	if o.subscriptions == nil {
		return nil, ClosedOlympusServerError{}
	}
	zoneIdentifier := ZoneIdentifier(host, zone)
	s, ok := o.subscriptions[zoneIdentifier]
	if ok == false {
		return nil, ZoneNotFoundError(zoneIdentifier)
	}
	if s.climate == nil {
		return nil, NoClimateRunningError(zoneIdentifier)
	}
	return s.climate.object, nil
}

func (o *Olympus) getTrackingLogger(host string) (TrackingLogger, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	if o.subscriptions == nil {
		return nil, ClosedOlympusServerError{}
	}
	zoneIdentifier := ZoneIdentifier(host, "box")

	s, ok := o.subscriptions[zoneIdentifier]
	if ok == false {
		return nil, ZoneNotFoundError(host)
	}
	if s.tracking == nil {
		return nil, NoTrackingRunningError(zoneIdentifier)
	}
	return s.tracking.object, nil
}

func (o *Olympus) getAlarmLogger(host, zone string) (AlarmLogger, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	if o.subscriptions == nil {
		return nil, ClosedOlympusServerError{}
	}
	zoneIdentifier := ZoneIdentifier(host, zone)

	s, ok := o.subscriptions[zoneIdentifier]
	if ok == false || s.alarmLogger == nil {
		return nil, ZoneNotFoundError(zoneIdentifier)
	}
	return s.alarmLogger, nil
}

// GetClimateTimeSeries returns the time series for a zone within a
// given window. window should be one of "10m","1h","1d", "1w".  It
// may return a ZoneNotFoundError.
func (o *Olympus) GetClimateTimeSerie(host, zone, window string) (api.ClimateTimeSeries, error) {
	z, err := o.getClimateLogger(host, zone)
	if err != nil {
		return api.ClimateTimeSeries{}, err
	}
	return z.GetClimateTimeSeries(window), nil
}

func (o *Olympus) GetZoneReport(host, zone string) (*api.ZoneReport, error) {
	z, errZone := o.getClimateLogger(host, zone)
	i, errTracking := o.getTrackingLogger(host)
	a, errAlarm := o.getAlarmLogger(host, zone)
	if errZone != nil && errTracking != nil && errAlarm != nil {
		return nil, errZone
	}
	res := &api.ZoneReport{
		Host: host,
		Name: zone,
	}
	if errZone == nil {
		res.Climate = z.GetClimateReport()
	}
	if errTracking == nil {
		res.Tracking = i.TrackingInfo()
	}
	if errAlarm == nil {
		res.Alarms = a.GetReports()
	}

	return res, nil
}

func (o *Olympus) GetAlarmReports(host, zone string) ([]api.AlarmReport, error) {
	a, err := o.getAlarmLogger(host, zone)
	if err != nil {
		return nil, err
	}
	return a.GetReports(), nil
}

func (o *Olympus) RegisterClimate(declaration *api.ClimateDeclaration) (*GrpcSubscription[ClimateLogger], error) {
	o.mx.Lock()
	defer o.mx.Unlock()
	if o.subscriptions == nil {
		return nil, ClosedOlympusServerError{}
	}

	zoneIdentifier := ZoneIdentifier(declaration.Host, declaration.Name)
	sub, ok := o.subscriptions[zoneIdentifier]
	if ok == true && sub.climate != nil {
		return nil, AlreadyExistError("zone '" + zoneIdentifier + "'")
	}

	if ok == false {
		alarmLogger := NewAlarmLogger()
		sub = &subscription{
			host:        declaration.Host,
			name:        declaration.Name,
			alarmLogger: alarmLogger,
		}
		o.subscriptions[zoneIdentifier] = sub
	}
	o.subscriptionWg.Add(1)

	sub.climate = &GrpcSubscription[ClimateLogger]{
		zone:        zoneIdentifier,
		object:      NewClimateLogger(declaration),
		alarmLogger: sub.alarmLogger,
		updates:     o.unfilteredAlarms,
	}

	go o.serviceLogger.Log(zoneIdentifier+".climate", true, true)

	return sub.climate, nil
}

func (o *Olympus) UnregisterClimate(host, name string, graceful bool) error {
	zoneIdentifier := ZoneIdentifier(host, name)
	o.mx.Lock()
	defer o.mx.Unlock()

	if o.subscriptions == nil {
		return ClosedOlympusServerError{}
	}

	s, ok := o.subscriptions[zoneIdentifier]
	if ok == false || s.climate == nil {
		return ZoneNotFoundError(zoneIdentifier)
	}

	defer func() {
		s.climate = nil
		o.subscriptionWg.Done()
	}()

	if s.tracking == nil {
		delete(o.subscriptions, zoneIdentifier)
	}

	o.serviceLogger.Log(zoneIdentifier+".climate", false, graceful)

	return nil
}

func (o *Olympus) RegisterTracking(declaration *api.TrackingDeclaration) (*GrpcSubscription[TrackingLogger], error) {
	if declaration.StreamServer != o.hostname+".local" {
		return nil, UnexpectedStreamServerError{Got: declaration.StreamServer,
			Expected: o.hostname + ".local"}
	}

	o.mx.Lock()
	defer o.mx.Unlock()

	if o.subscriptions == nil {
		return nil, ClosedOlympusServerError{}
	}

	zoneIdentifier := ZoneIdentifier(declaration.Hostname, "box")

	sub, ok := o.subscriptions[zoneIdentifier]
	if ok == true && sub.tracking != nil {
		return nil, AlreadyExistError("tracking '" + declaration.Hostname + "'")
	}

	if ok == false {
		alarmLogger := NewAlarmLogger()
		sub = &subscription{
			host:        declaration.Hostname,
			name:        "box",
			alarmLogger: alarmLogger,
		}
		o.subscriptions[zoneIdentifier] = sub
	}
	o.subscriptionWg.Add(1)

	sub.tracking = &GrpcSubscription[TrackingLogger]{
		zone:        zoneIdentifier,
		object:      NewTrackingLogger(declaration),
		alarmLogger: sub.alarmLogger,
		updates:     o.unfilteredAlarms,
	}

	o.serviceLogger.Log(zoneIdentifier+".tracking", true, true)

	return sub.tracking, nil
}

func (o *Olympus) UnregisterTracker(host string, graceful bool) error {
	o.mx.Lock()
	defer o.mx.Unlock()

	if o.subscriptions == nil {
		return ClosedOlympusServerError{}
	}
	zoneIdentifier := ZoneIdentifier(host, "box")
	s, ok := o.subscriptions[zoneIdentifier]
	if ok == false || s.tracking == nil {
		return ZoneNotFoundError(zoneIdentifier)
	}
	defer func() {
		s.climate = nil
		o.subscriptionWg.Done()
	}()

	if s.climate == nil {
		delete(o.subscriptions, zoneIdentifier)
	}

	o.serviceLogger.Log(zoneIdentifier+".tracking", false, graceful)
	return nil
}

func (o *Olympus) NotifyAlarm(zone string, update *api.AlarmUpdate) {
	o.mx.Lock()
	defer o.mx.Unlock()

	o.notifier.Incoming() <- ZonedAlarmUpdate{Zone: zone, Update: update}
}

func (o *Olympus) setRoutes(router *mux.Router) {
	o.setFetchRoutes(router)
	if o.csrfHandler != nil {
		o.setNotificationRoutes(router)
	} else {
		o.log.Printf("No CSRF handler set, notifications routes are disabled")
	}
}

func (o *Olympus) setFetchRoutes(router *mux.Router) {
	router.HandleFunc("/api/zones", func(w http.ResponseWriter, r *http.Request) {
		res := o.GetZones()
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		res := o.GetServiceLogs()
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}/climate", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetClimateTimeSerie(vars["hname"], vars["zname"], r.URL.Query().Get("window"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}/alarms", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetAlarmReports(vars["hname"], vars["zname"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetZoneReport(vars["hname"], vars["zname"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		type Version struct {
			Version string `json:"version"`
		}
		version := Version{Version: OLYMPUS_VERSION}
		JSONify(w, &version)
	}).Methods("GET")
}

func (o *Olympus) setNotificationRoutes(router *mux.Router) {
	router.Handle("/api/notifications/key",
		o.csrfHandler.SetCSRFCookie(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if len(o.serverPublicKey) == 0 {
					http.Error(w, "no key available", http.StatusNotFound)
					return
				}

				w.Header().Add("Cache-Control", "no-store")

				type publicKey struct {
					Key string `json:"serverPublicKey"`
				}

				JSONify(w, &publicKey{o.serverPublicKey})
			}))).Methods("GET")

	subrouter := router.PathPrefix("/api/notifications").Subrouter()

	subrouter.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")

		update, err := Golangify[api.NotificationSettingsUpdate](r)
		if err != nil {
			o.log.Printf("invalid notification settings update: %s", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := o.notifier.UpdatePushSubscription(update); err != nil {
			o.log.Printf("could not update notification settings: %s", err)
			if err == UnknownEndpointError {
				http.Error(w, "unknown subscription endpoint", http.StatusNotFound)
			} else {
				http.Error(w, "notification settings error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	subrouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")

		sub, err := Golangify[webpush.Subscription](r)
		if err != nil {
			o.log.Printf("invalid push subscription: %s", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := o.notifier.RegisterPushSubscription(sub); err != nil {
			o.log.Printf("could not register push subscription: %s", err)
			http.Error(w, "push subscription error", http.StatusInternalServerError)
			return
		}

	}).Methods("POST")

	subrouter.Use(o.csrfHandler.CheckCSRFCookie)
}
