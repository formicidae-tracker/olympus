package main

import (
	"time"

	. "gopkg.in/check.v1"
)

type NotifierSuite struct {
	datapath string

	notifier Notifier
}

var _ = Suite(&NotifierSuite{})

func (s *NotifierSuite) SetUpSuite(c *C) {
	s.datapath = _datapath
}

func (s *NotifierSuite) TearDownSuite(c *C) {
	_datapath = s.datapath
}

func (s *NotifierSuite) SetUpTest(c *C) {
	_datapath = c.MkDir()
	s.notifier = NewNotifier(5 * time.Millisecond)
}

func (s *NotifierSuite) TestClosingOnIncoming(c *C) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.notifier.Loop()
	}()

	close(s.notifier.Incoming())

	grace := 5 * time.Millisecond
	select {
	case <-time.After(grace):
		c.Errorf("notifier did not close after %s", grace)
	case <-done:
	}
}
