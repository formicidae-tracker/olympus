package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/formicidae-tracker/olympus/proto"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

type OlympusSuite struct {
	o                                      *Olympus
	somehostBox, anotherBox, anotherTunnel ZoneLogger
}

var _ = Suite(&OlympusSuite{})

func (s *OlympusSuite) SetUpTest(c *C) {
	var err error
	s.o, err = NewOlympus("")
	c.Assert(err, IsNil)
	s.somehostBox, _, err = s.o.RegisterZone(&proto.ZoneDeclaration{
		Host:        "somehost",
		Name:        "box",
		NumAux:      0,
		HasHumidity: true,
	})
	c.Assert(err, IsNil)

	s.anotherBox, _, err = s.o.RegisterZone(&proto.ZoneDeclaration{
		Host:   "another",
		Name:   "box",
		NumAux: 0,
	})
	c.Assert(err, IsNil)

	s.anotherTunnel, _, err = s.o.RegisterZone(&proto.ZoneDeclaration{
		Host:   "another",
		Name:   "tunnel",
		NumAux: 0,
	})
	c.Assert(err, IsNil)

	hostname, err := os.Hostname()
	c.Assert(err, IsNil)

	_, _, err = s.o.RegisterTracker(&proto.TrackingDeclaration{
		Hostname:       "somehost",
		StreamServer:   hostname + ".local",
		ExperimentName: "TEST-MODE",
	})
	c.Assert(err, IsNil)

	_, _, err = s.o.RegisterTracker(&proto.TrackingDeclaration{
		Hostname:       "fifou",
		StreamServer:   hostname + ".local",
		ExperimentName: "TEST-MODE",
	})
	c.Assert(err, IsNil)
}

func newInitialized[T any](v T) *T {
	res := new(T)
	*res = v
	return res
}

func (s *OlympusSuite) TearDownTest(c *C) {
	c.Check(s.o.Close(), IsNil)
}

func (s *OlympusSuite) TestReportClimate(c *C) {
	reports := make([]*proto.ClimateReport, 20)
	for i := 0; i < 20; i++ {
		reports[i] = &proto.ClimateReport{
			Humidity:     newInitialized[float32](20.0),
			Temperatures: []float32{20},
		}
	}
	s.somehostBox.PushReports(reports)

	windows := []string{"10m", "1h", "1d", "1w", "10-minutes", "10-minute", "hour", "day", "week", "will default to 10 minutes if window is not a valid one"}
	for _, w := range windows {
		series, err := s.o.GetClimateTimeSerie("somehost", "box", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, 20, Commentf("for window %s", w))
		c.Check(series.TemperatureAnt, HasLen, 20, Commentf("for window %s", w))

		series, err = s.o.GetClimateTimeSerie("another", "box", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, 0, Commentf("for window %s", w))
		c.Check(series.TemperatureAnt, HasLen, 0, Commentf("for window %s", w))

		series, err = s.o.GetClimateTimeSerie("another", "tunnel", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, 0, Commentf("for window %s", w))
		c.Check(series.TemperatureAnt, HasLen, 0, Commentf("for window %s", w))
	}

	_, err := s.o.GetClimateTimeSerie("fifou", "bar", "10m")
	c.Check(err, ErrorMatches, "olympus: unknown zone 'fifou.bar'")
	r, err := s.o.GetZoneReport("fifou", "bar")
	c.Check(err, IsNil)
	c.Check(r.Climate, IsNil)
	if c.Check(r.Stream, NotNil) == true {
		c.Check(r.Stream.StreamURL, Matches, "/olympus/hls/fifou.m3u8")
	}

	report, err := s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Check(*report.Climate.Humidity, Equals, float32(20.0))
	c.Check(*report.Climate.Temperature, Equals, float32(20.0))

	report, err = s.o.GetZoneReport("another", "box")
	c.Check(err, IsNil)
	c.Check(report.Climate.Humidity, IsNil)
	c.Check(report.Climate.Temperature, IsNil)
}

func isSorted(n int, comp func(i, j int) bool) bool {
	for i := 1; i < n; i++ {
		if comp(i-1, i) == false {
			return false
		}
	}
	return true
}

func (s *OlympusSuite) TestZoneSummary(c *C) {
	summary := s.o.GetZones()
	//we check that it is sorted
	c.Check(isSorted(len(summary), func(i, j int) bool {
		if summary[i].Host == summary[j].Host {
			return summary[i].Name < summary[j].Name
		}
		return summary[i].Host < summary[j].Host
	}), Equals, true)
	c.Assert(summary, HasLen, 4)
	c.Check(summary[0].Host, Matches, "another")
	c.Check(summary[0].Name, Matches, "box")
	c.Check(summary[0].Climate, NotNil)
	c.Check(summary[0].Stream, IsNil)

	c.Check(summary[1].Host, Matches, "another")
	c.Check(summary[1].Name, Matches, "tunnel")
	c.Check(summary[1].Climate, NotNil)
	c.Check(summary[1].Stream, IsNil)

	c.Check(summary[2].Host, Matches, "fifou")
	c.Check(summary[2].Name, Matches, "box")
	c.Check(summary[2].Climate, IsNil)
	if c.Check(summary[2].Stream, NotNil) == true {
		c.Check(summary[2].Stream.StreamURL, Matches, "/olympus/hls/fifou.m3u8")
		c.Check(summary[2].Stream.ThumbnailURL, Matches, "/olympus/fifou.png")
	}

	c.Check(summary[3].Host, Matches, "somehost")
	c.Check(summary[3].Name, Matches, "box")
	c.Check(summary[3].Climate, NotNil)
	if c.Check(summary[3].Stream, NotNil) == true {
		c.Check(summary[3].Stream.StreamURL, Matches, "/olympus/hls/somehost.m3u8")
		c.Check(summary[3].Stream.ThumbnailURL, Matches, "/olympus/somehost.png")
	}

}

func (s *OlympusSuite) TestMultipleError(c *C) {
	var err multipleError
	c.Check(err, IsNil)
	err = appendError(err, nil)
	c.Check(err, IsNil)
	c.Check(err.Error(), Matches, "")

	err = appendError(err, nil, errors.New("foo"), nil)
	c.Check(err, ErrorMatches, "foo")
	err = appendError(err, errors.New("bar"))
	c.Check(err, ErrorMatches, `multiple errors:
foo
bar`)
}

func (s *OlympusSuite) TestRoute(c *C) {
	testdata := []struct {
		Method string
		URL    string
		Error  string
	}{
		{"GET", "/api/zones", ""},
		{"GET", "/api/host/somehost/zone/box", ""},
		{"GET", "/api/host/somehost/zone/box/climate?window=1d", ""},
		{"GET", "/api/host/somehost/zone/box/alarms", ""},
		{"GET", "/api/host/somehosts/zone/box", "olympus: unknown zone 'somehosts.box'\n"},
		{"GET", "/api/host/somehosts/zone/box/climate", "olympus: unknown zone 'somehosts.box'\n"},
		{"GET", "/api/host/somehosts/zone/box/alarms", "olympus: unknown zone 'somehosts.box'\n"},
	}

	router := mux.NewRouter()
	s.o.route(router)

	for _, d := range testdata {
		var err error
		req, err := http.NewRequest(d.Method, d.URL, nil)
		c.Assert(err, IsNil, Commentf("%s %s", d.Method, d.URL))
		match := mux.RouteMatch{}

		c.Assert(router.Match(req, &match), Equals, true, Commentf("%s %s", d.Method, d.URL))
		req = mux.SetURLVars(req, match.Vars)
		w := httptest.NewRecorder()
		match.Handler.ServeHTTP(w, req)
		if len(d.Error) == 0 {
			rerr, _ := ioutil.ReadAll(w.Result().Body)
			c.Check(w.Result().StatusCode, Equals, http.StatusOK, Commentf("%s %s returned:", d.Method, d.URL, string(rerr)))

		} else {
			c.Check(w.Result().StatusCode, Equals, http.StatusInternalServerError, Commentf("%s %s", d.Method, d.URL))
			res, err := ioutil.ReadAll(w.Result().Body)
			c.Check(err, IsNil)
			c.Check(string(res), Matches, d.Error)
		}

	}
}
