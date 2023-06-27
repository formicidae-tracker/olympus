package olympus

import (
	"math/rand"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ClimateLoggerSuite struct {
	l ClimateLogger
}

var _ = Suite(&ClimateLoggerSuite{})

func (s *ClimateLoggerSuite) SetUpTest(c *C) {
	s.l = NewClimateLogger(&api.ClimateDeclaration{
		Host: "foo",
		Name: "bar",
	})
}

func (s *ClimateLoggerSuite) TearDownTest(c *C) {
}

func (s *ClimateLoggerSuite) TestLogsClimate(c *C) {
	start := time.Now().Round(0)
	reports := make([]*api.ClimateReport, 60)
	for i := 0; i < len(reports); i++ {
		reports[i] = &api.ClimateReport{
			Time:         timestamppb.New(start.Add(time.Duration(i*500+rand.Intn(20)-10) * time.Millisecond)),
			Humidity:     newInitialized[float32](55.0),
			Temperatures: []float32{21.0},
		}
	}
	s.l.PushReports(reports)

	checkReport := func(c *C, series api.ClimateTimeSeries, size int) {
		if c.Check(len(series.Humidity), Equals, size) == false {
			return
		}
		if c.Check(len(series.Temperature), Equals, size) == false {
			return
		}

		for i := range series.Humidity {
			if i == 0 {
				continue
			}
			c.Check(series.Humidity[i-1].X <= series.Humidity[i].X,
				Equals,
				true, Commentf("at index %i", i))
			c.Check(series.Temperature[i-1].X <= series.Temperature[i].X,
				Equals,
				true, Commentf("at index %i", i))
		}
	}

	checkReport(c, s.l.GetClimateTimeSeries("10m"), 60)
	checkReport(c, s.l.GetClimateTimeSeries("1h"), 30)
	checkReport(c, s.l.GetClimateTimeSeries("1d"), 2)
	checkReport(c, s.l.GetClimateTimeSeries("1w"), 1)
}

func (s *ClimateLoggerSuite) TestUninitialzedReport(c *C) {
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 0)
	c.Check(series.Temperature, HasLen, 0)
	report := s.l.GetClimateReport()
	c.Check(report.Humidity, IsNil)
	c.Check(report.Temperature, IsNil)
}

func (s *ClimateLoggerSuite) TestLogsHumididityOnlyClimate(c *C) {
	s.l.PushReports([]*api.ClimateReport{
		{
			Time:         timestamppb.New(time.Now()),
			Humidity:     newInitialized[float32](55.0),
			Temperatures: []float32{},
		},
	})
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 1)
	c.Check(series.Temperature, HasLen, 0)
	report := s.l.GetClimateReport()
	c.Check(*report.Humidity, Equals, float32(55.0))
	c.Check(report.Temperature, IsNil)
}

func (s *ClimateLoggerSuite) TestLogsTemperatureOnlyClimate(c *C) {
	s.l.PushReports([]*api.ClimateReport{
		{
			Time:         timestamppb.New(time.Now()),
			Humidity:     nil,
			Temperatures: []float32{20.0},
		},
	})
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 0)
	c.Check(series.Temperature, HasLen, 1)
	report := s.l.GetClimateReport()
	c.Check(report.Humidity, IsNil)
	c.Check(*report.Temperature, Equals, float32(20.0))
}

func newWithValue[T any](v T) *T {
	res := new(T)
	*res = v
	return res
}
