package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ZoneLogger interface {
	Host() string
	ZoneName() string
	ZoneIdentifier() string
	StateChannel() chan<- zeus.StateReport
	ReportChannel() chan<- zeus.ClimateReport
	AlarmChannel() chan<- zeus.AlarmEvent
	Timeouted() <-chan struct{}
	Done() <-chan struct{}
	GetClimateReportSeries(window string) ClimateReportTimeSerie
	GetAlarmsEventLog() []AlarmEvent
	GetReport() ZoneClimateReport
	Close() error
}

const (
	climateTenMinute = iota
	climateHour
	climateDay
	climateWeek
	logs
	report
)

type namedRequest struct {
	request int
	result  chan interface{}
}

type zoneLogger struct {
	done, timeout chan struct{}
	states        chan zeus.StateReport
	reports       chan zeus.ClimateReport
	alarms        chan zeus.AlarmEvent

	sampler ClimateReportSampler
	logs    []AlarmEvent

	requests      chan namedRequest
	host, name    string
	currentReport ZoneClimateReport

	timeoutPeriod time.Duration

	last                  map[string]int
	warnings, emergencies map[string]bool
}

func NewZoneLogger(reg zeus.ZoneRegistration) ZoneLogger {
	return newZoneLogger(reg, 30*time.Second)
}

func newZoneLogger(reg zeus.ZoneRegistration, timeoutPeriod time.Duration) ZoneLogger {

	res := &zoneLogger{
		done:          make(chan struct{}),
		timeout:       make(chan struct{}),
		states:        make(chan zeus.StateReport, 2),
		reports:       make(chan zeus.ClimateReport, 10),
		alarms:        make(chan zeus.AlarmEvent, 10),
		requests:      make(chan namedRequest),
		sampler:       NewClimateReportSampler(reg.NumAux),
		timeoutPeriod: timeoutPeriod,
		host:          reg.Host,
		name:          reg.Name,
		last:          make(map[string]int),
		warnings:      make(map[string]bool),
		emergencies:   make(map[string]bool),

		currentReport: ZoneClimateReport{
			ZoneClimateStatus: ZoneClimateStatus{
				Temperature: 0.0,
				TemperatureBounds: Bounds{
					Min: reg.MinTemperature,
					Max: reg.MaxTemperature,
				},
				Humidity: 0.0,
				HumidityBounds: Bounds{
					Min: reg.MinHumidity,
					Max: reg.MaxHumidity,
				},
			},
			NumAux: reg.NumAux,
		},
	}
	go res.mainLoop()
	return res
}

func (l *zoneLogger) pushClimate(cr zeus.ClimateReport) {
	l.sampler.Add(cr)
	if len(cr.Temperatures) > 0 {
		l.currentReport.Temperature = float64(cr.Temperatures[0])
	}
	l.currentReport.Humidity = float64(cr.Humidity)
}

func (l *zoneLogger) pushLog(ae zeus.AlarmEvent) {
	// insert sort the event, in most cases, it will simply append it
	i := BackLinearSearch(len(l.logs), func(i int) bool { return l.logs[i].Time.Before(ae.Time) }) + 1
	l.logs = append(l.logs, AlarmEvent{})
	copy(l.logs[i+1:], l.logs[i:])
	l.logs[i] = AlarmEvent{
		Reason: ae.Reason,
		Level:  zeus.MapPriority(ae.Flags),
		On:     ae.Status == zeus.AlarmOn,
		Time:   ae.Time,
	}
	if i < l.last[ae.Reason] {
		return
	}
	l.last[ae.Reason] = i
	set := l.warnings
	if ae.Flags&zeus.Emergency != 0 {
		set = l.emergencies
	}
	if ae.Status == zeus.AlarmOn {
		set[ae.Reason] = true
	} else {
		delete(set, ae.Reason)
	}
	l.currentReport.ActiveEmergencies = len(l.emergencies)
	l.currentReport.ActiveWarnings = len(l.warnings)
}

func (l *zoneLogger) pushState(sr zeus.StateReport) {
	l.currentReport.Current = &sr.Current
	l.currentReport.CurrentEnd = sr.CurrentEnd
	if sr.Next != nil && sr.NextTime != nil {
		l.currentReport.Next = sr.Next
		l.currentReport.NextTime = sr.NextTime
		l.currentReport.NextEnd = sr.NextEnd
	} else {
		l.currentReport.Next = nil
		l.currentReport.NextEnd = nil
		l.currentReport.NextTime = nil
	}
}

func (l *zoneLogger) handleRequest(r namedRequest) {
	requestHandlers := map[int]func() interface{}{
		climateTenMinute: func() interface{} { return l.sampler.LastTenMinutes() },
		climateHour:      func() interface{} { return l.sampler.LastHour() },
		climateDay:       func() interface{} { return l.sampler.LastDay() },
		climateWeek:      func() interface{} { return l.sampler.LastWeek() },
		logs:             func() interface{} { return append([]AlarmEvent(nil), l.logs...) },
		report:           func() interface{} { return l.currentReport.makeCopy() },
	}
	defer close(r.result)
	h, ok := requestHandlers[r.request]
	if ok == false {
		return
	}
	r.result <- h()
}

func (l *zoneLogger) mainLoop() {
	defer close(l.done)

	once := sync.Once{}
	seen := false
	tick := time.NewTicker(l.timeoutPeriod)
	defer tick.Stop()
	for {
		if l.reports == nil && l.alarms == nil && l.states == nil {
			return
		}

		select {
		case cr, ok := <-l.reports:
			if ok == false {
				l.reports = nil
				continue
			}
			seen = true
			l.pushClimate(cr)
		case ae, ok := <-l.alarms:
			if ok == false {
				l.alarms = nil
				continue
			}
			seen = true
			l.pushLog(ae)
		case sr, ok := <-l.states:
			if ok == false {
				l.states = nil
				continue
			}
			seen = true
			l.pushState(sr)
		case req := <-l.requests:
			l.handleRequest(req)
		case <-tick.C:
			if seen == false {
				once.Do(func() {
					close(l.timeout)
				})
			}
			seen = false
		}
	}
}

func (l *zoneLogger) StateChannel() chan<- zeus.StateReport {
	return l.states
}

func (l *zoneLogger) ReportChannel() chan<- zeus.ClimateReport {
	return l.reports
}

func (l *zoneLogger) AlarmChannel() chan<- zeus.AlarmEvent {
	return l.alarms
}

func (l *zoneLogger) Timeouted() <-chan struct{} {
	return l.timeout
}

func (l *zoneLogger) Done() <-chan struct{} {
	return l.done
}

var stringToRequestInt = map[string]int{
	"10-minute":  climateTenMinute,
	"10m":        climateTenMinute,
	"10-minutes": climateTenMinute,
	"1h":         climateHour,
	"hour":       climateHour,
	"1d":         climateDay,
	"day":        climateDay,
	"1w":         climateWeek,
	"week":       climateWeek,
}

func (l *zoneLogger) fromWindow(window string) int {
	rValue, ok := stringToRequestInt[window]
	if ok == false {
		return climateTenMinute
	}
	return rValue
}

func (l *zoneLogger) GetClimateReportSeries(window string) ClimateReportTimeSerie {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: l.fromWindow(window), result: returnChannel}
	res := <-returnChannel
	return res.(ClimateReportTimeSerie)
}

func (l *zoneLogger) GetAlarmsEventLog() []AlarmEvent {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: logs, result: returnChannel}
	res := <-returnChannel
	return res.([]AlarmEvent)
}

func (l *zoneLogger) GetReport() ZoneClimateReport {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: report, result: returnChannel}
	res := <-returnChannel
	return res.(ZoneClimateReport)
}

func (l *zoneLogger) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("ZoneLogger: already closed")
		}
		<-l.done
	}()

	close(l.reports)
	close(l.alarms)
	close(l.states)

	return nil
}

func (l *zoneLogger) Host() string {
	return l.host
}

func (l *zoneLogger) ZoneName() string {
	return l.name
}

func (l *zoneLogger) ZoneIdentifier() string {
	return zeus.ZoneIdentifier(l.host, l.name)
}
