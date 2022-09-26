package main

import (
	"fmt"
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
)

type ClimateReportSampler interface {
	Add(proto.ClimateReport) error
	LastTenMinutes() ClimateTimeSerie
	LastHour() ClimateTimeSerie
	LastDay() ClimateTimeSerie
	LastWeek() ClimateTimeSerie
}

type samplers struct {
	temperatures []DataRollingSampler
	humidity     DataRollingSampler
}

func (s *samplers) add(d time.Duration, temperatures []float32, humidity *float32) error {
	if len(temperatures) != len(s.temperatures) {
		return fmt.Errorf("invalid number of temperatures in report, got %d, expected %d",
			len(temperatures),
			len(s.temperatures),
		)
	}
	if s.humidity != nil && humidity == nil {
		return fmt.Errorf("expected a report with humidity, but none passed")
	}

	if s.humidity == nil && humidity != nil {
		return fmt.Errorf("expected a report without humidity, but got one")
	}

	if humidity != nil {
		s.humidity.Add(d, float64(*humidity))
	}

	for i, t := range temperatures {
		s.temperatures[i].Add(d, float64(t))
	}
	return nil
}

func (s *samplers) getTimeSeries() ClimateTimeSerie {
	res := ClimateTimeSerie{Humidity: nil, TemperatureAnt: nil, TemperatureAux: nil}

	if s.humidity != nil {
		res.Humidity = s.humidity.TimeSerie()
	}

	if len(s.temperatures) > 0 {
		res.TemperatureAnt = s.temperatures[0].TimeSerie()
	}

	if len(s.temperatures) > 1 {
		res.TemperatureAux = make([][]lttb.Point, 0, len(s.temperatures)-1)
		for _, t := range s.temperatures {
			res.TemperatureAux = append(res.TemperatureAux, t.TimeSerie())
		}
	}
	return res
}

func newSamplers(numTemperatures int, hasHumidity bool, window time.Duration, nbSamples int) samplers {
	res := samplers{temperatures: make([]DataRollingSampler, 0, numTemperatures)}
	for i := 0; i < numTemperatures; i++ {
		res.temperatures = append(res.temperatures, NewRollingSampler(window, nbSamples))
	}
	if hasHumidity == true {
		res.humidity = NewRollingSampler(window, nbSamples)
	}
	return res
}

type climateReportSampler struct {
	start             *time.Time
	tenMinuteSamplers samplers
	hourSamplers      samplers
	daySamplers       samplers
	weekSamplers      samplers
}

func (s *climateReportSampler) Add(report proto.ClimateReport) error {
	t := report.Time.AsTime()
	if s.start == nil {
		s.start = &t
	}
	ellapsed := t.Sub(*s.start)
	allSamplers := []*samplers{&s.tenMinuteSamplers, &s.hourSamplers, &s.daySamplers, &s.weekSamplers}

	for _, sampler := range allSamplers {
		err := sampler.add(ellapsed, report.Temperatures, report.Humidity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *climateReportSampler) LastTenMinutes() ClimateTimeSerie {
	return s.tenMinuteSamplers.getTimeSeries()
}
func (s *climateReportSampler) LastHour() ClimateTimeSerie {
	return s.hourSamplers.getTimeSeries()
}
func (s *climateReportSampler) LastDay() ClimateTimeSerie {
	return s.daySamplers.getTimeSeries()
}
func (s *climateReportSampler) LastWeek() ClimateTimeSerie {
	return s.weekSamplers.getTimeSeries()
}

type climateReportSamplerSetting struct {
	numAux, tenMinuteSample, hourSample, daySample, weekSample int
	hasHumidity                                                bool
}

func newClimateReportSampler(s climateReportSamplerSetting) ClimateReportSampler {
	if s.numAux < 0 {
		s.numAux = 0
	}
	if s.numAux > 3 {
		s.numAux = 3
	}
	res := &climateReportSampler{
		tenMinuteSamplers: newSamplers(s.numAux+1, s.hasHumidity, 10*time.Minute, s.tenMinuteSample),
		hourSamplers:      newSamplers(s.numAux+1, s.hasHumidity, 1*time.Hour, s.hourSample),
		daySamplers:       newSamplers(s.numAux+1, s.hasHumidity, 24*time.Hour, s.daySample),
		weekSamplers:      newSamplers(s.numAux+1, s.hasHumidity, 7*24*time.Hour, s.weekSample),
	}
	return res

}

func NewClimateReportSampler(numberOfAux int, hasHumidity bool) ClimateReportSampler {
	return newClimateReportSampler(climateReportSamplerSetting{
		numAux:          numberOfAux,
		hasHumidity:     hasHumidity,
		tenMinuteSample: 500,
		hourSample:      400,
		daySample:       300,
		weekSample:      300,
	})
}
