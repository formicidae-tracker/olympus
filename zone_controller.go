package main

import (
	"fmt"
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
}

func (l *zoneLogger) pushLog(ae zeus.AlarmEvent) {
}

func (l *zoneLogger) pushState(ae zeus.StateReport) {
}

func (l *zoneLogger) handleRequest(r namedRequest) {
	close(r.result)
}

func (l *zoneLogger) mainLoop() {
	defer close(l.done)

	once := sync.Once{}
	dp, deadline := NewDeadlinePusher(20 * time.Second)
	for {
		select {
		case cr, ok := <-l.reports:
			if ok == false {
				l.reports = nil
				continue
			}
			deadline = dp.Push()
			l.pushClimate(cr)
		case ae, ok := <-l.alarms:
			if ok == false {
				l.alarms = nil
				continue
			}
			deadline = dp.Push()
			l.pushLog(ae)
		case sr, ok := <-l.states:
			if ok == false {
				l.states = nil
				continue
			}
			deadline = dp.Push()
			l.pushState(sr)
		case req := <-l.requests:
			l.handleRequest(req)
		case now := <-deadline:
			deadline = dp.Check(now)
			if deadline != nil {
				continue
			}
			// we need to do it
			once.Do(func() {
				close(l.timeout)
			})
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
