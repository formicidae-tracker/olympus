package main

import (
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
)

type ClimateReportSampler interface {
	Add(*proto.ClimateReport)
	LastTenMinutes() ClimateTimeSerie
	LastHour() ClimateTimeSerie
	LastDay() ClimateTimeSerie
	LastWeek() ClimateTimeSerie
}

type climateData struct {
	time           []time.Time
	humidity       []float32
	temperatureAnt []float32
	temperatureAux [][]float32
}

func (d *climateData) push(report *proto.ClimateReport, window time.Duration) {
	t := report.Time.AsTime()
	index := BackLinearSearch(d.time, t,
		func(a, b time.Time) bool { return a.Before(b) })

	Insert(d.time, t, index)
	if report.Humidity != nil {
		Insert(d.humidity, *report.Humidity, index)
	}
	if len(report.Temperatures) > 0 {
		Insert(d.temperatureAnt, report.Temperatures[0], index)
	}
	if len(report.Temperatures) > 1 && d.temperatureAux == nil {
		d.temperatureAux = make([][]float32, len(report.Temperatures)-1)
	}
	for i := 1; i < len(report.Temperatures); i++ {
		Insert(d.temperatureAux[i-1], report.Temperatures[i], index)
	}
	d.rollOutOfWindow(window)
}

func (d *climateData) rollOutOfWindow(window time.Duration) {
	last := d.time[len(d.time)-1]
	minValue := last.Add(-window)
	newStart := LinearSearch(d.time, minValue,
		func(a, b time.Time) bool { return a.Before(b) })

	d.time = d.time[newStart:]
	if d.humidity != nil {
		d.humidity = d.humidity[newStart:]
	}
	if d.temperatureAnt != nil {
		d.temperatureAnt = d.temperatureAnt[newStart:]
	}
	for i, t := range d.temperatureAux {
		d.temperatureAux[i] = t[newStart:]
	}
}

func (d *climateData) timeVector() []float32 {
	if len(d.time) == 0 {
		return nil
	}
	lastTime := d.time[len(d.time)-1]
	res := make([]float32, len(d.time))
	for i, t := range d.time {
		res[i] = float32(t.Sub(lastTime).Seconds())
	}
	return res
}

func sample(times []float32, values []float32, samples int) []lttb.Point {
	if len(values) == 0 {
		return nil
	}
	points := make([]lttb.Point, len(values))
	for i, t := range times {
		points[i] = lttb.Point{X: float64(t), Y: float64(values[i])}
	}
	if len(points) <= samples {
		return points
	}
	return lttb.LTTB(points, samples)
}

func (d *climateData) computeCache(samples int) ClimateTimeSerie {
	times := d.timeVector()
	if len(times) == 0 {
		return ClimateTimeSerie{}
	}
	humidity := sample(times, d.humidity, samples)
	temperatureAnt := sample(times, d.temperatureAnt, samples)
	var temperatureAux [][]lttb.Point = nil
	if len(d.temperatureAux) > 0 {
		temperatureAux = make([][]lttb.Point, len(d.temperatureAux))
		for i, aux := range d.temperatureAux {
			temperatureAux[i] = sample(times, aux, samples)
		}
	}
	return ClimateTimeSerie{
		Humidity:       humidity,
		TemperatureAnt: temperatureAnt,
		TemperatureAux: temperatureAux,
	}
}

type cachedSeries struct {
	reports climateData
	series  ClimateTimeSerie
	window  time.Duration
	samples int
}

func (s *cachedSeries) push(report *proto.ClimateReport) {
	s.reports.push(report, s.window)
}

type climateReportSampler struct {
	mx                                          sync.RWMutex
	lastTenMinutes, lastHour, lastDay, lastWeek cachedSeries
	caching                                     bool
}

func (s *climateReportSampler) Add(report *proto.ClimateReport) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.push(report)

	if s.caching == false {
		s.caching = true
		go s.cacheSeries()
	}
}

func (s *climateReportSampler) push(report *proto.ClimateReport) {
	s.lastTenMinutes.push(report)
	s.lastHour.push(report)
	s.lastDay.push(report)
	s.lastWeek.push(report)
}

func (s *climateReportSampler) cacheSeries() {
	wg := sync.WaitGroup{}
	computeCache := func(c *cachedSeries) {
		s.mx.RLock()
		data := deepcopy.MustAnything(c.reports).(climateData)
		s.mx.RUnlock()

		res := data.computeCache(c.samples)

		s.mx.Lock()
		defer s.mx.Unlock()
		c.series = res
		wg.Done()
	}
	wg.Add(4)
	go computeCache(&s.lastTenMinutes)
	go computeCache(&s.lastHour)
	go computeCache(&s.lastDay)
	go computeCache(&s.lastWeek)

	wg.Wait()

	s.mx.Lock()
	defer s.mx.Unlock()
	s.caching = false
}

func (s *climateReportSampler) LastTenMinutes() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.lastTenMinutes.series
}
func (s *climateReportSampler) LastHour() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.lastHour.series
}
func (s *climateReportSampler) LastDay() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.lastDay.series
}

func (s *climateReportSampler) LastWeek() ClimateTimeSerie {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.lastWeek.series
}

type climateReportSamplerSetting struct {
	numAux, tenMinuteSample, hourSample, daySample, weekSample int
	hasHumidity                                                bool
}

func newClimateReportSampler(s climateReportSamplerSetting) ClimateReportSampler {
	if s.numAux < 0 {
		s.numAux = 0
	}
	return nil

}

func NewClimateReportSampler(numberOfAux int) ClimateReportSampler {
	return newClimateReportSampler(climateReportSamplerSetting{
		numAux:          numberOfAux,
		tenMinuteSample: 500,
		hourSample:      400,
		daySample:       300,
		weekSample:      300,
	})
}
