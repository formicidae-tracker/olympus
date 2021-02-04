package main

import (
	"time"

	. "gopkg.in/check.v1"
)

type TrackingWatcherSuite struct {
}

var _ = Suite(&TrackingWatcherSuite{})

func (s *TrackingWatcherSuite) TestTimeout(c *C) {
	w := newTrackingWatcher(TrackingWatcherArgs{"", "", ""}, 500*time.Microsecond)
	c.Check(pollClosed(w.Timeouted()), Equals, false)
	start := time.Now()
	for pollClosed(w.Timeouted()) == false && time.Since(start) < 10*time.Second {
		time.Sleep(5 * time.Millisecond)
	}
	c.Check(pollClosed(w.Timeouted()), Equals, true)
	c.Check(w.Close(), IsNil)
	c.Check(pollClosed(w.Done()), Equals, true)
}
