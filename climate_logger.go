package main

import (
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
)

type ClimateLogger interface {
	Host() string
	ZoneName() string
	ZoneIdentifier() string
	// PushTarget update current target for this logger.
	PushTarget(*api.ClimateTarget)
	// PushReports updates a list of reports to this logger.
	PushReports([]*api.ClimateReport)
	GetClimateTimeSeries(window string) api.ClimateTimeSeries
	GetClimateReport() *api.ZoneClimateReport
}

const (
	climateTenMinute = iota
	climateHour
	climateDay
	climateWeek
	logs
	report
)

type climateLogger struct {
	mx sync.RWMutex

	host, name     string
	currentReport  *api.ZoneClimateReport
	lastReportTime time.Time

	samplers         []ClimateDataDownsampler
	samplersByWindow map[string]ClimateDataDownsampler
}

func NewClimateLogger(declaration *api.ClimateDeclaration) ClimateLogger {
	samplers := []ClimateDataDownsampler{
		NewClimateDataDownsampler(10*time.Minute, time.Minute, 500),
		NewClimateDataDownsampler(1*time.Hour, time.Minute, 400),
		NewClimateDataDownsampler(24*time.Hour, time.Hour, 300),
		NewClimateDataDownsampler(7*24*time.Hour, 24*time.Hour, 300),
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

	res := &climateLogger{
		samplers:         samplers,
		samplersByWindow: samplersByWindow,
		host:             declaration.Host,
		name:             declaration.Name,
		currentReport: &api.ZoneClimateReport{
			Temperature: nil,
			TemperatureBounds: api.Bounds{
				Minimum: declaration.MinTemperature,
				Maximum: declaration.MaxTemperature,
			},
			Humidity: nil,
			HumidityBounds: api.Bounds{
				Minimum: declaration.MinHumidity,
				Maximum: declaration.MaxHumidity,
			},
		},
	}

	if declaration.Since != nil {
		res.currentReport.Since = declaration.Since.AsTime()
	} else {
		res.currentReport.Since = time.Now()
	}

	return res
}

func buildBatch(reports []*api.ClimateReport) TimedValues {
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

func (l *climateLogger) PushReports(reports []*api.ClimateReport) {
	if len(reports) == 0 {
		return
	}
	l.mx.Lock()
	defer l.mx.Unlock()

	for _, s := range l.samplers {
		// we should rebuild the batch for each samplers as slices
		// gets modified by frequency cut-offs
		s.Add(buildBatch(reports))
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

func (l *climateLogger) PushTarget(target *api.ClimateTarget) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.currentReport.Current = target.Current
	l.currentReport.CurrentEnd = target.CurrentEnd
	if target.Next != nil && target.NextTime != nil {
		l.currentReport.Next = target.Next
		l.currentReport.NextTime = new(time.Time)
		*l.currentReport.NextTime = target.NextTime.AsTime()
		l.currentReport.NextEnd = target.NextEnd
	} else {
		l.currentReport.Next = nil
		l.currentReport.NextEnd = nil
		l.currentReport.NextTime = nil
	}
}

func (l *climateLogger) fromWindow(window string) api.ClimateTimeSeries {
	sampler, ok := l.samplersByWindow[window]
	if ok == false {
		sampler = l.samplers[0]
	}
	return sampler.TimeSeries()
}

func (l *climateLogger) GetClimateTimeSeries(window string) api.ClimateTimeSeries {
	l.mx.RLock()
	defer l.mx.RUnlock()
	// already a data copy, so it is safe
	return l.fromWindow(window)
}

func (l *climateLogger) GetClimateReport() *api.ZoneClimateReport {
	l.mx.RLock()
	defer l.mx.RUnlock()
	res := deepcopy.MustAnything(l.currentReport).(*api.ZoneClimateReport)
	if l.currentReport.NextTime != nil {
		res.NextTime = l.currentReport.NextTime
	}
	return res
}

func (l *climateLogger) Host() string {
	return l.host
}

func (l *climateLogger) ZoneName() string {
	return l.name
}

func (l *climateLogger) ZoneIdentifier() string {
	return ZoneIdentifier(l.host, l.name)
}
