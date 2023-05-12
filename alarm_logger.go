package main

import (
	"sync"

	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AlarmLogger interface {
	ActiveAlarmsCount() (warnings int32, emergencies int32)
	GetReports() []*api.AlarmReport
	// PushAlarms adds a list of AlarmEvents to this logger.
	PushAlarms([]*api.AlarmEvent)
}
type alarmLogger struct {
	mx      sync.RWMutex
	reports map[string]*api.AlarmReport

	warnings, emergencies int32
}

func NewAlarmLogger() AlarmLogger {
	return &alarmLogger{
		reports: make(map[string]*api.AlarmReport),
	}
}

func (l *alarmLogger) ActiveAlarmsCount() (int32, int32) {
	return l.warnings, l.emergencies
}

func (l *alarmLogger) GetReports() []*api.AlarmReport {
	return nil
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

func timestampEqual(a, b *timestamppb.Timestamp) bool {
	return a.Seconds == b.Seconds && a.Nanos == b.Nanos
}

func (l *alarmLogger) pushEventToLog(event *api.AlarmEvent) {
	if event.Identification == nil || event.Level == nil {
		return
	}
	report, ok := l.reports[*event.Identification]
	if ok == false {
		report = &api.AlarmReport{
			Identification: *event.Identification,
			Level:          *event.Level,
		}
		l.reports[report.Identification] = report
	}
	report.Events = BackInsertionSort(report.Events,
		&api.AlarmEvent{
			Time:   event.Time,
			Status: event.Status,
		},
		func(a, b *api.AlarmEvent) bool {
			return timestampBefore(a.Time, b.Time)
		})

	lastEvent := report.Events[len(report.Events)-1]
	if timestampEqual(lastEvent.Time, event.Time) && event.Description != nil {
		report.Description = *event.Description
	}
}

func (l *alarmLogger) computeActives() {
	l.emergencies = 0
	l.warnings = 0
	for _, report := range l.reports {
		lastEvent := report.Events[len(report.Events)-1]
		if lastEvent.Status != api.AlarmStatus_ON {
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
