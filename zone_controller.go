package main

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ZoneLogger interface {
	StateChannel() chan<- zeus.StateReport
	ReportChannel() chan<- zeus.ClimateReport
	AlarmChannel() chan<- zeus.AlarmEvent
	Timeouted() <-chan struct{}
	Done() <-chan struct{}
	GetClimateReportSeries(window string) ClimateReportTimeSerie
	GetAlarmsEventLog() []zeus.AlarmEvent
	Close() error
}

type namedRequest struct {
	request string
	result  chan interface{}
}

type zoneLogger struct {
	done, timeout chan struct{}
	states        chan zeus.StateReport
	reports       chan zeus.ClimateReport
	alarms        chan zeus.AlarmEvent

	sampler ClimateReportSampler
	logs    []zeus.AlarmEvent

	requests chan namedRequest
}

func NewZoneLogger() ZoneLogger {
	return &zoneLogger{
		done:     make(chan struct{}),
		timeout:  make(chan struct{}),
		states:   make(chan zeus.StateReport, 2),
		reports:  make(chan zeus.ClimateReport, 10),
		alarms:   make(chan zeus.AlarmEvent, 10),
		requests: make(chan namedRequest),
	}

}

func (l *zoneLogger) pushClimate(cr zeus.ClimateReport) {
	l.sampler.Add(cr)
}

func (l *zoneLogger) pushLog(ae zeus.AlarmEvent) {
	needSort := false
	if len(l.logs) > 0 {
		needSort = l.logs[len(l.logs)-1].Time.After(ae.Time)
	}
	ae.Zone = ""
	l.logs = append(l.logs, ae)
	if needSort == true {
		sort.Slice(l.logs, func(i int, j int) bool {
			return l.logs[i].Time.Before(l.logs[j].Time)
		})
	}
}

func (l *zoneLogger) pushState(ae zeus.StateReport) {

}

func (l *zoneLogger) handleRequest(r namedRequest) {
	close(r.result)
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

func (l *zoneLogger) GetClimateReportSeries(window string) ClimateReportTimeSerie {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{result: returnChannel}
	res := <-returnChannel
	return res.(ClimateReportTimeSerie)
}

func (l *zoneLogger) GetAlarmsEventLog() []zeus.AlarmEvent {
	returnChannel := make(chan interface{})
	l.requests <- namedRequest{result: returnChannel}
	res := <-returnChannel
	return res.([]zeus.AlarmEvent)
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
