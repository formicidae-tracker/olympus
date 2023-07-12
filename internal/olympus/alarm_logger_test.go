package olympus

import (
	"math/rand"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
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

	eventList := []*api.AlarmUpdate{
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

	events := make([]*api.AlarmUpdate, 500)
	for i := 0; i < 500; i++ {
		r := rand.Intn(2000000)
		t := start.Add(time.Duration(r) * time.Millisecond)
		on := i%2 == 0
		event := eventList[i%3].Clone()
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
		c.Check(len(r.Events) > 0, Equals, true)
		for i, e := range r.Events {
			if i == 0 {
				continue
			}
			previousEvent := r.Events[i-1]
			if c.Check(previousEvent.End, Not(IsNil)) == false {
				continue
			}
			c.Check(previousEvent.End.After(e.Start), Equals, false,
				Commentf("Got consecutive events %+v %+v", previousEvent, e))
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

	activeFailures, activeEmergencies, activeWarnings := s.l.ActiveAlarmsCount()
	c.Check(activeFailures, Equals, 0)
	c.Check(activeEmergencies, Equals, expectedEmergency)
	c.Check(activeWarnings, Equals, expectedWarning)
}
