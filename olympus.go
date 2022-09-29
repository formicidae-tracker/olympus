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
	log               *log.Logger
	zoneSubscriptions *sync.Map
	hostSubscriptions *sync.Map

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
		zoneSubscriptions: &sync.Map{},
		hostSubscriptions: &sync.Map{},
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
	errs := closeAll[subscription[ZoneLogger]](o.zoneSubscriptions)
	errsHost := closeAll[subscription[trackingLogger]](o.hostSubscriptions)
	if errsHost != nil {
		errs = append(errs, errsHost...)
	}
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

func (o *Olympus) ZoneIsRegistered(hostname, zone string) bool {
	_, ok := o.zoneSubscriptions.Load(ZoneIdentifier(hostname, zone))
	return ok
}

func (o *Olympus) GetZones() []ZoneReportSummary {
	nbActiveZones := len(o.climateLogger.OnServices())
	nbActiveTrackers := len(o.climateLogger.OnServices())
	res := make([]ZoneReportSummary, 0, nbActiveZones+nbActiveTrackers)

	o.zoneSubscriptions.Range(func(key, s any) bool {
		sub, ok := s.(subscription[ZoneLogger])
		if ok == false {
			return true
		}
		z := sub.object
		r := z.GetClimateReport()
		sum := ZoneReportSummary{
			Host:    z.Host(),
			Name:    z.ZoneName(),
			Climate: &r,
		}

		if s, ok := o.hostSubscriptions.Load(z.Host()); ok == true {
			sub, ok := s.(subscription[TrackingLogger])
			if ok == true {
				sum.Stream = sub.object.StreamInfo()
			}
		}

		res = append(res, sum)

		return true
	})

	o.hostSubscriptions.Range(func(key, s any) bool {
		host := key.(string)
		sub, ok := s.(subscription[TrackingLogger])
		if ok == false {
			return true
		}
		if _, ok := o.zoneSubscriptions.Load(ZoneIdentifier(host, "box")); ok == true {
			return true
		}
		res = append(res, ZoneReportSummary{
			Host:   host,
			Name:   "box",
			Stream: sub.object.StreamInfo(),
		})
		return true
	})

	sort.Slice(res, func(i, j int) bool {
		if res[i].Host == res[j].Host {
			return res[i].Name < res[j].Name
		}
		return res[i].Host < res[j].Host
	})

	return res
}

func (o *Olympus) getZoneLogger(host, zone string) (ZoneLogger, error) {
	zoneIdentifier := ZoneIdentifier(host, zone)
	s, ok := o.zoneSubscriptions.Load(zoneIdentifier)
	if ok == false {
		return nil, ZoneNotFoundError(zoneIdentifier)
	}
	z, ok := s.(subscription[ZoneLogger])
	if ok == false {
		return nil, fmt.Errorf("Internal variable conversion error")
	}

	return z.object, nil
}

func (o *Olympus) getTracking(host string) (TrackingLogger, error) {
	s, ok := o.hostSubscriptions.Load(host)
	if ok == false {
		return nil, HostNotFoundError(host)
	}
	sub, ok := s.(subscription[TrackingLogger])
	if ok == false {
		return nil, fmt.Errorf("internal variable conversion error")
	}
	return sub.object, nil
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
	z := NewZoneLogger(declaration)
	finish := make(chan struct{})
	_, loaded := o.zoneSubscriptions.LoadOrStore(z.ZoneIdentifier(),
		subscription[ZoneLogger]{object: z, finish: finish})
	if loaded == true {
		return nil, nil, fmt.Errorf("Zone '%s' is already registered", z.ZoneIdentifier())
	}
	return z, finish, nil
}

func (o *Olympus) UnregisterZone(zoneIdentifier string) error {
	s, loaded := o.zoneSubscriptions.LoadAndDelete(zoneIdentifier)
	if loaded == false {
		return ZoneNotFoundError(zoneIdentifier)
	}
	return s.(subscription[ZoneLogger]).Close()
}

func (o *Olympus) RegisterTracker(declaration *proto.TrackingDeclaration) (TrackingLogger, <-chan struct{}, error) {
	if declaration.StreamServer != o.hostname+".local" {
		return nil, nil, fmt.Errorf("unexpected server %s (expect: %s)", declaration.StreamServer, o.hostname+".local")
	}
	t := NewTrackingLogger(declaration)
	finish := make(chan struct{})
	_, loaded := o.hostSubscriptions.LoadOrStore(declaration.Hostname, subscription[TrackingLogger]{
		object: t,
		finish: finish,
	})

	if loaded == true {
		return nil, nil, fmt.Errorf("Tracking host '%s' is already registered", declaration.Hostname)
	}

	return t, finish, nil
}

func (o *Olympus) UnregisterTracker(host string) error {
	s, loaded := o.hostSubscriptions.LoadAndDelete(host)
	if loaded == false {
		return HostNotFoundError(host)
	}
	return s.(subscription[TrackingLogger]).Close()
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
