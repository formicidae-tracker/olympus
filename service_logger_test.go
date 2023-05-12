package main

import . "gopkg.in/check.v1"

type ServiceLoggerSuite struct{}

var _ = Suite(&ServiceLoggerSuite{})

func (s *ServiceLoggerSuite) TestKeepsLogsSorted(c *C) {
	l := NewServiceLogger()
	l.Log("zeLast", true, true)
	l.Log("aFirst", true, true)
	l.Log("aFirst", false, false)
	l.Log("zeLast", false, true)
	logs := l.Logs()
	c.Assert(logs, HasLen, 2)
	c.Assert(logs[0].Zone, Equals, "aFirst")
	c.Assert(logs[1].Zone, Equals, "zeLast")
	c.Assert(logs[0].Events, HasLen, 2)
	c.Assert(logs[1].Events, HasLen, 2)
	for i := 0; i < 2; i++ {
		c.Check(logs[0].Events[i].On, Equals, i == 0)
		c.Check(logs[1].Events[i].On, Equals, i == 0)
	}
}

func (s *ServiceLoggerSuite) TestEnforceGracefulCorrectness(c *C) {
	l := NewServiceLogger()
	l.Log("a", true, false)
	l.Log("a", true, true)
	l.Log("a", false, false)
	l.Log("a", false, true)
	logs := l.Logs()
	c.Assert(logs, HasLen, 1)
	c.Assert(logs[0].Events, HasLen, 4)
	c.Check(logs[0].Events[0].Graceful, Equals, true)
	c.Check(logs[0].Events[1].Graceful, Equals, true)
	c.Check(logs[0].Events[2].Graceful, Equals, false)
	c.Check(logs[0].Events[3].Graceful, Equals, true)
}

func (s *ServiceLoggerSuite) TestFetchLastStatus(c *C) {
	l := NewServiceLogger()
	l.Log("a", true, true)
	l.Log("b", true, true)
	l.Log("c", true, true)
	l.Log("a", false, true)
	l.Log("a", true, true)
	l.Log("b", false, true)
	c.Check(l.OnServices(), DeepEquals, []string{"a", "c"})
	c.Check(l.OffServices(), DeepEquals, []string{"b"})
}
