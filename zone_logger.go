package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ZoneReport struct {
	Host              string
	Name              string
	NumAux            int
	Temperature       float64
	TemperatureBounds Bounds
	Humidity          float64
	HumidityBounds    Bounds
	Current           *zeus.State
	CurrentEnd        *zeus.State
	Next              *zeus.State
	NextEnd           *zeus.State
	NextTime          *time.Time
}

type ZoneLogger interface {
	StateChannel() chan<- zeus.StateReport
	ReportChannel() chan<- zeus.ClimateReport
	AlarmChannel() chan<- zeus.AlarmEvent
	Timeouted() <-chan struct{}
	Done() <-chan struct{}
	GetClimateReportSeries(window string) ClimateReportTimeSerie
	GetAlarmsEventLog() []zeus.AlarmEvent
	GetReport() ZoneReport
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
	logs    []zeus.AlarmEvent

	requests      chan namedRequest
	currentReport ZoneReport
}

func NewZoneLogger(reg zeus.ZoneRegistration) ZoneLogger {

	return &zoneLogger{
		done:     make(chan struct{}),
		timeout:  make(chan struct{}),
		states:   make(chan zeus.StateReport, 2),
		reports:  make(chan zeus.ClimateReport, 10),
		alarms:   make(chan zeus.AlarmEvent, 10),
		requests: make(chan namedRequest),
		sampler:  NewClimateReportSampler(reg.NumAux),
		currentReport: ZoneReport{
			Host:        reg.Host,
			Name:        reg.Name,
			NumAux:      reg.NumAux,
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
	}
}

func (l *zoneLogger) pushClimate(cr zeus.ClimateReport) {
	l.sampler.Add(cr)
	if len(cr.Temperatures) > 0 {
		l.currentReport.Temperature = float64(cr.Temperatures[0])
	}
	l.currentReport.Humidity = float64(cr.Humidity)
}

func (l *zoneLogger) pushLog(ae zeus.AlarmEvent) {
	ae.Zone = ""
	// insert sort the event, in most cases, it will simply append it
	i := BackLinearSearch(len(l.logs), func(i int) bool { return l.logs[i].Time.Before(ae.Time) }) + 1
	l.logs = append(l.logs, zeus.AlarmEvent{})
	copy(l.logs[i+1:], l.logs[i:])
	l.logs[i] = ae
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
		logs:             func() interface{} { return append([]zeus.AlarmEvent(nil), l.logs...) },
		report:           func() interface{} { return l.currentReport },
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
	tick := time.NewTicker(20 * time.Second)
	defer tick.Stop()
	for {
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
		if l.reports == nil && l.alarms == nil && l.states == nil {
			return
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

func (l *zoneLogger) GetAlarmsEventLog() []zeus.AlarmEvent {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: logs, result: returnChannel}
	res := <-returnChannel
	return res.([]zeus.AlarmEvent)
}

func (l *zoneLogger) GetReport() ZoneReport {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: report, result: returnChannel}
	res := <-returnChannel
	return res.(ZoneReport)
}

func (l *zoneLogger) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("already closed")
		}
		<-l.done
	}()

	close(l.reports)
	close(l.alarms)
	close(l.states)

	return nil
}
