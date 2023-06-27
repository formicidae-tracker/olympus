package olympus

import (
	"sync"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
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
	timepoints     []alarmTimePoint
}

func (l *alarmLog) getReport() api.AlarmReport {
	return api.AlarmReport{
		Identification: l.identification,
		Level:          l.level,
		Description:    l.description,
		Events:         l.buildEvents(),
	}
}

func (l *alarmLog) buildEvents() []api.AlarmEvent {
	events := make([]api.AlarmEvent, 0, len(l.timepoints)/2+1)
	var start *time.Time = nil
	for _, tp := range l.timepoints {
		if start == nil {
			if tp.on == false {
				continue
			}
			start = &tp.time
		} else {
			if tp.on == true {
				continue
			}
			e := api.AlarmEvent{
				Start: *start,
				End:   new(time.Time),
			}
			*e.End = tp.time
			events = append(events, e)
			start = nil
		}
	}

	if start != nil {
		events = append(events, api.AlarmEvent{Start: *start})
	}

	return events
}

func (l *alarmLog) needDecimate() bool {
	lastOn := false
	for _, tp := range l.timepoints {
		if tp.on == lastOn {
			return true
		}
		lastOn = tp.on
	}
	return false
}

func (l *alarmLog) decimate() {
	if l.needDecimate() == false {
		return
	}
	events := l.buildEvents()
	l.timepoints = make([]alarmTimePoint, 0, len(events)*2)
	for _, e := range events {
		l.timepoints = append(l.timepoints, alarmTimePoint{time: e.Start, on: true})
		if e.End != nil {
			l.timepoints = append(l.timepoints, alarmTimePoint{time: *e.End, on: false})
		}
	}
}

func (l *alarmLog) on() bool {
	if len(l.timepoints) == 0 {
		return false
	}
	return l.timepoints[len(l.timepoints)-1].on
}

func (l *alarmLog) pushUpdate(u *api.AlarmUpdate) {
	updateTime := u.Time.AsTime()
	l.timepoints = BackInsertionSort(l.timepoints,
		alarmTimePoint{
			time: updateTime,
			on:   u.Status == api.AlarmStatus_ON,
		},
		func(a, b alarmTimePoint) bool {
			return a.time.Before(b.time)
		})

	updateHasDescription := len(u.Description) != 0
	hasNoDescription := len(l.description) == 0
	isLastUpdate := l.timepoints[len(l.timepoints)-1].time.Equal(updateTime)
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
	l.decimateLogs(updates)
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

func (l *alarmLogger) decimateLogs(updates []*api.AlarmUpdate) {
	idts := make(map[string]bool)
	for _, u := range updates {
		idts[u.Identification] = true
	}

	for idt := range idts {
		l.logs[idt].decimate()
	}
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
