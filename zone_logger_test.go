package main

import (
	"math/rand"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ZoneLoggerSuite struct {
	l ZoneLogger
}

var _ = Suite(&ZoneLoggerSuite{})

func (s *ZoneLoggerSuite) SetUpTest(c *C) {
	s.l = NewZoneLogger(&proto.ZoneDeclaration{
		Host:        "foo",
		Name:        "bar",
		NumAux:      0,
		HasHumidity: true,
	})
}

func (s *ZoneLoggerSuite) TearDownTest(c *C) {
}

func (s *ZoneLoggerSuite) TestLogsClimate(c *C) {
	start := time.Now().Round(0)
	reports := make([]*proto.ClimateReport, 20)
	for i := 0; i < 20; i++ {
		reports[i] = &proto.ClimateReport{
			Time:         timestamppb.New(start.Add(time.Duration(rand.Intn(20000)) * time.Millisecond)),
			Humidity:     newInitialized[float32](20.0),
			Temperatures: []float32{20.0},
		}
	}
	s.l.PushReports(reports)

	checkReport := func(c *C, series ClimateTimeSerie) {
		if c.Check(len(series.Humidity), Equals, 20) == false {
			return
		}
		if c.Check(len(series.TemperatureAnt), Equals, 20) == false {
			return
		}

		for i, _ := range series.Humidity {
			if i == 0 {
				continue
			}
			c.Check(series.Humidity[i-1].X <= series.Humidity[i].X,
				Equals,
				true, Commentf("at index %i", i))
			c.Check(series.TemperatureAnt[i-1].X <= series.TemperatureAnt[i].X,
				Equals,
				true, Commentf("at index %i", i))
		}
	}

	checkReport(c, s.l.GetClimateTimeSeries("10m"))
	checkReport(c, s.l.GetClimateTimeSeries("1h"))
	checkReport(c, s.l.GetClimateTimeSeries("1d"))
	checkReport(c, s.l.GetClimateTimeSeries("1w"))
}

func (s *ZoneLoggerSuite) TestUninitialzedReport(c *C) {
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 0)
	c.Check(series.TemperatureAnt, HasLen, 0)
	report := s.l.GetClimateReport()
	c.Check(report.Humidity, IsNil)
	c.Check(report.Temperature, IsNil)
}

func (s *ZoneLoggerSuite) TestLogsHumididityOnlyClimate(c *C) {
	s.l.PushReports([]*proto.ClimateReport{
		{
			Time:         timestamppb.New(time.Now()),
			Humidity:     newInitialized[float32](55.0),
			Temperatures: []float32{},
		},
	})
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 1)
	c.Check(series.TemperatureAnt, HasLen, 0)
	report := s.l.GetClimateReport()
	c.Check(*report.Humidity, Equals, float32(55.0))
	c.Check(report.Temperature, IsNil)
}

func (s *ZoneLoggerSuite) TestLogsTemperatureOnlyClimate(c *C) {
	s.l.PushReports([]*proto.ClimateReport{
		{
			Time:         timestamppb.New(time.Now()),
			Humidity:     nil,
			Temperatures: []float32{20.0},
		},
	})
	series := s.l.GetClimateTimeSeries("")
	c.Check(series.Humidity, HasLen, 0)
	c.Check(series.TemperatureAnt, HasLen, 1)
	report := s.l.GetClimateReport()
	c.Check(report.Humidity, IsNil)
	c.Check(*report.Temperature, Equals, float32(20.0))
}

func (s *ZoneLoggerSuite) TestLogsAlarms(c *C) {
	start := time.Now().Round(0)

	eventList := []*proto.AlarmEvent{
		{
			Reason: "foo",
			Level:  1,
		},
		{
			Reason: "bar",
			Level:  2,
		},
		{
			Reason: "baz",
			Level:  1,
		},
	}

	lastState := map[string]struct {
		Time time.Time
		On   bool
	}{}

	events := make([]*proto.AlarmEvent, 300)
	for i := 0; i < 300; i++ {
		r := rand.Intn(2000000)
		t := start.Add(time.Duration(r) * time.Millisecond)
		on := r%2 == 0
		event := deepcopy.MustAnything(eventList[i%3]).(*proto.AlarmEvent)
		if on {
			event.Status = proto.AlarmStatus_ALARM_ON
		} else {
			event.Status = proto.AlarmStatus_ALARM_OFF
		}
		event.Time = timestamppb.New(t)
		ls := lastState[event.Reason]
		if ls.Time.Before(t) {
			ls.Time = t
			ls.On = on
			lastState[event.Reason] = ls
		}
		events[i] = event
	}
	s.l.PushAlarms(events)

	reports := s.l.GetAlarmReports()
	for _, r := range reports {
		switch r.Reason {
		case "foo":
			c.Check(r.Level, Equals, 1)
		case "bar":
			c.Check(r.Level, Equals, 2)
		case "baz":
			c.Check(r.Level, Equals, 1)
		}
		c.Check(r.Events, HasLen, 100)
		for i, e := range r.Events {
			if i == 0 {
				continue
			}
			c.Check(r.Events[i-1].Time.After(e.Time), Equals, false)
		}
	}

	expectedWarning := 0
	expectedEmergency := 0
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
