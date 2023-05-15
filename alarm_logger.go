package main

import (
	"sync"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AlarmLogger interface {
	ActiveAlarmsCount() (warnings int, emergencies int)
	GetReports() []api.AlarmReport
	// PushAlarms adds a list of AlarmEvents to this logger.
	PushAlarms([]*api.AlarmEvent)
}
type alarmLogger struct {
	mx      sync.RWMutex
	reports map[string]*api.AlarmReport

	warnings, emergencies int
}

func NewAlarmLogger() AlarmLogger {
	return &alarmLogger{
		reports: make(map[string]*api.AlarmReport),
	}
}

func (l *alarmLogger) ActiveAlarmsCount() (int, int) {
	l.mx.RLock()
	defer l.mx.RUnlock()
	return l.warnings, l.emergencies
}

func (l *alarmLogger) GetReports() []api.AlarmReport {
	l.mx.RLock()
	defer l.mx.RUnlock()
	res := make([]api.AlarmReport, 0, len(l.reports))
	for _, report := range l.reports {
		res = append(res, *deepcopy.MustAnything(report).(*api.AlarmReport))
	}
	return res
}

func (l *alarmLogger) PushAlarms(events []*api.AlarmEvent) {
	l.mx.Lock()
	defer l.mx.Unlock()
	for _, e := range events {
		l.pushEventToLog(e)
	}
	l.computeActives()
}

func timestampBefore(a, b *timestamppb.Timestamp) bool {
	if a.Seconds == b.Seconds {
		return a.Nanos < b.Nanos
	}
	return a.Seconds < b.Seconds
}

func (l *alarmLogger) pushEventToLog(event *api.AlarmEvent) {
	report, ok := l.reports[event.Identification]
	if ok == false {
		report = &api.AlarmReport{
			Identification: event.Identification,
			Level:          event.Level,
		}
		l.reports[report.Identification] = report
	}
	report.Events = BackInsertionSort(report.Events,
		api.AlarmTimePoint{
			Time: event.Time.AsTime(),
			On:   event.Status == api.AlarmStatus_ON,
		},
		func(a, b api.AlarmTimePoint) bool {
			return a.Time.Before(b.Time)
		})

	lastTimePoint := report.Events[len(report.Events)-1]
	if lastTimePoint.Time.Equal(event.Time.AsTime()) && len(event.Description) == 0 {
		report.Description = event.Description
	}
}

func (l *alarmLogger) computeActives() {
	l.emergencies = 0
	l.warnings = 0
	for _, report := range l.reports {
		lastTimepoint := report.Events[len(report.Events)-1]
		if lastTimepoint.On == false {
			continue
		}
		switch report.Level {
		case api.AlarmLevel_EMERGENCY:
			l.emergencies += 1
		case api.AlarmLevel_WARNING:
			l.warnings += 1
		}
	}
}
