package main

import (
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ClimateReportSamplerSuite struct {
}

var _ = Suite(&ClimateReportSamplerSuite{})

func (s *ClimateReportSamplerSuite) SetUpSuite(c *C) {
	c.Skip("for now")
}

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
		hasHumidity:     true,
		tenMinuteSample: 10,
		hourSample:      30,
		daySample:       60,
		weekSample:      200,
	}
	r := newClimateReportSampler(s)
	start := time.Now().Round(0)
	for ellapsed := time.Duration(0); ellapsed < 6*time.Minute; ellapsed += 500 * time.Millisecond {
		report := &proto.ClimateReport{
			Time:         timestamppb.New(start.Add(ellapsed)),
			Humidity:     newInitialized[float32](60.0),
			Temperatures: []float32{20.0, 21.0},
		}
		r.Add(report)
		checkClimateReport(c, r.LastTenMinutes(), 10*time.Minute, s.tenMinuteSample)
		checkClimateReport(c, r.LastHour(), 1*time.Hour, s.hourSample)
		checkClimateReport(c, r.LastDay(), 24*time.Hour, s.daySample)
		checkClimateReport(c, r.LastWeek(), 7*24*time.Hour, s.weekSample)
	}
}
