package main

import (
	"time"

	. "gopkg.in/check.v1"
)

type TrackingWatcherSuite struct {
}

var _ = Suite(&TrackingWatcherSuite{})

func (s *TrackingWatcherSuite) TestTimeout(c *C) {
	w := newTrackingWatcher("", "", 500*time.Microsecond)
	c.Check(pollClosed(w.Timeouted()), Equals, false)
	time.Sleep(2 * time.Millisecond)
	c.Check(pollClosed(w.Timeouted()), Equals, true)
	c.Check(w.Close(), IsNil)
	c.Check(pollClosed(w.Done()), Equals, true)
}
