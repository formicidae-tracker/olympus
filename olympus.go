package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/formicidae-tracker/olympus/proto"
	"github.com/gorilla/mux"
)

type ZoneNotFoundError string

func (z ZoneNotFoundError) Error() string {
	return fmt.Sprintf("olympus: unknown zone '%s'", string(z))
}

type HostNotFoundError string

func (h HostNotFoundError) Error() string {
	return fmt.Sprintf("olympus: unknown host '%s'", string(h))
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
	proto.UnimplementedOlympusServer
	mx                sync.RWMutex
	log               *log.Logger
	zoneSubscriptions map[string]subscription[ZoneLogger]
	hostSubscriptions map[string]subscription[TrackingLogger]

	trackingLogger ServiceLogger
	climateLogger  ServiceLogger

	hostname string
	slackURL string
}

type subscription[T any] struct {
	object T
	finish chan struct{}
}

func (s subscription[T]) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("already closed")
		}
	}()
	close(s.finish)
	return nil
}

func NewOlympus(slackURL string) (*Olympus, error) {
	res := &Olympus{
		log:               log.New(os.Stderr, "[olympus] :", log.LstdFlags),
		zoneSubscriptions: make(map[string]subscription[ZoneLogger]),
		hostSubscriptions: make(map[string]subscription[TrackingLogger]),
		trackingLogger:    NewServiceLogger(),
		climateLogger:     NewServiceLogger(),
		slackURL:          slackURL,
	}
	var err error
	res.hostname, err = os.Hostname()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func closeAll[T io.Closer](m *sync.Map) []error {
	var errs []error = nil
	m.Range(func(key, s any) bool {
		o, ok := s.(T)
		if ok == false {
			return true
		}
		err := o.Close()
		if err != nil {
			errs = append(errs, err)
		}
		return true
	})
	return errs
}

func (o *Olympus) Close() error {
	o.mx.Lock()
	defer o.mx.Unlock()

	var errs []error
	for _, s := range o.zoneSubscriptions {
		err := s.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	o.zoneSubscriptions = nil

	for _, s := range o.hostSubscriptions {
		err := s.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	o.hostSubscriptions = nil

	if errs == nil {
		return nil
	}
	return multipleError(errs)
}

func (o *Olympus) GetServiceLogs() ServiceLogs {
	return ServiceLogs{
		Climates: o.climateLogger.Logs(),
		Tracking: o.trackingLogger.Logs(),
	}
}

func (o *Olympus) ZoneIsRegistered(host, zone string) bool {
	o.mx.RLock()
	defer o.mx.RUnlock()
	_, ok := o.zoneSubscriptions[ZoneIdentifier(host, zone)]
	return ok
}

func (o *Olympus) GetZones() []ZoneReportSummary {
	o.mx.RLock()
	defer o.mx.RUnlock()

	nbActiveZones := len(o.climateLogger.OnServices())
	nbActiveTrackers := len(o.climateLogger.OnServices())
	res := make([]ZoneReportSummary, 0, nbActiveZones+nbActiveTrackers)

	for _, s := range o.zoneSubscriptions {
		z := s.object
		r := z.GetClimateReport()
		sum := ZoneReportSummary{
			Host:    z.Host(),
			Name:    z.ZoneName(),
			Climate: &r,
		}

		if t, ok := o.hostSubscriptions[z.Host()]; ok == true {
			if ok == true {
				sum.Stream = t.object.StreamInfo()
			}
		}

		res = append(res, sum)

	}

	for host, s := range o.hostSubscriptions {
		t := s.object
		if _, ok := o.zoneSubscriptions[ZoneIdentifier(host, "box")]; ok == true {
			continue
		}
		res = append(res, ZoneReportSummary{
			Host:   host,
			Name:   "box",
			Stream: t.StreamInfo(),
		})
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Host == res[j].Host {
			return res[i].Name < res[j].Name
		}
		return res[i].Host < res[j].Host
	})

	return res
}

func (o *Olympus) getZoneLogger(host, zone string) (ZoneLogger, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()
	zoneIdentifier := ZoneIdentifier(host, zone)
	s, ok := o.zoneSubscriptions[zoneIdentifier]
	if ok == false {
		return nil, ZoneNotFoundError(zoneIdentifier)
	}
	return s.object, nil
}

func (o *Olympus) getTracking(host string) (TrackingLogger, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	s, ok := o.hostSubscriptions[host]
	if ok == false {
		return nil, HostNotFoundError(host)
	}
	return s.object, nil
}

// GetClimateTimeSeries returns the time series for a zone within a
// given window. window should be one of "10m","1h","1d", "1w".  It
// may return a ZoneNotFoundError.
func (o *Olympus) GetClimateTimeSerie(host, zone, window string) (ClimateTimeSeries, error) {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return ClimateTimeSeries{}, err
	}
	return z.GetClimateTimeSeries(window), nil
}

func (o *Olympus) GetZoneReport(host, zone string) (ZoneReport, error) {
	z, errZone := o.getZoneLogger(host, zone)
	i, errTracking := o.getTracking(host)
	if errZone != nil && errTracking != nil {
		return ZoneReport{}, errZone
	}
	res := ZoneReport{
		Host: host,
		Name: zone,
	}
	if errZone == nil {
		r := z.GetClimateReport()
		res.Climate = &r
		res.Alarms = z.GetAlarmReports()
	}
	if errTracking == nil {
		res.Stream = i.StreamInfo()
	}
	return res, nil
}

func (o *Olympus) GetAlarmReports(host, zone string) ([]AlarmReport, error) {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return nil, err
	}
	return z.GetAlarmReports(), nil
}

func (o *Olympus) RegisterZone(declaration *proto.ZoneDeclaration) (ZoneLogger, <-chan struct{}, error) {
	o.mx.Lock()
	defer o.mx.Unlock()
	zID := ZoneIdentifier(declaration.Host, declaration.Name)
	if _, ok := o.zoneSubscriptions[zID]; ok == true {
		return nil, nil, fmt.Errorf("Zone '%s' is already registered", zID)
	}
	z := NewZoneLogger(declaration)
	finish := make(chan struct{})
	o.zoneSubscriptions[zID] = subscription[ZoneLogger]{
		object: z,
		finish: finish,
	}
	go o.climateLogger.Log(z.ZoneIdentifier(), true, true)
	return z, finish, nil
}

func (o *Olympus) UnregisterZone(zoneIdentifier string, graceful bool) error {
	o.mx.Lock()
	defer o.mx.Unlock()

	s, ok := o.zoneSubscriptions[zoneIdentifier]
	if ok == false {
		return ZoneNotFoundError(zoneIdentifier)
	}
	delete(o.zoneSubscriptions, zoneIdentifier)
	o.climateLogger.Log(zoneIdentifier, false, graceful)
	if graceful == false {
		go o.postToSlack(":warning: Climate on *%s* ended unexpectedly.", zoneIdentifier)
	}

	return s.Close()
}

func (o *Olympus) RegisterTracker(declaration *proto.TrackingDeclaration) (TrackingLogger, <-chan struct{}, error) {
	o.mx.Lock()
	defer o.mx.Unlock()

	if declaration.StreamServer != o.hostname+".local" {
		return nil, nil, fmt.Errorf("unexpected server %s (expect: %s)", declaration.StreamServer, o.hostname+".local")
	}

	if _, ok := o.hostSubscriptions[declaration.Hostname]; ok == true {
		return nil, nil, fmt.Errorf("Tracking host '%s' is already registered", declaration.Hostname)
	}

	t := NewTrackingLogger(declaration)
	finish := make(chan struct{})
	o.hostSubscriptions[declaration.Hostname] = subscription[TrackingLogger]{
		object: t,
		finish: finish,
	}

	o.trackingLogger.Log(declaration.Hostname, true, true)

	return t, finish, nil
}

func (o *Olympus) UnregisterTracker(host string, graceful bool) error {
	o.mx.Lock()
	defer o.mx.Unlock()

	s, ok := o.hostSubscriptions[host]
	if ok == false {
		return HostNotFoundError(host)
	}
	delete(o.hostSubscriptions, host)
	o.trackingLogger.Log(host, false, graceful)
	if graceful == false {
		o.postToSlack(":warning: Tracking experiment `%s` on %s ended unexpectedly.", "", host)
	}
	return s.Close()
}

func (o *Olympus) route(router *mux.Router) {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}/alarms", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetAlarmReports(vars["hname"], vars["zname"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetZoneReport(vars["hname"], vars["zname"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

}

func (o *Olympus) encodeSlackMessage(message string) (*bytes.Buffer, error) {
	type SlackMessage struct {
		Text string `json:"text"`
	}

	buffer := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buffer)
	if err := enc.Encode(SlackMessage{Text: message}); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (o *Olympus) postToSlackError(message string) error {
	buffer, err := o.encodeSlackMessage(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(o.slackURL, "application/json", buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("response %d: %s", resp.StatusCode, string(d))
	}
	return nil
}

func (o *Olympus) postToSlack(message string, args ...interface{}) {
	if len(o.slackURL) == 0 {
		return
	}
	if err := o.postToSlackError(fmt.Sprintf(message, args...)); err != nil {
		o.log.Printf("slack error: %s", err)
	}
}
