package main

import (
	"math/rand"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type AlarmLoggerSuite struct {
	l AlarmLogger
}

var _ = Suite(&AlarmLoggerSuite{})

func (s *AlarmLoggerSuite) SetUpTest(c *C) {
	s.l = NewAlarmLogger()
}

func (s *AlarmLoggerSuite) TearDownTest(c *C) {
}

func (s *AlarmLoggerSuite) TestLogsAlarms(c *C) {
	start := time.Now().Round(0)

	eventList := []*api.AlarmEvent{
		{
			Identification: "foo",
			Level:          api.AlarmLevel_WARNING,
		},
		{
			Identification: "bar",
			Level:          api.AlarmLevel_EMERGENCY,
		},
		{
			Identification: "baz",
			Level:          api.AlarmLevel_WARNING,
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
		ls := lastState[event.Identification]
		if ls.Time.Before(t) {
			ls.Time = t
			ls.On = on
			lastState[event.Identification] = ls
		}
		events[i] = event
	}
	s.l.PushAlarms(events)

	reports := s.l.GetReports()
	for _, r := range reports {
		switch r.Identification {
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
			c.Check(r.Events[i-1].Time.After(e.Time), Equals, false)
		}
	}

	var expectedWarning int = 0
	var expectedEmergency int = 0
	if lastState["bar"].On == true {
		expectedEmergency = 1
	}
	if lastState["foo"].On == true {
		expectedWarning += 1
	}
	if lastState["baz"].On == true {
		expectedWarning += 1
	}

	activeWarnings, activeEmergencies := s.l.ActiveAlarmsCount()
	c.Check(activeEmergencies, Equals, expectedEmergency)
	c.Check(activeWarnings, Equals, expectedWarning)
}
