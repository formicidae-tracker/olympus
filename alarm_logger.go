package main

import (
	"sync"
	"time"

	"github.com/formicidae-tracker/olympus/api"
)

type AlarmLogger interface {
	ActiveAlarmsCount() (warnings int, emergencies int)
	GetReports() []api.AlarmReport
	// PushAlarms adds a list of AlarmEvents to this logger.
	PushAlarms([]*api.AlarmUpdate)
}

type alarmTimePoint struct {
	time time.Time
	on   bool
}

type alarmLog struct {
	identification string
	level          api.AlarmLevel
	description    string
	logs           []alarmTimePoint
}

func (l *alarmLog) getReport() api.AlarmReport {
	res := api.AlarmReport{
		Identification: l.identification,
		Level:          l.level,
		Description:    l.description,
		Events:         make([]api.AlarmEvent, 0, len(l.logs)/2+1),
	}
	var start *time.Time = nil
	for _, u := range l.logs {
		if start == nil {
			if u.on == false {
				continue
			}
			start = &u.time
		} else {
			if u.on == true {
				continue
			}
			e := api.AlarmEvent{
				Start: *start,
				End:   new(time.Time),
			}
			*e.End = u.time
			res.Events = append(res.Events, e)
			start = nil
		}
	}

	if start != nil {
		res.Events = append(res.Events, api.AlarmEvent{Start: *start})
	}

	return res
}

func (l *alarmLog) on() bool {
	if len(l.logs) == 0 {
		return false
	}
	return l.logs[len(l.logs)-1].on
}

func (l *alarmLog) pushUpdate(u *api.AlarmUpdate) {
	updateTime := u.Time.AsTime()
	l.logs = BackInsertionSort(l.logs,
		alarmTimePoint{
			time: updateTime,
			on:   u.Status == api.AlarmStatus_ON,
		},
		func(a, b alarmTimePoint) bool {
			return a.time.Before(b.time)
		})

	updateHasDescription := len(u.Description) != 0
	hasNoDescription := len(l.description) == 0
	isLastUpdate := l.logs[len(l.logs)-1].time.Equal(updateTime)
	if updateHasDescription && (hasNoDescription || isLastUpdate) {
		l.description = u.Description
	}
}

type alarmLogger struct {
	mx   sync.RWMutex
	logs map[string]*alarmLog

	warnings, emergencies int
}

func NewAlarmLogger() AlarmLogger {
	return &alarmLogger{
		logs: make(map[string]*alarmLog),
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
	res := make([]api.AlarmReport, 0, len(l.logs))
	for _, log := range l.logs {
		res = append(res, log.getReport())
	}
	return res
}

func (l *alarmLogger) PushAlarms(updates []*api.AlarmUpdate) {
	l.mx.Lock()
	defer l.mx.Unlock()
	for _, e := range updates {
		l.pushUpdateToLog(e)
	}
	l.computeActives()
}

func (l *alarmLogger) pushUpdateToLog(update *api.AlarmUpdate) {
	log, ok := l.logs[update.Identification]
	if ok == false {
		log = &alarmLog{
			identification: update.Identification,
			level:          update.Level,
		}
		l.logs[log.identification] = log
	}
	log.pushUpdate(update)
}

func (l *alarmLogger) computeActives() {
	l.emergencies = 0
	l.warnings = 0
	for _, log := range l.logs {
		if log.on() == false {
			continue
		}
		switch log.level {
		case api.AlarmLevel_EMERGENCY:
			l.emergencies += 1
		case api.AlarmLevel_WARNING:
			l.warnings += 1
		}
	}
}
