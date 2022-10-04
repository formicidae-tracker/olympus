package main

import (
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/olympuspb"
)

type ZoneLogger interface {
	Host() string
	ZoneName() string
	ZoneIdentifier() string
	// PushTarget update current target for this logger.
	PushTarget(*olympuspb.ClimateTarget)
	// PushReports updates a list of reports to this logger.
	PushReports([]*olympuspb.ClimateReport)
	// PushAlarms adds a list of AlarmEvents to this logger.
	PushAlarms([]*olympuspb.AlarmEvent)
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
	alarmReports map[string]*AlarmReport

	host, name     string
	currentReport  ZoneClimateReport
	lastReportTime time.Time

	samplers         []ClimateDataDownsampler
	samplersByWindow map[string]ClimateDataDownsampler
}

func NewZoneLogger(declaration *olympuspb.ZoneDeclaration) ZoneLogger {
	samplers := []ClimateDataDownsampler{
		NewClimateDataDownsampler(10*time.Minute, 500),
		NewClimateDataDownsampler(1*time.Hour, 400),
		NewClimateDataDownsampler(24*time.Hour, 300),
		NewClimateDataDownsampler(7*24*time.Hour, 300),
	}
	samplersByWindow := map[string]ClimateDataDownsampler{
		"10-minute":  samplers[0],
		"10m":        samplers[0],
		"10-minutes": samplers[0],
		"1h":         samplers[1],
		"hour":       samplers[1],
		"1d":         samplers[2],
		"day":        samplers[2],
		"1w":         samplers[3],
		"week":       samplers[3],
	}

	res := &zoneLogger{
		samplers:         samplers,
		samplersByWindow: samplersByWindow,
		host:             declaration.Host,
		name:             declaration.Name,
		alarmReports:     make(map[string]*AlarmReport),
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
		},
	}
	return res
}

func buildBatch(reports []*olympuspb.ClimateReport) TimedValues {
	if len(reports) == 0 {
		return TimedValues{}
	}
	times := make([]time.Time, len(reports))
	values := make([][]float32, 0)

	appendValue := func(values []float32, v float32, size int) []float32 {
		missing := v
		if len(values) > 0 {
			missing = values[len(values)-1]
		}
		for len(values) < size-1 {
			values = append(values, missing)
		}
		return append(values, v)
	}

	for i, r := range reports {
		times[i] = r.Time.AsTime()
		if r.Humidity != nil {
			if len(values) == 0 {
				values = append(values, nil)
			}
			values[0] = appendValue(values[0], *r.Humidity, i+1)
		}
		for j, t := range r.Temperatures {
			for len(values) <= j+1 {
				values = append(values, nil)
			}
			values[j+1] = appendValue(values[j+1], t, i+1)
		}
	}
	return TimedValues{times: times, values: values}
}

func (l *zoneLogger) PushReports(reports []*olympuspb.ClimateReport) {
	if len(reports) == 0 {
		return
	}
	l.mx.Lock()
	defer l.mx.Unlock()

	toAdd := buildBatch(reports)
	for _, s := range l.samplers {
		s.Add(toAdd)
	}

	lastReport := reports[len(reports)-1]
	lastReportTime := lastReport.Time.AsTime()
	if lastReportTime.After(l.lastReportTime) == false {
		return
	}
	l.lastReportTime = lastReportTime

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

func (l *zoneLogger) PushAlarms(events []*olympuspb.AlarmEvent) {
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
		if olympuspb.AlarmLevel(r.Level) == olympuspb.AlarmLevel_WARNING {
			l.currentReport.ActiveWarnings += 1
		} else {
			l.currentReport.ActiveEmergencies += 1
		}
	}

}
func (l *zoneLogger) pushEventToLog(event *olympuspb.AlarmEvent) {
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
			On:   (event.Status == olympuspb.AlarmStatus_ON),
		},
		func(a, b AlarmEvent) bool { return a.Time.Before(b.Time) })

}

func (l *zoneLogger) PushTarget(target *olympuspb.ClimateTarget) {
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

func (l *zoneLogger) fromWindow(window string) ClimateTimeSeries {
	sampler, ok := l.samplersByWindow[window]
	if ok == false {
		sampler = l.samplers[0]
	}
	return sampler.TimeSeries()
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
