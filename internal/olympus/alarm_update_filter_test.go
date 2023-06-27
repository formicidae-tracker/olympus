package olympus

import (
	"sort"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"golang.org/x/exp/constraints"
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
		// A simple alarm
		{At: 0, Update: ZonedAlarmUpdate{Zone: "nice.zone",
			Update: &api.AlarmUpdate{
				Identification: "nice",
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 40 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "nice.zone",
			Update: &api.AlarmUpdate{
				Identification: "nice",
				Status:         api.AlarmStatus_OFF,
			},
		}},

		// An increasing alarm
		{At: 3 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
				Level:          api.AlarmLevel_WARNING,
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 10 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
				Level:          api.AlarmLevel_EMERGENCY,
				Status:         api.AlarmStatus_ON,
			},
		}},

		{At: 17 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "diff",
				Status:         api.AlarmStatus_OFF,
			},
		}},

		// Should update descriptions
		// An increasing alarm
		{At: 9 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "changing",
				Level:          api.AlarmLevel_WARNING,
				Status:         api.AlarmStatus_ON,
				Description:    "a",
			},
		}},
		{At: 10 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "changing",
				Level:          api.AlarmLevel_WARNING,
				Status:         api.AlarmStatus_ON,
				Description:    "b",
			},
		}},
		{At: 10 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "ramping.up",
			Update: &api.AlarmUpdate{
				Identification: "changing",
				Level:          api.AlarmLevel_EMERGENCY,
				Status:         api.AlarmStatus_ON,
				Description:    "",
			},
		}},

		// Warning should not fire twice
		{At: 6 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "twice",
				Status:         api.AlarmStatus_ON,
				Level:          api.AlarmLevel_WARNING,
			},
		}},
		{At: 13 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "twice",
				Status:         api.AlarmStatus_ON,
				Level:          api.AlarmLevel_WARNING,
			},
		}},
		{At: 20 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "twice",
				Status:         api.AlarmStatus_OFF,
				Level:          api.AlarmLevel_WARNING,
			},
		}},

		//Spurious alarms should be discarded
		{At: 0, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "spurious",
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 1 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "spurious",
				Status:         api.AlarmStatus_OFF,
			},
		}},
		{At: 2 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "spurious",
				Status:         api.AlarmStatus_ON,
			},
		}},
		{At: 4 * time.Millisecond, Update: ZonedAlarmUpdate{Zone: "faulty",
			Update: &api.AlarmUpdate{
				Identification: "spurious",
				Status:         api.AlarmStatus_OFF,
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
		{Zone: "nice.zone", Update: &api.AlarmUpdate{Identification: "nice", Level: api.AlarmLevel_WARNING}},
		{Zone: "ramping.up", Update: &api.AlarmUpdate{Identification: "diff", Level: api.AlarmLevel_WARNING}},
		{Zone: "faulty", Update: &api.AlarmUpdate{Identification: "twice", Level: api.AlarmLevel_WARNING}},
		{Zone: "ramping.up", Update: &api.AlarmUpdate{Identification: "changing", Level: api.AlarmLevel_EMERGENCY, Description: "b"}},
		{Zone: "ramping.up", Update: &api.AlarmUpdate{Identification: "diff", Level: api.AlarmLevel_EMERGENCY}},
	}

	received := make([]ZonedAlarmUpdate, 0, len(expected))
	for update := range s.filtered {
		received = append(received, update)
	}

	c.Check(len(received), Equals, len(expected))

	for i := 0; i < Min(len(expected), len(received)); i++ {
		comment := Commentf("alarm #%d", i)
		c.Check(received[i].ID(), Equals, expected[i].ID(), comment)
		c.Check(received[i].Update.Level, Equals, expected[i].Update.Level, comment)
		c.Check(received[i].Update.Description, Equals, expected[i].Update.Description, comment)
		c.Check(received[i].Update.Status, Equals, api.AlarmStatus_ON)
	}

	for _, u := range expected[Min(len(received), len(expected)):] {
		c.Errorf("Did not receive alarm %v", u)
	}

	for _, u := range received[Min(len(received), len(expected)):] {
		c.Errorf("Received unexpected alarm %v", u)
	}

}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
