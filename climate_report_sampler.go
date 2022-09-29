package main

import (
	"sync"
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
)

type ClimateReportSampler interface {
	// Adds a batch of reports to the sampler.
	Add(reports []*proto.ClimateReport)
	LastTenMinutes() ClimateTimeSeries
	LastHour() ClimateTimeSeries
	LastDay() ClimateTimeSeries
	LastWeek() ClimateTimeSeries
}

type climateReportSampler struct {
	mx                                          sync.RWMutex
	lastTenMinutes, lastHour, lastDay, lastWeek ClimateDataDownsampler
}

func buildBatch(reports []*proto.ClimateReport) TimedValues {
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
			values[0] = appendValue(values[0], *r.Humidity, i+1)
		}
		for j, t := range r.Temperatures {
			values[j+1] = appendValue(values[j+1], t, i+1)
		}
	}
	return TimedValues{times: times, values: values}
}

func buildClimateTimeSeries(data [][]lttb.Point[float32]) ClimateTimeSeries {
	if len(data) == 0 {
		return ClimateTimeSeries{}
	}
	res := ClimateTimeSeries{
		Humidity: data[0],
	}
	if len(data) > 1 {
		res.TemperatureAnt = data[1]
	}
	if len(data) > 2 {
		res.TemperatureAux = data[2:]
	}
	return res
}

func (s *climateReportSampler) Add(reports []*proto.ClimateReport) {
	s.mx.Lock()
	defer s.mx.Unlock()

	values := buildBatch(reports)
	s.lastTenMinutes.Add(values)
	s.lastHour.Add(values)
	s.lastDay.Add(values)
	s.lastWeek.Add(values)
}

func (s *climateReportSampler) LastTenMinutes() ClimateTimeSeries {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.lastTenMinutes.TimeSeries()
}
func (s *climateReportSampler) LastHour() ClimateTimeSeries {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.lastHour.TimeSeries()
}
func (s *climateReportSampler) LastDay() ClimateTimeSeries {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.lastDay.TimeSeries()
}

func (s *climateReportSampler) LastWeek() ClimateTimeSeries {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.lastWeek.TimeSeries()
}

type climateReportSamplerArgs struct {
	tenMinuteSamples, hourSamples, daySamples, weekSamples int
}

func newClimateReportSampler(a climateReportSamplerArgs) ClimateReportSampler {
	return &climateReportSampler{
		lastTenMinutes: NewClimateDataDownsampler(10*time.Minute, a.tenMinuteSamples, 200*time.Millisecond),
		lastHour:       NewClimateDataDownsampler(1*time.Hour, a.hourSamples, time.Second),
		lastDay:        NewClimateDataDownsampler(24*time.Hour, a.daySamples, 30*time.Second),
		lastWeek:       NewClimateDataDownsampler(24*7*time.Hour, a.weekSamples, 3*time.Minute),
	}
}

func NewClimateReportSampler() ClimateReportSampler {
	return newClimateReportSampler(climateReportSamplerArgs{
		tenMinuteSamples: 500,
		hourSamples:      400,
		daySamples:       300,
		weekSamples:      300,
	})
}
