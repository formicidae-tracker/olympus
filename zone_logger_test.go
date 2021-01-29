package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
	. "gopkg.in/check.v1"
)

type ZoneLoggerSuite struct {
	l ZoneLogger
}

var _ = Suite(&ZoneLoggerSuite{})

func (s *ZoneLoggerSuite) SetUpTest(c *C) {
	s.l = NewZoneLogger(zeus.ZoneRegistration{
		Host:   "foo",
		Name:   "bar",
		NumAux: 0,
	}, 1*time.Millisecond)
}

func (s *ZoneLoggerSuite) TearDownTest(c *C) {
	err := s.l.Close()
	c.Check(err, IsNil)
}

func (s *ZoneLoggerSuite) TestDoubleCloseDontPanic(c *C) {
	c.Check(s.l.Close(), IsNil)
	c.Check(s.l.Close(), ErrorMatches, "ZoneLogger: already closed")
	s.SetUpTest(c) // to avoid error on tear down
}

func (s *ZoneLoggerSuite) TestLogsClimate(c *C) {
	start := time.Now().Round(0)
	for i := 0; i < 20; i++ {
		s.l.ReportChannel() <- zeus.ClimateReport{
			Time:         start.Add(time.Duration(rand.Intn(20000)) * time.Millisecond),
			Humidity:     40.0,
			Temperatures: []zeus.Temperature{20.0},
		}
	}
	for {
		received := len(s.l.GetClimateReportSeries("10m").Humidity)
		if received == 20 {
			break
		}
		time.Sleep(10 * time.Millisecond)
		if time.Now().After(start.Add(500*time.Millisecond)) == true {
			c.Fatalf("did not received all report after 500ms")
		}
	}

	checkReport := func(c *C, series ClimateReportTimeSerie) {
		if c.Check(len(series.Humidity), Equals, len(series.TemperatureAnt)) == false {
			return
		}
		for i, _ := range series.Humidity {
			if i == 0 {
				continue
			}
			c.Check(series.Humidity[i-1].X <= series.Humidity[i].X,
				Equals,
				true, Commentf("at index %i", i))
			c.Check(series.TemperatureAnt[i-1].X <= series.TemperatureAnt[i].X,
				Equals,
				true, Commentf("at index %i", i))

		}
	}

	checkReport(c, s.l.GetClimateReportSeries("10m"))
	checkReport(c, s.l.GetClimateReportSeries("1h"))
	checkReport(c, s.l.GetClimateReportSeries("1d"))
	checkReport(c, s.l.GetClimateReportSeries("1w"))
}

func (s *ZoneLoggerSuite) TestLogsAlarms(c *C) {
	start := time.Now().Round(0)

	events := []zeus.AlarmEvent{
		zeus.AlarmEvent{
			Reason: "foo",
			Flags:  zeus.Warning | zeus.InstantNotification,
		},
		zeus.AlarmEvent{
			Reason: "bar",
			Flags:  zeus.Emergency,
		},
		zeus.AlarmEvent{
			Reason: "baz",
			Flags:  zeus.Warning,
		},
	}

	lastState := map[string]struct {
		Time time.Time
		On   bool
	}{}

	for i := 0; i < 200; i++ {
		r := rand.Intn(2000000)
		t := start.Add(time.Duration(r) * time.Millisecond)
		on := r%2 == 0
		ae := events[i%3]
		ae.Time = t
		ae.Status = zeus.AlarmStatus(r % 2)
		ls := lastState[ae.Reason]
		if ls.Time.Before(t) {
			ls.Time = t
			ls.On = on
			lastState[ae.Reason] = ls
		}
		s.l.AlarmChannel() <- ae
	}

	logs := s.l.GetAlarmsEventLog()
	for ; len(logs) < 200; logs = s.l.GetAlarmsEventLog() {
		if time.Now().After(start.Add(500*time.Millisecond)) == true {
			c.Fatalf("Did not received all logs afer 500ms")
		}
		time.Sleep(10 * time.Millisecond)
	}
	for i, l := range logs {
		switch l.Reason {
		case "foo":
			c.Check(l.Level, Equals, 1)
		case "bar":
			c.Check(l.Level, Equals, 2)
		case "baz":
			c.Check(l.Level, Equals, 1)
		}
		if i == 0 {
			continue
		}
		c.Check(logs[i-1].Time.After(logs[i].Time), Equals, false)

	}

	expectedWarning := 0
	expectedEmergency := 0
	if lastState["bar"].On == true {
		expectedEmergency = 1
	}
	if lastState["foo"].On == true {
		expectedWarning += 1
	}
	if lastState["baz"].On == true {
		expectedWarning += 1
	}

	reports := s.l.GetReport()
	c.Check(reports.ActiveEmergencies, Equals, expectedEmergency)
	c.Check(reports.ActiveWarnings, Equals, expectedWarning)
}

func (s *ZoneLoggerSuite) TestStressConcurrentAccess(c *C) {
	allgo := make(chan struct{})
	start := time.Now().Round(0)
	wg := sync.WaitGroup{}
	for i := 0; i < 300; i++ {
		wg.Add(3)
		t := start.Add(time.Duration(rand.Intn(20000)) * time.Millisecond)
		go func(t time.Time) {
			r := zeus.ClimateReport{
				Time:         t,
				Humidity:     20,
				Temperatures: []zeus.Temperature{0},
			}

			<-allgo
			s.l.ReportChannel() <- r
			wg.Done()
		}(t)
		go func(t time.Time) {
			ae := zeus.AlarmEvent{
				Reason: "heheh",
				Time:   t,
			}
			<-allgo
			s.l.AlarmChannel() <- ae
			wg.Done()
		}(t)
		go func() {
			sr := zeus.StateReport{}
			<-allgo
			s.l.StateChannel() <- sr
			wg.Done()
		}()
	}
	for i := 0; i < 100; i++ {
		wg.Add(3)
		go func() {
			<-allgo
			s.l.GetClimateReportSeries("10m")
			wg.Done()
		}()
		go func() {
			<-allgo
			s.l.GetAlarmsEventLog()
			wg.Done()
		}()
		go func() {
			<-allgo
			s.l.GetReport()
			wg.Done()
		}()
	}

	//now we should release the hounds
	close(allgo)
	// we wait for the stress test to end
	wg.Wait()
}

func (s *ZoneLoggerSuite) TestSignalTimeout(c *C) {
	pollTimeout := func() bool {
		select {
		case <-s.l.Timeouted():
			return true
		default:
			return false
		}
	}
	period := s.l.(*zoneLogger).timeoutPeriod
	c.Check(pollTimeout(), Equals, false)
	for i := 0; i < 10; i++ {
		s.l.StateChannel() <- zeus.StateReport{}
		c.Check(pollTimeout(), Equals, false)
		// if we maitain communication for much less than the time
		// period, we should not timeout
		time.Sleep(period / 4)
	}

	// if we do nothing for more than the timeout period, we should
	// timeout
	time.Sleep(2 * period)
	c.Check(pollTimeout(), Equals, true)
	// timeouting twice won't panic
	time.Sleep(2 * period)
}
