package main

import (
	"sort"
	"time"

	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type AlarmUpdateFilterSuite struct {
	unfiltered, filtered chan ZonedAlarmUpdate
}

var _ = Suite(&AlarmUpdateFilterSuite{})

func (s *AlarmUpdateFilterSuite) SetUpTest(c *C) {
	s.unfiltered = make(chan ZonedAlarmUpdate)
	s.filtered = make(chan ZonedAlarmUpdate)
}

func (s *AlarmUpdateFilterSuite) TearDownTest(c *C) {
	if s.unfiltered != nil {
		close(s.unfiltered)
	}
}

func (s *AlarmUpdateFilterSuite) TestPropagateClosing(c *C) {
	done := make(chan struct{})

	go func(incoming <-chan ZonedAlarmUpdate) {
		defer close(done)
		FilterAlarmUpdates(1*time.Millisecond)(s.filtered, incoming)
	}(s.unfiltered)

	close(s.unfiltered)
	s.unfiltered = nil

	grace := 10 * time.Millisecond
	select {
	case <-done:
	case <-time.After(grace):
		c.Fatalf("filter did not exit after %s", grace)
	}
	select {
	case _, ok := <-s.filtered:
		c.Check(ok, Equals, false)
	default:
		c.Errorf("filter did not close outgoing channel on close")
	}

}

func (s *AlarmUpdateFilterSuite) TestFilterSpuriousAlarm(c *C) {
	testdata := []struct {
		At     time.Duration
		Update ZonedAlarmUpdate
	}{
		{At: 0, Update: ZonedAlarmUpdate{Zone: "nice.zone",
			Update: &api.AlarmUpdate{
				Identification: "nice",
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 6 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "nice.zone",
			Update: &api.AlarmUpdate{
				Identification: "nice",
			},
		}},

		{At: 500 * time.Microsecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
				Level:          api.AlarmLevel_WARNING,
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 7 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
				Level:          api.AlarmLevel_EMERGENCY,
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 13 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
			},
		}},
	}

	sort.Slice(testdata, func(i, j int) bool {
		return testdata[i].At < testdata[j].At
	})

	go func(incoming <-chan ZonedAlarmUpdate) {
		FilterAlarmUpdates(5*time.Millisecond)(s.filtered, incoming)
	}(s.unfiltered)

	go func() {
		i := 0
		start := time.Now()
		for now := range time.Tick(100 * time.Microsecond) {
			ellapsed := now.Sub(start)
			for ; i < len(testdata); i++ {
				if testdata[i].At > ellapsed {
					break
				}
				update := testdata[i].Update
				update.Update.Time = timestamppb.New(now)
				s.unfiltered <- testdata[i].Update
			}
			if i >= len(testdata) {
				break
			}
		}
		close(s.unfiltered)
		s.unfiltered = nil
	}()

	expected := []ZonedAlarmUpdate{
		{Zone: "nice.zone", Update: &api.AlarmUpdate{Identification: "nice"}},
		{Zone: "ramping.up", Update: &api.AlarmUpdate{Identification: "diff", Level: api.AlarmLevel_WARNING}},
		{Zone: "ramping.up", Update: &api.AlarmUpdate{Identification: "diff", Level: api.AlarmLevel_EMERGENCY}},
	}

	j := 0
	for update := range s.filtered {

		if c.Check(j < len(expected), Equals, true, Commentf("got unexpected update %d:%v", j, update)) == false {
			continue
		}
		comment := Commentf("expected %d: %v", j, expected[j])
		c.Check(update.ID(), Equals, expected[j].ID(), comment)
		c.Check(update.Update.Level, Equals, expected[j].Update.Level, comment)
		j++
	}
	for _, update := range expected[j:] {
		c.Errorf("Did not receive expected update %v", update)
	}
}
