package main

import (
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ClimateReportSamplerSuite struct {
}

var _ = Suite(&ClimateReportSamplerSuite{})

func checkSeries(c *C, series []lttb.Point[float32], window time.Duration, samples int) {
	c.Assert(len(series) <= samples, Equals, true)
	c.Assert(len(series) > 0, Equals, true)
	actualWindow := float64(series[len(series)-1].X - series[0].X)
	c.Check(actualWindow <= window.Seconds(), Equals, true, Commentf("window:%s obtained:%f", window, actualWindow))
}

func checkClimateReport(c *C, report ClimateTimeSeries, window time.Duration, samples int) {
	checkSeries(c, report.Humidity, window, samples)
	checkSeries(c, report.TemperatureAnt, window, samples)
	for _, serie := range report.TemperatureAux {
		checkSeries(c, serie, window, samples)
	}
}

func (su *ClimateReportSamplerSuite) TestClimateReportSuite(c *C) {
	a := climateReportSamplerArgs{
		tenMinuteSamples: 10,
		hourSamples:      30,
		daySamples:       60,
		weekSamples:      200,
	}
	r := newClimateReportSampler(a)

	start := time.Now().Round(0)
	for ellapsed := time.Duration(0); ellapsed < 6*time.Minute; ellapsed += 500 * time.Millisecond {
		report := &proto.ClimateReport{
			Time:         timestamppb.New(start.Add(ellapsed)),
			Humidity:     newInitialized[float32](60.0),
			Temperatures: []float32{20.0, 21.0},
		}
		r.Add([]*proto.ClimateReport{report})
		checkClimateReport(c, r.LastTenMinutes(), 10*time.Minute, a.tenMinuteSamples)
		checkClimateReport(c, r.LastHour(), 1*time.Hour, a.hourSamples)
		checkClimateReport(c, r.LastDay(), 24*time.Hour, a.daySamples)
		checkClimateReport(c, r.LastWeek(), 7*24*time.Hour, a.weekSamples)
	}
}
