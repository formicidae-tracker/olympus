package main

import (
	"math/rand"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ZoneLoggerSuite struct {
	l ZoneLogger
}

var _ = Suite(&ZoneLoggerSuite{})

func (s *ZoneLoggerSuite) SetUpTest(c *C) {
	s.l = NewZoneLogger(&api.ZoneDeclaration{
		Host: "foo",
		Name: "bar",
	})
}

func (s *ZoneLoggerSuite) TearDownTest(c *C) {
}

func (s *ZoneLoggerSuite) TestLogsClimate(c *C) {
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

		for i, _ := range series.Humidity {
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

func (s *ZoneLoggerSuite) TestUninitialzedReport(c *C) {
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 0)
	c.Check(series.Temperature, HasLen, 0)
	report := s.l.GetClimateReport()
	c.Check(report.Humidity, IsNil)
	c.Check(report.Temperature, IsNil)
}

func (s *ZoneLoggerSuite) TestLogsHumididityOnlyClimate(c *C) {
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

func (s *ZoneLoggerSuite) TestLogsTemperatureOnlyClimate(c *C) {
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

func (s *ZoneLoggerSuite) TestLogsAlarms(c *C) {
	start := time.Now().Round(0)

	eventList := []*api.AlarmEvent{
		{
			Reason: newWithValue("foo"),
			Level:  newWithValue(api.AlarmLevel_WARNING),
		},
		{
			Reason: newWithValue("bar"),
			Level:  newWithValue(api.AlarmLevel_EMERGENCY),
		},
		{
			Reason: newWithValue("baz"),
			Level:  newWithValue(api.AlarmLevel_WARNING),
		},
	}

	lastState := map[string]struct {
		Time time.Time
		On   bool
	}{}

	events := make([]*api.AlarmEvent, 300)
	for i := 0; i < 300; i++ {
		r := rand.Intn(2000000)
		t := start.Add(time.Duration(r) * time.Millisecond)
		on := r%2 == 0
		event := deepcopy.MustAnything(eventList[i%3]).(*api.AlarmEvent)
		if on {
			event.Status = api.AlarmStatus_ON
		} else {
			event.Status = api.AlarmStatus_OFF
		}
		event.Time = timestamppb.New(t)
		ls := lastState[*event.Reason]
		if ls.Time.Before(t) {
			ls.Time = t
			ls.On = on
			lastState[*event.Reason] = ls
		}
		events[i] = event
	}
	s.l.PushAlarms(events)

	reports := s.l.GetAlarmReports()
	for _, r := range reports {
		switch r.Reason {
		case "foo":
			c.Check(api.AlarmLevel(r.Level), Equals, api.AlarmLevel_WARNING)
		case "bar":
			c.Check(api.AlarmLevel(r.Level), Equals, api.AlarmLevel_EMERGENCY)
		case "baz":
			c.Check(api.AlarmLevel(r.Level), Equals, api.AlarmLevel_WARNING)
		}
		c.Check(r.Events, HasLen, 100)
		for i, e := range r.Events {
			if i == 0 {
				continue
			}
			c.Check(r.Events[i-1].Time.AsTime().After(e.Time.AsTime()), Equals, false)
		}
	}

	var expectedWarning int32 = 0
	var expectedEmergency int32 = 0
	if lastState["bar"].On == true {
		expectedEmergency = 1
	}
	if lastState["foo"].On == true {
		expectedWarning += 1
	}
	if lastState["baz"].On == true {
		expectedWarning += 1
	}

	report := s.l.GetClimateReport()
	c.Check(report.ActiveEmergencies, Equals, expectedEmergency)
	c.Check(report.ActiveWarnings, Equals, expectedWarning)
}
