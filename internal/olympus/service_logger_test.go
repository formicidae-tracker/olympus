package olympus

import (
	"context"

	. "gopkg.in/check.v1"
)

type ServiceLoggerSuite struct{}

var _ = Suite(&ServiceLoggerSuite{})

func (s *ServiceLoggerSuite) SetUpTest(c *C) {
	_datapath = c.MkDir()
}

func (s *ServiceLoggerSuite) TestKeepsLogsSorted(c *C) {
	l := NewServiceLogger()
	ctx := context.Background()
	l.Log(ctx, "zeLast", true, true)
	l.Log(ctx, "aFirst", true, true)
	l.Log(ctx, "aFirst", false, false)
	l.Log(ctx, "zeLast", false, true)
	logs := l.Logs()
	c.Assert(logs, HasLen, 2)
	c.Assert(logs[0].Zone, Equals, "aFirst")
	c.Assert(logs[1].Zone, Equals, "zeLast")
	c.Assert(logs[0].Events, HasLen, 1)
	c.Assert(logs[1].Events, HasLen, 1)

	c.Check(logs[0].Events[0].End, Not(IsNil))
	c.Check(logs[1].Events[0].End, Not(IsNil))

	c.Check(logs[0].Events[0].Graceful, Equals, false)
	c.Check(logs[1].Events[0].Graceful, Equals, true)

}

func (s *ServiceLoggerSuite) TestEnforceGracefulCorrectness(c *C) {
	l := NewServiceLogger()
	ctx := context.Background()
	l.Log(ctx, "a", true, false)
	l.Log(ctx, "a", true, true)
	l.Log(ctx, "a", false, false)
	l.Log(ctx, "a", false, true)
	logs := l.Logs()
	c.Assert(logs, HasLen, 1)
	c.Assert(logs[0].Events, HasLen, 1)
	c.Check(logs[0].Events[0].Graceful, Equals, false)
}

func (s *ServiceLoggerSuite) TestFetchLastStatus(c *C) {
	l := NewServiceLogger()
	ctx := context.Background()
	l.Log(ctx, "a", true, true)
	l.Log(ctx, "b", true, true)
	l.Log(ctx, "c", true, true)
	l.Log(ctx, "a", false, true)
	l.Log(ctx, "a", true, true)
	l.Log(ctx, "b", false, true)
	c.Check(l.OnServices(), DeepEquals, []string{"a", "c"})
	c.Check(l.OffServices(), DeepEquals, []string{"b"})
}
