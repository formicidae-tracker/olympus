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
	GetClimateTimeSeries(window string) ClimateTimeSerie
	GetAlarmReports() []AlarmReport
	GetClimateReport() *ZoneClimateReport
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

	sampler      ClimateReportSampler
	alarmReports map[string]*AlarmReport

	requests      chan namedRequest
	host, name    string
	currentReport ZoneClimateReport

	timeoutPeriod time.Duration
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
		alarmReports:  make(map[string]*AlarmReport),
		currentReport: ZoneClimateReport{
			Temperature: -1000.0,
			TemperatureBounds: Bounds{
				Min: reg.MinTemperature,
				Max: reg.MaxTemperature,
			},
			Humidity: -1000.0,
			HumidityBounds: Bounds{
				Min: reg.MinHumidity,
				Max: reg.MaxHumidity,
			},
			NumAux: reg.NumAux,
		},
	}
	go res.mainLoop()
	return res
}

func (l *zoneLogger) pushClimate(cr zeus.ClimateReport) {
	l.sampler.Add(cr)
	if len(cr.Temperatures) > 0 && zeus.IsUndefined(cr.Temperatures[0]) == false {
		l.currentReport.Temperature = float64(cr.Temperatures[0])
	}
	if zeus.IsUndefined(cr.Humidity) == false {
		l.currentReport.Humidity = cr.Humidity.Value()
	}
}

func eventInsertionSort(events []AlarmEvent, ae AlarmEvent) []AlarmEvent {
	i := BackLinearSearch(len(events), func(i int) bool { return events[i].Time.Before(ae.Time) }) + 1
	events = append(events, AlarmEvent{})
	copy(events[i+1:], events[i:])
	events[i] = ae
	return events
}

func (a *AlarmReport) On() bool {
	if len(a.Events) == 0 {
		return false
	}
	return a.Events[len(a.Events)-1].On
}

func (l *zoneLogger) pushLog(ae zeus.AlarmEvent) {
	// insert sort the event, in most cases, it will simply append it
	r, ok := l.alarmReports[ae.Reason]
	if ok == false {
		r = &AlarmReport{
			Reason: ae.Reason,
			Level:  zeus.MapPriority(ae.Flags),
		}
		l.alarmReports[r.Reason] = r
	}
	r.Events = eventInsertionSort(r.Events, AlarmEvent{
		Time: ae.Time,
		On:   (ae.Status == zeus.AlarmOn),
	})
	l.currentReport.ActiveEmergencies = 0
	l.currentReport.ActiveWarnings = 0
	for _, r := range l.alarmReports {
		if r.On() == false {
			continue
		}
		if r.Level == 1 {
			l.currentReport.ActiveWarnings += 1
		} else {
			l.currentReport.ActiveEmergencies += 1
		}
	}

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

func (l *zoneLogger) copyReports() []AlarmReport {
	reports := make([]AlarmReport, 0, len(l.alarmReports))
	for _, r := range l.alarmReports {
		reports = append(reports, *r)
	}
	return reports
}

func (l *zoneLogger) handleRequest(r namedRequest) {
	requestHandlers := map[int]func() interface{}{
		climateTenMinute: func() interface{} { return l.sampler.LastTenMinutes() },
		climateHour:      func() interface{} { return l.sampler.LastHour() },
		climateDay:       func() interface{} { return l.sampler.LastDay() },
		climateWeek:      func() interface{} { return l.sampler.LastWeek() },
		logs:             func() interface{} { return l.copyReports() },
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

func (l *zoneLogger) GetClimateTimeSeries(window string) ClimateTimeSerie {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: l.fromWindow(window), result: returnChannel}
	res := <-returnChannel
	return res.(ClimateTimeSerie)
}

func (l *zoneLogger) GetAlarmReports() []AlarmReport {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: logs, result: returnChannel}
	res := <-returnChannel
	return res.([]AlarmReport)
}

func (l *zoneLogger) GetClimateReport() *ZoneClimateReport {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{request: report, result: returnChannel}
	res := <-returnChannel
	return res.(*ZoneClimateReport)
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
