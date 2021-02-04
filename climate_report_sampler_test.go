package main

import (
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/zeus"
	. "gopkg.in/check.v1"
)

type ClimateReportSamplerSuite struct {
}

var _ = Suite(&ClimateReportSamplerSuite{})

func checkSeries(c *C, series []lttb.Point, window time.Duration, samples int) {
	c.Assert(len(series) <= samples, Equals, true)
	c.Assert(len(series) > 0, Equals, true)
	actualWindow := series[len(series)-1].X - series[0].X
	c.Check(actualWindow <= window.Seconds(), Equals, true, Commentf("window:%s obtained:%f", window, actualWindow))
}

func checkClimateReport(c *C, report ClimateTimeSerie, window time.Duration, samples int) {
	checkSeries(c, report.Humidity, window, samples)
	checkSeries(c, report.TemperatureAnt, window, samples)
	for _, serie := range report.TemperatureAux {
		checkSeries(c, serie, window, samples)
	}
}

func (su *ClimateReportSamplerSuite) TestClimateReportSuite(c *C) {
	s := climateReportSamplerSetting{
		numAux:          1,
		tenMinute:       10 * time.Second,
		tenMinuteSample: 10,
		hour:            30 * time.Second,
		hourSample:      30,
		day:             1 * time.Minute,
		daySample:       60,
		week:            2 * time.Minute,
		weekSample:      200,
	}
	r := newClimateReportSampler(s)
	start := time.Now().Round(0)
	for ellapsed := time.Duration(0); ellapsed < 6*time.Minute; ellapsed += 500 * time.Millisecond {
		cr := zeus.ClimateReport{
			Time:         start.Add(ellapsed),
			Humidity:     60.0,
			Temperatures: []zeus.Temperature{20.0, 21.0},
		}
		c.Assert(r.Add(cr), IsNil)
		checkClimateReport(c, r.LastTenMinutes(), s.tenMinute, s.tenMinuteSample)
		checkClimateReport(c, r.LastHour(), s.hour, s.hourSample)
		checkClimateReport(c, r.LastDay(), s.day, s.daySample)
		checkClimateReport(c, r.LastWeek(), s.week, s.weekSample)
	}
}
