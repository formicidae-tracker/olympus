package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/formicidae-tracker/zeus"
	"github.com/gorilla/mux"
)

type ZoneNotFoundError struct {
	Zone string
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

func (z ZoneNotFoundError) Error() string {
	return fmt.Sprintf("olympus: unknown zone '%s'", z.Zone)
}

type Olympus struct {
	mxClimate, mxTracking sync.RWMutex
	zones                 map[string]ZoneLogger
	log                   *log.Logger

	watchers       map[string]TrackingWatcher
	trackingLogger ServiceLogger
	climateLogger  ServiceLogger
}

func NewOlympus(streamServerBaseURL string) *Olympus {
	res := &Olympus{
		zones:          make(map[string]ZoneLogger),
		log:            log.New(os.Stderr, "[olympus] :", log.LstdFlags),
		watchers:       make(map[string]TrackingWatcher),
		trackingLogger: NewServiceLogger(),
		climateLogger:  NewServiceLogger(),
	}
	go res.fetchOnlineTracking(streamServerBaseURL)
	return res
}

func (o *Olympus) Close() error {
	o.mxClimate.Lock()
	o.mxTracking.Lock()
	defer o.mxTracking.Unlock()
	defer o.mxClimate.Unlock()

	var err multipleError = nil

	for _, logger := range o.zones {
		err = appendError(err, logger.Close())
	}
	o.zones = nil
	for _, watcher := range o.watchers {
		err = appendError(err, watcher.Close())
	}
	o.watchers = nil
	if len(err) == 0 {
		return nil
	}
	return err
}

func (o *Olympus) GetZones() []ZoneReportSummary {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()
	o.mxTracking.RLock()
	defer o.mxTracking.RUnlock()
	return o.getZones()
}

func (o *Olympus) GetClimateReport(zoneIdentifier, window string) (ClimateReportTimeSerie, error) {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.getClimateReport(zoneIdentifier, window)
}

func (o *Olympus) GetZone(zoneIdentifier string) (ZoneClimateReport, error) {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.getZoneReport(zoneIdentifier)
}

func (o *Olympus) GetAlarmEventLog(zoneIdentifier string) ([]AlarmEvent, error) {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.getAlarmEventLog(zoneIdentifier)
}

func (o *Olympus) RegisterZone(zr zeus.ZoneRegistration) error {
	o.mxClimate.Lock()
	defer o.mxClimate.Unlock()

	return o.registerZone(zr)
}

func (o *Olympus) UnregisterZone(zr zeus.ZoneUnregistration) error {
	o.mxClimate.Lock()
	defer o.mxClimate.Unlock()

	return o.unregisterZone(zr.Host, zr.Name, true)
}

func (o *Olympus) ZoneIsRegistered(hostname, zoneName string) bool {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	_, res := o.zones[zeus.ZoneIdentifier(hostname, zoneName)]
	return res
}

func (o *Olympus) ReportClimate(cr zeus.NamedClimateReport) error {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.reportClimate(cr)
}

func (o *Olympus) ReportAlarm(ae zeus.AlarmEvent) error {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.reportAlarm(ae)
}

func (o *Olympus) ReportState(sr zeus.StateReport) error {
	o.mxClimate.RLock()
	defer o.mxClimate.RUnlock()

	return o.reportState(sr)
}

func (o *Olympus) RegisterTracker(args LetoTrackingRegister) error {
	o.mxTracking.Lock()
	defer o.mxTracking.Unlock()

	return o.registerTracker(args.Host, args.URL)
}

func (o *Olympus) UnregisterTracker(hostname string) error {
	o.mxTracking.Lock()
	defer o.mxTracking.Unlock()
	return o.unregisterTracker(hostname, true)
}

func (o *Olympus) getZones() []ZoneReportSummary {
	res := make([]ZoneReportSummary, 0, len(o.zones)+len(o.watchers))
	for _, z := range o.zones {
		status := z.GetReport().ZoneClimateStatus
		sum := ZoneReportSummary{
			Host:    z.Host(),
			Name:    z.ZoneName(),
			Climate: &status,
		}
		if w, ok := o.watchers[z.Host()]; ok == true {
			sum.StreamURL = w.StreamURL()
		}

		res = append(res, sum)
	}
	for host, watcher := range o.watchers {
		if _, ok := o.zones[zeus.ZoneIdentifier(host, "box")]; ok == true {
			continue
		}
		res = append(res, ZoneReportSummary{
			Host:      host,
			Name:      "box",
			StreamURL: watcher.StreamURL(),
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

func (o *Olympus) getClimateReport(zoneIdentifier, window string) (ClimateReportTimeSerie, error) {
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return ClimateReportTimeSerie{}, ZoneNotFoundError{zoneIdentifier}
	}
	return z.GetClimateReportSeries(window), nil
}

func (o *Olympus) getZoneReport(zoneIdentifier string) (ZoneClimateReport, error) {
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return ZoneClimateReport{}, ZoneNotFoundError{zoneIdentifier}
	}
	return z.GetReport(), nil
}

func (o *Olympus) getAlarmEventLog(zoneIdentifier string) ([]AlarmEvent, error) {
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return nil, ZoneNotFoundError{zoneIdentifier}
	}
	return z.GetAlarmsEventLog(), nil
}

func (o *Olympus) reportClimate(cr zeus.NamedClimateReport) error {
	z, ok := o.zones[cr.ZoneIdentifier]
	if ok == false {
		return ZoneNotFoundError{cr.ZoneIdentifier}
	}
	z.ReportChannel() <- cr.ClimateReport
	return nil
}

func (o *Olympus) reportAlarm(ae zeus.AlarmEvent) error {
	z, ok := o.zones[ae.ZoneIdentifier]
	if ok == false {
		return ZoneNotFoundError{ae.ZoneIdentifier}
	}
	z.AlarmChannel() <- ae
	return nil
}

func (o *Olympus) reportState(sr zeus.StateReport) error {
	z, ok := o.zones[sr.ZoneIdentifier]
	if ok == false {
		return ZoneNotFoundError{sr.ZoneIdentifier}
	}
	z.StateChannel() <- sr
	return nil
}

func (o *Olympus) unregisterZone(hostName, zoneName string, graceful bool) error {
	zoneIdentifier := zeus.ZoneIdentifier(hostName, zoneName)
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return ZoneNotFoundError{zoneIdentifier}
	}
	o.climateLogger.Log(hostName+"."+zoneName, false, graceful)
	o.log.Printf("unregistering %s", zoneIdentifier)
	err := z.Close()
	if err != nil {
		o.log.Printf("unregister %s: error: %s:", zoneIdentifier, err)
	}
	delete(o.zones, zoneIdentifier)

	return nil
}

func (o *Olympus) watchTimeout(logger ZoneLogger, hostName, zoneName string) {
	select {
	case <-logger.Done():
		return
	case <-logger.Timeouted():
		o.log.Printf("%s timeouted, unregistering", logger.ZoneIdentifier())
		o.mxClimate.Lock()
		defer o.mxClimate.Unlock()
		err := o.unregisterZone(hostName, zoneName, false)
		if err != nil {
			o.log.Printf("timeout unregistering error: %s", err)
		}
	}
}

func minInt(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func (o *Olympus) fetchBackClimateLog(logger ZoneLogger, c *rpc.Client, name string, size int) error {
	batchSize := 500

	for i := 0; i < size; i += batchSize {
		args := zeus.ZeusLogArgs{
			Start:    i,
			End:      minInt(i+batchSize, size),
			ZoneName: name,
		}
		reply := zeus.ZeusClimateLogReply{}
		if err := c.Call("Zeus.ClimateLog", args, &reply); err != nil {
			return err
		}
		for _, cr := range reply.Data {
			logger.ReportChannel() <- cr
		}
	}
	return nil
}

func (o *Olympus) fetchBackAlarmLog(logger ZoneLogger, c *rpc.Client, name string, size int) error {
	batchSize := 200

	for i := 0; i < size; i += batchSize {
		args := zeus.ZeusLogArgs{
			Start:    i,
			End:      minInt(i+batchSize, size),
			ZoneName: name,
		}
		reply := zeus.ZeusAlarmLogReply{}
		if err := c.Call("Zeus.AlarmLog", args, &reply); err != nil {
			return err
		}
		for _, ae := range reply.Data {
			logger.AlarmChannel() <- ae
		}
	}
	return nil
}

func (o *Olympus) fetchBackLogError(logger ZoneLogger, zr zeus.ZoneRegistration) (err error) {
	// race condition, since we are not behind the mutex anymore,
	// logger could be closed before we finish fetching logs. We just recover in that case.
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("zone was unregistered before backlogging terminated")
		}
	}()

	c, err := rpc.DialHTTP("tcp", zr.RPCAddress)
	if err != nil {
		return err
	}
	err = o.fetchBackClimateLog(logger, c, zr.Name, zr.SizeClimateLog)
	if err != nil {
		return err
	}
	return o.fetchBackAlarmLog(logger, c, zr.Name, zr.SizeAlarmLog)
}

func (o *Olympus) fetchBackLog(logger ZoneLogger, zr zeus.ZoneRegistration) {
	if len(zr.RPCAddress) == 0 || (zr.SizeClimateLog == 0 && zr.SizeAlarmLog == 0) {
		return
	}
	o.log.Printf("%s declared backlog data {ClimateReport:%d,AlarmEvent:%d}, fetching it", logger.ZoneIdentifier(), zr.SizeClimateLog, zr.SizeAlarmLog)
	if err := o.fetchBackLogError(logger, zr); err != nil {
		o.log.Printf("could not fetch backlog for %s: %s", logger.ZoneIdentifier(), err)
	}
}

func (o *Olympus) registerZone(zr zeus.ZoneRegistration) error {
	zoneIdentifier := zr.ZoneIdentifier()
	_, ok := o.zones[zoneIdentifier]
	if ok == true {
		return fmt.Errorf("%s is already registered", zoneIdentifier)
	}
	o.log.Printf("registering new zone %s", zoneIdentifier)
	o.climateLogger.Log(zr.Host+"."+zr.Name, true, true)
	logger := NewZoneLogger(zr)
	o.zones[zoneIdentifier] = logger
	go o.watchTimeout(logger, zr.Host, zr.Name)
	go o.fetchBackLog(logger, zr)
	return nil
}

func (o *Olympus) watchTracking(watcher TrackingWatcher) {
	select {
	case <-watcher.Done():
		return
	case <-watcher.Timeouted():
		o.log.Printf("%s timeouted, unregistering", watcher.Hostname())
		o.mxTracking.Lock()
		defer o.mxTracking.Unlock()
		err := o.unregisterTracker(watcher.Hostname(), false)
		if err != nil {
			o.log.Printf("timeout unregistering error: %s", err)
		}
	}
}

func (o *Olympus) registerTracker(hostname, streamURL string) error {
	if _, ok := o.watchers[hostname]; ok == true {
		return fmt.Errorf("tracking on %s is already marked on", hostname)
	}
	watcher := NewTrackingWatcher(hostname, streamURL)
	o.watchers[hostname] = watcher
	o.trackingLogger.Log(hostname, true, true)

	go o.watchTracking(watcher)
	return nil
}

func (o *Olympus) fetchOnlineTrackingError(streamServerBaseURL string) error {
	resp, err := http.Get(path.Join(streamServerBaseURL, "master.m3u8"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	reader := bufio.NewReader(resp.Body)
	for {
		l, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		l = strings.TrimSpace(l)
		ext := path.Ext(l)
		if len(l) == 0 || l[0] == '#' || ext != ".m3u8" {
			continue
		}
		URL := path.Clean(path.Join(streamServerBaseURL, l))
		host := strings.TrimSuffix(path.Base(l), ext)
		go func() {
			o.mxTracking.Lock()
			defer o.mxTracking.Unlock()
			o.registerTracker(host, URL)
		}()
	}
}

func (o *Olympus) fetchOnlineTracking(streamServerBaseURL string) {
	if err := o.fetchOnlineTrackingError(streamServerBaseURL); err != nil {
		o.log.Printf("could not fetch online tracker: %s", err)
	}
}

func (o *Olympus) unregisterTracker(hostname string, graceful bool) error {
	watcher, ok := o.watchers[hostname]
	if ok == false {
		return fmt.Errorf("tracking on %s is already not marked on", hostname)
	}
	err := watcher.Close()
	if err != nil {
		fmt.Errorf("could not close tracking watcher: %s", err)
	}
	o.trackingLogger.Log(hostname, false, graceful)
	delete(o.watchers, hostname)
	return nil

}

func (o *Olympus) route(router *mux.Router) {
	router.HandleFunc("/api/zones", func(w http.ResponseWriter, r *http.Request) {
		res := o.GetZones()
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}/climate", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetClimateReport(zeus.ZoneIdentifier(vars["hname"], vars["zname"]), r.URL.Query().Get("window"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}/alarms", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetAlarmEventLog(zeus.ZoneIdentifier(vars["hname"], vars["zname"]))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

	router.HandleFunc("/api/host/{hname}/zone/{zname}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := o.GetZone(zeus.ZoneIdentifier(vars["hname"], vars["zname"]))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONify(w, &res)
	}).Methods("GET")

}
