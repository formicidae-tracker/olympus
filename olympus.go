package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ZoneNotFoundError struct {
	Zone string
}

func (z ZoneNotFoundError) Error() string {
	return fmt.Sprintf("olympus: unknown zone '%s'", z.Zone)
}

type ServiceEvent struct {
	ZoneIdentifier string
	Time           time.Time
	On             bool
	Graceful       bool
}

type Olympus struct {
	mx    *sync.RWMutex
	zones map[string]ZoneLogger
	log   *log.Logger
}

func (o *Olympus) GetZones() []ZoneReportSummary {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.getZones()
}

func (o *Olympus) GetClimateReport(zoneIdentifier, window string) (ClimateReportTimeSerie, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.getClimateReport(zoneIdentifier, window)
}

func (o *Olympus) GetZone(zoneIdentifier string) (ZoneReport, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.getZoneReport(zoneIdentifier)
}

func (o *Olympus) GetAlarmEventLog(zoneIdentifier string) ([]AlarmEvent, error) {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.getAlarmEventLog(zoneIdentifier)
}

func (o *Olympus) RegisterZone(zr *zeus.ZoneRegistration, unused *int) error {
	o.mx.Lock()
	defer o.mx.Unlock()

	return o.registerZone(zr)
}

func (o *Olympus) UnregisterZone(zr *zeus.ZoneUnregistration, unused *int) error {
	o.mx.Lock()
	defer o.mx.Unlock()

	return o.unregisterZone(zeus.ZoneIdentifier(zr.Host, zr.Name), true)
}

func (o *Olympus) ReportClimate(cr *zeus.NamedClimateReport, unused *int) error {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.reportClimate(*cr)
}

func (o *Olympus) ReportAlarm(ae *zeus.AlarmEvent, unused *int) error {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.reportAlarm(*ae)
}

func (o *Olympus) ReportState(sr *zeus.StateReport, unused *int) error {
	o.mx.RLock()
	defer o.mx.RUnlock()

	return o.reportState(*sr)
}

func (o *Olympus) RegisterTracker(hostname string, unused *int) error {
	o.mx.Lock()
	defer o.mx.RUnlock()

	return fmt.Errorf("not yet implemented")
}

func (o *Olympus) UnregisterTracker(hostname string, unused *int) error {
	o.mx.Lock()
	defer o.mx.RUnlock()

	return fmt.Errorf("not yet implemented")
}

func (o *Olympus) getZones() []ZoneReportSummary {
	res := make([]ZoneReportSummary, 0, len(o.zones))
	for _, z := range o.zones {
		res = append(res, z.GetReport().ZoneReportSummary)
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

func (o *Olympus) getZoneReport(zoneIdentifier string) (ZoneReport, error) {
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return ZoneReport{}, ZoneNotFoundError{zoneIdentifier}
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

func (o *Olympus) unregisterZone(zoneIdentifier string, graceful bool) error {
	z, ok := o.zones[zoneIdentifier]
	if ok == false {
		return ZoneNotFoundError{zoneIdentifier}
	}
	o.log.Printf("unregistering %s", zoneIdentifier)
	err := z.Close()
	o.log.Printf("unregister %s: error: %s:", zoneIdentifier, err)
	delete(o.zones, zoneIdentifier)

	return nil
}

func (o *Olympus) watchTimeout(logger ZoneLogger) {
	select {
	case <-logger.Done():
		return
	case <-logger.Timeouted():
		o.log.Printf("%s timeouted, unregistering", logger.Fullname())
		o.mx.Lock()
		defer o.mx.Unlock()
		err := o.unregisterZone(logger.Fullname(), false)
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

func (o *Olympus) fetchBackLogError(logger ZoneLogger, zr *zeus.ZoneRegistration) (err error) {
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

func (o *Olympus) fetchBackLog(logger ZoneLogger, zr *zeus.ZoneRegistration) {
	if len(zr.RPCAddress) == 0 || (zr.SizeClimateLog == 0 && zr.SizeAlarmLog == 0) {
		return
	}
	if err := o.fetchBackLogError(logger, zr); err != nil {
		o.log.Printf("could not fetch backlog for %s: %s", logger.Fullname(), err)
	}
}

func (o *Olympus) registerZone(zr *zeus.ZoneRegistration) error {
	zoneIdentifier := zr.ZoneIdentifier()
	_, ok := o.zones[zoneIdentifier]
	if ok == true {
		return fmt.Errorf("%s is already registered", zoneIdentifier)
	}
	logger := NewZoneLogger(*zr, 30*time.Second)
	o.zones[zoneIdentifier] = logger
	go o.watchTimeout(logger)
	go o.fetchBackLog(logger, zr)
	return nil
}

func NewOlympus() *Olympus {
	return &Olympus{
		mx:    &sync.RWMutex{},
		zones: make(map[string]ZoneLogger),
		log:   log.New(os.Stderr, "[olympus] :", log.LstdFlags),
	}
}
