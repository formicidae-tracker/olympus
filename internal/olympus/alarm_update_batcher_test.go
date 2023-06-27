package olympus

import (
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	. "gopkg.in/check.v1"
)

type BatchAlarmUpdateSuite struct {
	incoming chan ZonedAlarmUpdate
	outgoing chan []ZonedAlarmUpdate
}

var _ = Suite(&BatchAlarmUpdateSuite{})

func (s *BatchAlarmUpdateSuite) SetUpTest(c *C) {
	s.incoming = make(chan ZonedAlarmUpdate)
	s.outgoing = make(chan []ZonedAlarmUpdate)
}

func (s *BatchAlarmUpdateSuite) TearDownTest(c *C) {
	if s.incoming != nil {
		close(s.incoming)
	}
}

func (s *BatchAlarmUpdateSuite) TestCloseOutgoing(c *C) {
	done := make(chan struct{})

	go func(incoming chan ZonedAlarmUpdate) {
		defer close(done)
		BatchAlarmUpdate(1*time.Millisecond)(s.outgoing, incoming)
	}(s.incoming)

	close(s.incoming)
	s.incoming = nil

	grace := 5 * time.Millisecond
	select {
	case <-done:
	case <-time.After(grace):
		c.Errorf("filter did not terminate after %s", grace)
	}

	select {
	case _, ok := <-s.outgoing:
		c.Check(ok, Equals, false)
	default:
		c.Errorf("filter did not close outgoing channel")
	}
}

func buildAlarm(zones string) []ZonedAlarmUpdate {
	res := make([]ZonedAlarmUpdate, 0, len(zones))
	for i := range zones {
		res = append(res, ZonedAlarmUpdate{
			Zone: string([]byte{zones[i]}),
			Update: &api.AlarmUpdate{
				Identification: "foo",
				Status:         api.AlarmStatus_ON,
			},
		})
	}
	return res
}

func (s *BatchAlarmUpdateSuite) TestBatching(c *C) {
	inputs := buildAlarm("abcdefghijklmn")
	pauses := map[int]time.Duration{
		4: 5 * time.Millisecond,
		9: 5 * time.Millisecond,
	}
	expected := [][]ZonedAlarmUpdate{
		buildAlarm("a"),
		buildAlarm("bcde"),
		buildAlarm("fghij"),
		buildAlarm("klmn"),
	}

	go func(incoming <-chan ZonedAlarmUpdate) {
		BatchAlarmUpdate(5*time.Millisecond)(s.outgoing, incoming)
	}(s.incoming)

	go func() {
		defer func() {
			close(s.incoming)
			s.incoming = nil
		}()
		for i, u := range inputs {
			s.incoming <- u
			if p, ok := pauses[i]; ok == true {
				time.Sleep(p)
			}
		}
		time.Sleep(6 * time.Millisecond)
	}()

	result := make([][]ZonedAlarmUpdate, 0, len(expected))

	for r := range s.outgoing {
		result = append(result, r)
	}

	c.Check(len(result), Equals, len(expected))
	size := Min(len(result), len(expected))
	for i := 0; i < size; i++ {
		c.Check(len(result[i]), Equals, len(expected[i]), Commentf("batch #%d", i))
		resSize := Min(len(result[i]), len(expected[i]))
		for j := 0; j < resSize; j++ {
			c.Check(result[i][j].ID(), Equals, expected[i][j].ID(), Commentf("result #%d.%d", i, j))
		}

		for _, r := range result[i][resSize:] {
			c.Errorf("Unexpected result %v in batch #%d", r, i)
		}

		for _, r := range expected[i][resSize:] {
			c.Errorf("Missing result %v in batch #%d", r, i)
		}

	}

	for _, b := range result[size:] {
		c.Errorf("Unexpected batch %s", b)
	}

	for _, b := range expected[size:] {
		c.Errorf("Missing batch %s", b)
	}

}

func (s *BatchAlarmUpdateSuite) TestDisabledBatching(c *C) {
	inputs := buildAlarm("abcdefghijklmn")
	expected := make([][]ZonedAlarmUpdate, 0, len(inputs))
	for _, i := range inputs {
		expected = append(expected, []ZonedAlarmUpdate{i})
	}

	go func(incoming <-chan ZonedAlarmUpdate) {
		BatchAlarmUpdate(0)(s.outgoing, incoming)
	}(s.incoming)

	go func() {
		defer func() {
			close(s.incoming)
			s.incoming = nil
		}()
		for _, u := range inputs {
			s.incoming <- u
			time.Sleep(10 * time.Microsecond)
		}
		time.Sleep(1 * time.Millisecond)
	}()

	result := make([][]ZonedAlarmUpdate, 0, len(expected))

	for r := range s.outgoing {
		result = append(result, r)
	}

	c.Check(len(result), Equals, len(expected))

}
