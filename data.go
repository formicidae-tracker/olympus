package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/formicidae-tracker/dieu"
	"github.com/dgryski/go-lttb"
)

type ClimateReportTimeSerie struct {
	Humidity        []lttb.Point
	TemperatureAnt  []lttb.Point
	TemperatureAux1 []lttb.Point
	TemperatureAux2 []lttb.Point
	TemperatureAux3 []lttb.Point
}

type ClimateReportManager interface {
	Sample()
	Inbound() chan<- dieu.ClimateReport
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
	inbound      chan dieu.ClimateReport
	requests     chan request
	quit         chan struct{}
	wg           sync.WaitGroup
	downsamplers []rollingDownsampler
	start        *time.Time
}

type rollingDownsampler struct {
	xPeriod   float64
	threshold int
	sampled   bool
	points    []lttb.Point
}

func newRollingDownsampler(period float64, nbSamples int) rollingDownsampler {
	res := rollingDownsampler{
		xPeriod:   period,
		threshold: nbSamples,
		points:    make([]lttb.Point, 0, nbSamples),
	}
	return res
}

func (d *rollingDownsampler) add(p lttb.Point) {
	d.points = append(d.points, p)
	idx := 0
	last := d.points[len(d.points)-1].X
	for {
		if (last - d.points[idx].X) <= d.xPeriod {
			break
		}
		idx += 1
	}

	if idx != 0 {
		d.points = d.points[idx:]
	}
}

func (d *rollingDownsampler) getPoints() []lttb.Point {
	return lttb.LTTB(d.points, d.threshold)
}

func (m *climateReportManager) addReportUnsafe(r *dieu.ClimateReport) {
	if m.start == nil {
		m.start = &time.Time{}
		*m.start = r.Time
	}
	ellapsed := r.Time.Sub(*m.start).Seconds()
	for i := 0; i < int(NbWindow); i++ {
		m.downsamplers[5*i].add(lttb.Point{X: ellapsed, Y: float64(r.Humidity)})
		for j := 0; j < 4; j++ {
			m.downsamplers[5*i+j+1].add(lttb.Point{X: ellapsed, Y: float64(r.Temperatures[j])})
		}
	}
}

const (
	humidityTenMinutesIdx = iota
	temperatureAntTenMinutesIdx
	temperatureAux1TenMinutesIdx
	temperatureAux2TenMinutesIdx
	temperatureAux3TenMinutesIdx
	humidityHourIdx
	temperatureAntHourIdx
	temperatureAux1HourIdx
	temperatureAux2HourIdx
	temperatureAux3HourIdx
	humidityDayIdx
	temperatureAntDayIdx
	temperatureAux1DayIdx
	temperatureAux2DayIdx
	temperatureAux3DayIdx
	humidityWeekIdx
	temperatureAntWeekIdx
	temperatureAux1WeekIdx
	temperatureAux2WeekIdx
	temperatureAux3WeekIdx
)

func (m *climateReportManager) Sample() {
	m.quit = make(chan struct{})
	defer func() {
		close(m.quit)
		m.wg.Wait()
	}()
	for {
		select {
		case r := <-m.requests:
			switch r.w {
			case tenMinutes:
				r.result <- ClimateReportTimeSerie{
					Humidity:        m.downsamplers[humidityTenMinutesIdx].getPoints(),
					TemperatureAnt:  m.downsamplers[temperatureAntTenMinutesIdx].getPoints(),
					TemperatureAux1: m.downsamplers[temperatureAux1TenMinutesIdx].getPoints(),
					TemperatureAux2: m.downsamplers[temperatureAux2TenMinutesIdx].getPoints(),
					TemperatureAux3: m.downsamplers[temperatureAux3TenMinutesIdx].getPoints(),
				}
			case hour:
				r.result <- ClimateReportTimeSerie{
					Humidity:        m.downsamplers[humidityHourIdx].getPoints(),
					TemperatureAnt:  m.downsamplers[temperatureAntHourIdx].getPoints(),
					TemperatureAux1: m.downsamplers[temperatureAux1HourIdx].getPoints(),
					TemperatureAux2: m.downsamplers[temperatureAux2HourIdx].getPoints(),
					TemperatureAux3: m.downsamplers[temperatureAux3HourIdx].getPoints(),
				}
			case day:
				r.result <- ClimateReportTimeSerie{
					Humidity:        m.downsamplers[humidityDayIdx].getPoints(),
					TemperatureAnt:  m.downsamplers[temperatureAntDayIdx].getPoints(),
					TemperatureAux1: m.downsamplers[temperatureAux1DayIdx].getPoints(),
					TemperatureAux2: m.downsamplers[temperatureAux2DayIdx].getPoints(),
					TemperatureAux3: m.downsamplers[temperatureAux3DayIdx].getPoints(),
				}
			case week:
				r.result <- ClimateReportTimeSerie{
					Humidity:        m.downsamplers[humidityWeekIdx].getPoints(),
					TemperatureAnt:  m.downsamplers[temperatureAntWeekIdx].getPoints(),
					TemperatureAux1: m.downsamplers[temperatureAux1WeekIdx].getPoints(),
					TemperatureAux2: m.downsamplers[temperatureAux2WeekIdx].getPoints(),
					TemperatureAux3: m.downsamplers[temperatureAux3WeekIdx].getPoints(),
				}
			default:
				r.result <- ClimateReportTimeSerie{}
			}
		case r, ok := <-m.inbound:
			if ok == false {
				return
			}
			m.addReportUnsafe(&r)
		}
	}
}

func (m *climateReportManager) Inbound() chan<- dieu.ClimateReport {
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

func NewClimateReportManager() ClimateReportManager {
	return &climateReportManager{
		inbound:  make(chan dieu.ClimateReport),
		requests: make(chan request),
		downsamplers: []rollingDownsampler{
			newRollingDownsampler(10*time.Minute.Seconds(), tenMinutesSamples),
			newRollingDownsampler(10*time.Minute.Seconds(), tenMinutesSamples),
			newRollingDownsampler(10*time.Minute.Seconds(), tenMinutesSamples),
			newRollingDownsampler(10*time.Minute.Seconds(), tenMinutesSamples),
			newRollingDownsampler(10*time.Minute.Seconds(), tenMinutesSamples),
			newRollingDownsampler(time.Hour.Seconds(), hourSamples),
			newRollingDownsampler(time.Hour.Seconds(), hourSamples),
			newRollingDownsampler(time.Hour.Seconds(), hourSamples),
			newRollingDownsampler(time.Hour.Seconds(), hourSamples),
			newRollingDownsampler(time.Hour.Seconds(), hourSamples),
			newRollingDownsampler(24*time.Hour.Seconds(), daySamples),
			newRollingDownsampler(24*time.Hour.Seconds(), daySamples),
			newRollingDownsampler(24*time.Hour.Seconds(), daySamples),
			newRollingDownsampler(24*time.Hour.Seconds(), daySamples),
			newRollingDownsampler(24*time.Hour.Seconds(), daySamples),
			newRollingDownsampler(7*24*time.Hour.Seconds(), weekSamples),
			newRollingDownsampler(7*24*time.Hour.Seconds(), weekSamples),
			newRollingDownsampler(7*24*time.Hour.Seconds(), weekSamples),
			newRollingDownsampler(7*24*time.Hour.Seconds(), weekSamples),
			newRollingDownsampler(7*24*time.Hour.Seconds(), weekSamples),
		},
	}
}

var stubClimateReporter ClimateReportManager

func setClimateReporterStub() {

	stubClimateReporter = NewClimateReportManager()
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	go stubClimateReporter.Sample()
	go func() {
		for t := start; t.Before(end); t = t.Add(500 * time.Millisecond) {
			ellapsed := t.Sub(start).Seconds()

			toAdd := dieu.ClimateReport{
				Time:     t,
				Humidity: dieu.Humidity(40.0 + 3*math.Cos(2*math.Pi/200.0*ellapsed) + 0.5*rand.NormFloat64()),
				Temperatures: [4]dieu.Temperature{
					dieu.Temperature(20.0 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					dieu.Temperature(20.5 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					dieu.Temperature(21.0 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
					dieu.Temperature(21.5 + 0.5*math.Cos(2*math.Pi/1800.0*ellapsed) + 0.1*rand.NormFloat64()),
				},
			}
			stubClimateReporter.Inbound() <- toAdd
		}
		toPrint := []interface{}{}
		for _, d := range stubClimateReporter.(*climateReportManager).downsamplers {
			toPrint = append(toPrint, len(d.points))
		}
	}()
}
