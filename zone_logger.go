package main

import (
	"math"
	"sync"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/proto"
)

type ZoneLogger interface {
	Host() string
	ZoneName() string
	ZoneIdentifier() string
	// PushTarget update current target for this logger.
	PushTarget(*proto.ClimateTarget)
	// PushReports updates a list of reports to this logger.
	PushReports([]*proto.ClimateReport)
	// PushAlarms adds a list of AlarmEvents to this logger.
	PushAlarms([]*proto.AlarmEvent)
	GetClimateTimeSeries(window string) ClimateTimeSeries
	GetAlarmReports() []AlarmReport
	GetClimateReport() ZoneClimateReport
}

const (
	climateTenMinute = iota
	climateHour
	climateDay
	climateWeek
	logs
	report
)

type zoneLogger struct {
	mx           sync.RWMutex
	sampler      ClimateReportSampler
	alarmReports map[string]*AlarmReport

	host, name    string
	currentReport ZoneClimateReport
}

func naNIfNil(v *float32) float32 {
	if v == nil {
		return float32(math.NaN())
	}
	return *v
}

func NewZoneLogger(declaration *proto.ZoneDeclaration) ZoneLogger {
	res := &zoneLogger{
		sampler:      NewClimateReportSampler(),
		host:         declaration.Host,
		name:         declaration.Name,
		alarmReports: make(map[string]*AlarmReport),
		currentReport: ZoneClimateReport{
			Temperature: nil,
			TemperatureBounds: Bounds{
				Min: declaration.MinTemperature,
				Max: declaration.MaxTemperature,
			},
			Humidity: nil,
			HumidityBounds: Bounds{
				Min: declaration.MinHumidity,
				Max: declaration.MaxHumidity,
			},
			NumAux: int(declaration.NumAux),
		},
	}
	return res
}

func (l *zoneLogger) PushReports(reports []*proto.ClimateReport) {
	if len(reports) == 0 {
		return
	}
	l.mx.Lock()
	defer l.mx.Unlock()

	l.sampler.Add(reports)

	lastReport := reports[len(reports)-1]
	if len(lastReport.Temperatures) > 0 {
		if l.currentReport.Temperature == nil {
			l.currentReport.Temperature = new(float32)
		}
		*l.currentReport.Temperature = lastReport.Temperatures[0]
	}
	if lastReport.Humidity != nil {
		if l.currentReport.Humidity == nil {
			l.currentReport.Humidity = new(float32)
		}
		*l.currentReport.Humidity = *lastReport.Humidity
	}
}

func (a *AlarmReport) On() bool {
	if len(a.Events) == 0 {
		return false
	}
	return a.Events[len(a.Events)-1].On
}

func (l *zoneLogger) PushAlarms(events []*proto.AlarmEvent) {
	l.mx.Lock()
	defer l.mx.Unlock()
	for _, e := range events {
		l.pushEventToLog(e)
	}
	l.updateActiveAlarmCounts()
}

func (l *zoneLogger) updateActiveAlarmCounts() {
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
func (l *zoneLogger) pushEventToLog(event *proto.AlarmEvent) {
	// insert sort the event, in most cases, it will simply append it
	r, ok := l.alarmReports[event.Reason]
	if ok == false {
		r = &AlarmReport{
			Reason: event.Reason,
			Level:  int(event.Level),
		}
		l.alarmReports[r.Reason] = r
	}
	r.Events = BackInsertionSort(r.Events,
		AlarmEvent{
			Time: event.Time.AsTime(),
			On:   (event.Status == proto.AlarmStatus_ALARM_ON),
		},
		func(a, b AlarmEvent) bool { return a.Time.Before(b.Time) })

}

func (l *zoneLogger) PushTarget(target *proto.ClimateTarget) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.currentReport.Current = target.Current
	l.currentReport.CurrentEnd = target.CurrentEnd
	if target.Next != nil && target.NextTime != nil {
		l.currentReport.Next = target.Next
		t := target.NextTime.AsTime()
		l.currentReport.NextTime = &t
		l.currentReport.NextEnd = target.NextEnd
	} else {
		l.currentReport.Next = nil
		l.currentReport.NextEnd = nil
		l.currentReport.NextTime = nil
	}
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

func (l *zoneLogger) fromWindow(window string) ClimateTimeSeries {
	windowType, ok := stringToRequestInt[window]
	if ok == false {
		windowType = climateTenMinute
	}

	switch windowType {
	case climateTenMinute:
		return l.sampler.LastTenMinutes()
	case climateHour:
		return l.sampler.LastHour()
	case climateDay:
		return l.sampler.LastDay()
	case climateWeek:
		return l.sampler.LastWeek()
	default:
		return l.sampler.LastTenMinutes()
	}
}

func (l *zoneLogger) GetClimateTimeSeries(window string) ClimateTimeSeries {
	l.mx.RLock()
	defer l.mx.RUnlock()
	// already a data copy, so it is safe
	return l.fromWindow(window)
}

func (l *zoneLogger) GetAlarmReports() []AlarmReport {
	l.mx.RLock()
	defer l.mx.RUnlock()
	res := make([]AlarmReport, 0, len(l.alarmReports))
	for _, report := range l.alarmReports {
		res = append(res, *report)
	}
	return res
}

func (l *zoneLogger) GetClimateReport() ZoneClimateReport {
	l.mx.RLock()
	defer l.mx.RUnlock()
	return deepcopy.MustAnything(l.currentReport).(ZoneClimateReport)
}

func (l *zoneLogger) Host() string {
	return l.host
}

func (l *zoneLogger) ZoneName() string {
	return l.name
}

func (l *zoneLogger) ZoneIdentifier() string {
	return ZoneIdentifier(l.host, l.name)
}
