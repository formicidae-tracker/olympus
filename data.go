package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/zeus"
)

type ClimateReportTimeSerie struct {
	NumAux         int
	Humidity       []lttb.Point
	TemperatureAnt []lttb.Point
	TemperatureAux [][]lttb.Point
}

type ClimateReportManager interface {
	Sample()
	Inbound() chan<- zeus.ClimateReport
	LastTenMinutes() ClimateReportTimeSerie
	LastHour() ClimateReportTimeSerie
	LastDay() ClimateReportTimeSerie
	LastWeek() ClimateReportTimeSerie
}

type window int

const (
	tenMinutes window = iota
	hour
	day
	week
	NbWindow
)

type request struct {
	w      window
	result chan<- ClimateReportTimeSerie
}

type climateReportManager struct {
	inbound      chan zeus.ClimateReport
	requests     chan request
	quit         chan struct{}
	wg           sync.WaitGroup
	numAux       int
	downsamplers [][]DataRollingSampler
	start        *time.Time
}

func (m *climateReportManager) addReportUnsafe(r *zeus.ClimateReport) {
	if m.start == nil {
		m.start = &time.Time{}
		*m.start = r.Time
	}
	if len(r.Temperatures) != 1+m.numAux {
		return
	}
	ellapsed := r.Time.Sub(*m.start)
	for i := 0; i < int(NbWindow); i++ {
		for _, downsamplers := range m.downsamplers {
			downsamplers[0].Add(ellapsed, float64(r.Humidity))
			for i, d := range downsamplers[1:] {
				d.Add(ellapsed, float64(r.Temperatures[i]))
			}
		}
	}
}

const (
	TenMinutesIdx = 0
	HourIdx       = 1
	DayIdx        = 2
	WeekIdx       = 3
)

func windowToIndex(w window) int {
	switch w {
	case tenMinutes:
		return TenMinutesIdx
	case hour:
		return HourIdx
	case day:
		return DayIdx
	case week:
		return WeekIdx
	default:
		return -1
	}
}

func (m *climateReportManager) reportSeries(w window) ClimateReportTimeSerie {
	idx := windowToIndex(w)
	if idx == -1 {
		return ClimateReportTimeSerie{}
	}
	d := m.downsamplers[idx]
	res := ClimateReportTimeSerie{
		NumAux:         m.numAux,
		Humidity:       d[0].TimeSerie(),
		TemperatureAnt: d[1].TimeSerie(),
		TemperatureAux: nil,
	}
	if m.numAux > 0 {
		for _, auxD := range d[2:] {
			res.TemperatureAux = append(res.TemperatureAux, auxD.TimeSerie())
		}
	}

	return res
}

func (m *climateReportManager) Sample() {
	m.quit = make(chan struct{})
	defer func() {
		close(m.quit)
		m.wg.Wait()
	}()
	for {
		select {
		case r := <-m.requests:
			r.result <- m.reportSeries(r.w)
		case r, ok := <-m.inbound:
			if ok == false {
				return
			}
			m.addReportUnsafe(&r)
		}
	}
}

func (m *climateReportManager) Inbound() chan<- zeus.ClimateReport {
	return m.inbound
}

func (m *climateReportManager) lastReport(w window) ClimateReportTimeSerie {
	m.wg.Add(1)
	res := make(chan ClimateReportTimeSerie)
	defer func() {
		close(res)
		m.wg.Done()
	}()
	go func() {
		m.requests <- request{w: w, result: res}
	}()
	select {
	case <-m.quit:
		return ClimateReportTimeSerie{}
	case r := <-res:
		return r
	}
}

func (m *climateReportManager) LastTenMinutes() ClimateReportTimeSerie {
	return m.lastReport(tenMinutes)
}

func (m *climateReportManager) LastHour() ClimateReportTimeSerie {
	return m.lastReport(hour)
}

func (m *climateReportManager) LastDay() ClimateReportTimeSerie {
	return m.lastReport(day)
}

func (m *climateReportManager) LastWeek() ClimateReportTimeSerie {
	return m.lastReport(week)
}

const (
	tenMinutesSamples = 500
	hourSamples       = 500
	daySamples        = 500
	weekSamples       = 500
)

func NewClimateReportManager(numAux int) ClimateReportManager {
	res := &climateReportManager{
		numAux:   numAux,
		inbound:  make(chan zeus.ClimateReport),
		requests: make(chan request),
	}
	cData := []struct {
		NbSample int
		Window   time.Duration
	}{
		{500, 10 * time.Minute},
		{500, 1 * time.Hour},
		{500, 24 * time.Hour},
		{500, 7 * 24 * time.Hour},
	}
	for _, d := range cData {
		var dsamplers []DataRollingSampler
		dsamplers = append(dsamplers, NewRollingSampler(d.Window, d.NbSample))
		dsamplers = append(dsamplers, NewRollingSampler(d.Window, d.NbSample))
		for i := 0; i < numAux; i++ {
			dsamplers = append(dsamplers, NewRollingSampler(d.Window, d.NbSample))
		}
		res.downsamplers = append(res.downsamplers, dsamplers)
	}
	return res
}

var stubClimateReporter ClimateReportManager

func setClimateReporterStub() {

	stubClimateReporter = NewClimateReportManager(3)
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	go stubClimateReporter.Sample()
	go func() {
		for t := start; t.Before(end); t = t.Add(500 * time.Millisecond) {
			ellapsed := t.Sub(start).Seconds()

			toAdd := zeus.ClimateReport{
				Time:     t,
				Humidity: zeus.Humidity(40.0 + 3*math.Cos(2*math.Pi/200.0*ellapsed) + 0.5*rand.NormFloat64()),
				Temperatures: []zeus.Temperature{
					zeus.Temperature(20.0 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					zeus.Temperature(20.5 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					zeus.Temperature(21.0 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					zeus.Temperature(21.5 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
				},
			}
			stubClimateReporter.Inbound() <- toAdd
		}
		toPrint := []interface{}{}
		for _, dsamplers := range stubClimateReporter.(*climateReportManager).downsamplers {
			for _, d := range dsamplers {
				toPrint = append(toPrint, len(d.(*rollingDownsampler).points))
			}
		}
	}()
}
