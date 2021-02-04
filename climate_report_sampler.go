package main

import (
	"fmt"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ClimateReportSampler interface {
	Add(zeus.ClimateReport) error
	LastTenMinutes() ClimateTimeSerie
	LastHour() ClimateTimeSerie
	LastDay() ClimateTimeSerie
	LastWeek() ClimateTimeSerie
}

type climateReportSampler struct {
	start             *time.Time
	tenMinuteSamplers []DataRollingSampler
	hourSamplers      []DataRollingSampler
	daySamplers       []DataRollingSampler
	weekSamplers      []DataRollingSampler
	allSamplers       [][]DataRollingSampler
}

func (s *climateReportSampler) Add(r zeus.ClimateReport) error {
	if len(r.Temperatures)+1 != len(s.tenMinuteSamplers) {
		return fmt.Errorf("invalid number of value in report, got %d, expected %d", len(r.Temperatures)+1, len(s.tenMinuteSamplers))
	}
	if s.start == nil {
		s.start = &r.Time
	}
	ellapsed := r.Time.Sub(*s.start)
	for _, windowSamplers := range s.allSamplers {
		if zeus.IsUndefined(r.Humidity) == false {
			windowSamplers[0].Add(ellapsed, float64(r.Humidity))
		}
		for i, temperatureSampler := range windowSamplers[1:] {
			t := r.Temperatures[i]
			if zeus.IsUndefined(t) {
				continue
			}
			temperatureSampler.Add(ellapsed, t.Value())
		}
	}
	return nil
}

func (s *climateReportSampler) timeSerieUnsafe(samplers []DataRollingSampler) ClimateTimeSerie {
	res := ClimateTimeSerie{
		Humidity:       samplers[0].TimeSerie(),
		TemperatureAnt: samplers[1].TimeSerie(),
	}

	if len(samplers) == 2 {
		return res
	}
	for _, s := range samplers[2:] {
		res.TemperatureAux = append(res.TemperatureAux, s.TimeSerie())
	}
	return res
}

func (s *climateReportSampler) LastTenMinutes() ClimateTimeSerie {
	return s.timeSerieUnsafe(s.tenMinuteSamplers)
}
func (s *climateReportSampler) LastHour() ClimateTimeSerie {
	return s.timeSerieUnsafe(s.hourSamplers)
}
func (s *climateReportSampler) LastDay() ClimateTimeSerie {
	return s.timeSerieUnsafe(s.daySamplers)
}
func (s *climateReportSampler) LastWeek() ClimateTimeSerie {
	return s.timeSerieUnsafe(s.weekSamplers)
}

type climateReportSamplerSetting struct {
	tenMinute, hour, day, week                                 time.Duration
	numAux, tenMinuteSample, hourSample, daySample, weekSample int
}

func newClimateReportSampler(s climateReportSamplerSetting) ClimateReportSampler {
	if s.numAux < 0 {
		s.numAux = 0
	}
	if s.numAux > 3 {
		s.numAux = 3
	}
	res := &climateReportSampler{}
	for i := 0; i < s.numAux+2; i++ {
		res.tenMinuteSamplers = append(res.tenMinuteSamplers, NewRollingSampler(s.tenMinute, s.tenMinuteSample))
		res.hourSamplers = append(res.hourSamplers, NewRollingSampler(s.hour, s.hourSample))
		res.daySamplers = append(res.daySamplers, NewRollingSampler(s.day, s.daySample))
		res.weekSamplers = append(res.weekSamplers, NewRollingSampler(s.week, s.weekSample))
	}
	res.allSamplers = [][]DataRollingSampler{
		res.tenMinuteSamplers,
		res.hourSamplers,
		res.daySamplers,
		res.weekSamplers,
	}
	return res

}

func NewClimateReportSampler(numberOfAux int) ClimateReportSampler {
	return newClimateReportSampler(climateReportSamplerSetting{
		numAux:          numberOfAux,
		tenMinute:       10 * time.Minute,
		tenMinuteSample: 500,
		hour:            time.Hour,
		hourSample:      400,
		day:             24 * time.Hour,
		daySample:       300,
		week:            24 * 7 * time.Hour,
		weekSample:      300,
	})
}
