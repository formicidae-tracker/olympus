package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/proto"
	"github.com/formicidae-tracker/zeus"
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
	log           *log.Logger
	subscriptions *sync.Map

	trackingLogger ServiceLogger
	climateLogger  ServiceLogger

	hostname string
	slackURL string
}

type subscription[T any] struct {
	object T
	finish chan bool
}

func (s subscription[T]) Close() error {
	select {
	case s.finish <- true:
	default:
		//default in case connection already finished
	}
}

func NewOlympus(slackURL string) (*Olympus, error) {
	res := &Olympus{
		log:            log.New(os.Stderr, "[olympus] :", log.LstdFlags),
		subscriptions:  &sync.Map{},
		trackingLogger: NewServiceLogger(),
		climateLogger:  NewServiceLogger(),
		slackURL:       slackURL,
	}
	var err error
	res.hostname, err = os.Hostname()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (o *Olympus) Close() error {
	var errs multipleError = nil

	o.subscriptions.Range(func(key, s any) bool {
		sub, ok := s.(subscription[ZoneLogger])
		if ok == true {
			err := sub.Close()
			if err != nil {
				errs = appendError(errs, err)
			}
		}
		sub2, ok := s.(subscription[StreamInfo])
		if ok == true {
			err := sub2.Close()
			if err != nil {
				errs = appendError(errs, err)
			}
		}
		return true
	})
	if errs == nil {
		return nil
	}
	return errs
}

func (o *Olympus) GetServiceLogs() ServiceLogs {
	return ServiceLogs{
		Climates: o.climateLogger.Logs(),
		Tracking: o.trackingLogger.Logs(),
	}
}

func (o *Olympus) ZoneIsRegistered(hostname, zone string) bool {
	_, ok := o.subscriptions.Load(ZoneIdentifier(hostname, zone))
	return ok
}

func (o *Olympus) GetZones() []ZoneReportSummary {
	nbActiveZones := len(o.climateLogger.OnServices())
	nbActiveTrackers := len(o.climateLogger.OnServices())
	res := make([]ZoneReportSummary, 0, nbActiveZones+nbActiveTrackers)

	appendZone := func(z ZoneLogger) {
		r := z.GetClimateReport()
		sum := ZoneReportSummary{
			Host:    z.Host(),
			Name:    z.ZoneName(),
			Climate: &r,
		}

		if s, ok := o.subscriptions.Load(z.Host()); ok == true {
			sub, ok := s.(subscription[StreamInfo])
			if ok == true {
				sum.Stream = deepcopy.MustAnything(&sub.object).(*StreamInfo)
			}
		}

		res = append(res, sum)
	}

	appendTracking := func(host string, i StreamInfo) {
		if _, ok := o.subscriptions.Load(ZoneIdentifier(host, "box")); ok == true {
			return
		}
		res = append(res, ZoneReportSummary{
			Host:   host,
			Name:   "box",
			Stream: deepcopy.MustAnything(&i).(*StreamInfo),
		})
	}

	o.subscriptions.Range(func(key, s any) bool {
		sub, ok := s.(subscription[ZoneLogger])
		if ok == true {
			appendZone(sub.object)
		}
		sub2, ok := s.(subscription[StreamInfo])
		if ok == true {
			appendTracking(key.(string), sub2.object)
		}
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
	s, ok := o.subscriptions.Load(zoneIdentifier)
	if ok == false {
		return nil, ZoneNotFoundError(zoneIdentifier)
	}
	z, ok := s.(subscription[ZoneLogger])
	if ok == false {
		return nil, fmt.Errorf("Internal variable conversion error")
	}

	return z, nil
}

func (o *Olympus) getTracking(host string) (StreamInfo, error) {
	s, ok := o.subscriptions.Load(host)
	if ok == false {
		return StreamInfo{}, HostNotFoundError(host)
	}
	i, ok := s.(subscription[StreamInfo])
	if ok == false {
		return StreamInfo{}, fmt.Errorf("internal variable conversion error")
	}
	return i, nil
}

func (o *Olympus) GetClimateTimeSerie(host, zone, window string) (ClimateTimeSerie, error) {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return ClimateTimeSerie{}, err
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
		res.Stream = deepcopy.MustAnything(&i).(*StreamInfo)
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

func (o *Olympus) ReportClimates(host, zone string, reports []proto.ClimateReport) error {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return err
	}
	z.PushReports(reports)
	return nil
}

func (o *Olympus) ReportAlarms(host, zone string, events []proto.AlarmEvent) error {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return err
	}
	z.PushAlarms(events)
	return nil
}

func (o *Olympus) ReportTarget(host, zone string, target proto.ClimateTarget) error {
	z, err := o.getZoneLogger(host, zone)
	if err != nil {
		return err
	}
	z.PushTarget(target)
	return nil
}

func (o *Olympus) RegisterZone(declaration proto.ZoneDeclaration) error {
	zoneIdentifier := ZoneIdentifier(declaration.Host, declaration.Name)
	_, loaded := o.subscriptions.LoadOrStore(zoneIdentifier, NewZoneLogger(declaration))
	if loaded == true {
		return fmt.Errorf("%s is already registered", zoneIdentifier)
	}
	o.log.Printf("registered new zone %s", zoneIdentifier)
	return nil
}

func (o *Olympus) UnregisterZone(host, zone string, graceful bool) error {
	zoneIdentifier := ZoneIdentifier(host, zone)
	_, loaded := o.subscriptions.LoadAndDelete(zoneIdentifier)
	if loaded == false {
		return ZoneNotFoundError(zoneIdentifier)
	}
	o.climateLogger.Log(zoneIdentifier, false, graceful)
	if graceful == false {
		go o.postToSlack(":warning: Climate on *%s.%s* ended unexpectedly.", host, zone)
	} else {
		// go o.postToSlack(":information_source: Climate on *%s.%s* ended normally.", hostName, zoneName)
	}

	return nil
}

func (o *Olympus) RegisterTracker(host, streamServer, experimentName string) error {
	if streamServer != o.hostname+".local" {
		return fmt.Errorf("unexpected server %s (expect: %s)", streamServer, o.hostname+".local")
	}
	infos := &StreamInfo{
		StreamURL:    path.Join("/olympus/hls/", host+".m3u8"),
		ThumbnailURL: path.Join("/olympus", host+".png"),
	}
	_, loaded := o.subscriptions.LoadOrStore(host, infos)
	if loaded == true {
		return fmt.Errorf("Tracking on '%s' is already registered", host)
	}
	o.log.Printf("registered tracker %s", host)

	// go o.postToSlack(":information_source: Tracking experiment `%s` on *%s* started.", args.Experiment, hostname)
	return nil
}

func (o *Olympus) UnregisterTracker(host string, graceful bool) error {
	_, loaded := o.subscriptions.LoadAndDelete(host)
	if loaded == false {
		return HostNotFoundError(host)
	}
	o.trackingLogger.Log(host, false, graceful)
	o.log.Printf("unregistered tracking on '%s'", host)
	if graceful == false {
		go o.postToSlack(":warning: Tracking experiment `%s` on *%s* ended unexpectedly.", watcher.Experiment(), hostname)
	} else {
		//go o.postToSlack(":information_source: Tracking experiment `%s` on *%s* ended normally.", watcher.Experiment(), hostname)
	}

	return nil

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
		res, err := o.GetAlarmReports(zeus.ZoneIdentifier(vars["hname"], vars["zname"]))
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

func (o *Olympus) Zone(stream proto.Olympus_ZoneServer) error {
	zoneIdentifier := "unknown"
	finish := make(chan bool)
	ctx := stream.Context()

	for {
		select {
		case <- finish:
			o.logger.Printf("closing stream for %s",zoneIdentifier)
			return nil
		case <- ctx.Done():
			log.Printf("zone %s has disconnected")
		default:



}
