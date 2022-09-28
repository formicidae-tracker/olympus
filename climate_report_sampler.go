package main

import (
	"sync"
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
)

type ClimateReportSampler interface {
	// Adds a batch of reports to the sampler. If async is true, the
	// result may not be immediatly reported in the time series.
	Add(reports []*proto.ClimateReport, async bool)
	LastTenMinutes() ClimateTimeSerie
	LastHour() ClimateTimeSerie
	LastDay() ClimateTimeSerie
	LastWeek() ClimateTimeSerie
}

type climateReportSampler struct {
	mx                                          sync.RWMutex
	lastTenMinutes, lastHour, lastDay, lastWeek DataRollingSampler
}

func buildBatch(reports []*proto.ClimateReport) ([]time.Time, [][]*float32) {
	if len(reports) == 0 {
		return nil, nil
	}
	times := make([]time.Time, len(reports))
	values := make([][]*float32, len(reports))

	for i, r := range reports {
		times[i] = r.Time.AsTime()
		reportValues := make([]*float32, 1, len(r.Temperatures)+1)
		reportValues[0] = r.Humidity
		for _, t := range r.Temperatures {
			reportValues = append(reportValues, &t)
		}
		values[i] = reportValues
	}
	return times, values
}

func buildClimateTimeSeries(data [][]lttb.Point[float32]) ClimateTimeSerie {
	if len(data) == 0 {
		return ClimateTimeSerie{}
	}
	res := ClimateTimeSerie{
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

func (s *climateReportSampler) Add(reports []*proto.ClimateReport, async bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	var mx *sync.RWMutex
	if async == true {
		mx = &s.mx
	}

	times, values := buildBatch(reports)
	s.lastTenMinutes.AddBatch(times, values, mx)
	s.lastHour.AddBatch(times, values, mx)
	s.lastDay.AddBatch(times, values, mx)
	s.lastWeek.AddBatch(times, values, mx)
}

func (s *climateReportSampler) LastTenMinutes() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return buildClimateTimeSeries(s.lastTenMinutes.TimeSerie())
}
func (s *climateReportSampler) LastHour() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return buildClimateTimeSeries(s.lastHour.TimeSerie())
}
func (s *climateReportSampler) LastDay() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return buildClimateTimeSeries(s.lastDay.TimeSerie())
}

func (s *climateReportSampler) LastWeek() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return buildClimateTimeSeries(s.lastWeek.TimeSerie())
}

type climateReportSamplerArgs struct {
	tenMinuteSamples, hourSamples, daySamples, weekSamples int
}

func newClimateReportSampler(a climateReportSamplerArgs) ClimateReportSampler {
	return &climateReportSampler{
		lastTenMinutes: NewRollingSampler(10*time.Minute, a.tenMinuteSamples, 200*time.Millisecond),
		lastHour:       NewRollingSampler(1*time.Hour, a.hourSamples, time.Second),
		lastDay:        NewRollingSampler(24*time.Hour, a.daySamples, 30*time.Second),
		lastWeek:       NewRollingSampler(24*7*time.Hour, a.weekSamples, 3*time.Minute),
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
