package olympus

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

func init() {
	logrus.SetOutput(io.Discard)
}

type OlympusSuite struct {
	o                                      *Olympus
	somehostBox, anotherBox, anotherTunnel *GrpcSubscription[ClimateLogger]
	somehostTracking, fifouTracking        *GrpcSubscription[TrackingLogger]
	somehostTrackingDefinition             *api.TrackingDeclaration
	somehostClimateDefinition              *api.ClimateDeclaration
}

var _ = Suite(&OlympusSuite{})

var hostname string

func init() {
	var err error
	if err != nil {
		panic(err.Error())
	}
}

func (s *OlympusSuite) SetUpTest(c *C) {
	_datapath = c.MkDir()

	var err error
	s.o, err = NewOlympus()
	c.Assert(err, IsNil)
	s.somehostClimateDefinition = &api.ClimateDeclaration{
		Host: "somehost",
		Name: "box",
	}
	s.somehostBox, err = s.o.RegisterClimate(context.Background(),
		s.somehostClimateDefinition)
	c.Assert(err, IsNil)

	s.anotherBox, err = s.o.RegisterClimate(context.Background(),
		&api.ClimateDeclaration{
			Host: "another",
			Name: "box",
		})
	c.Assert(err, IsNil)

	s.anotherTunnel, err = s.o.RegisterClimate(context.Background(),
		&api.ClimateDeclaration{
			Host: "another",
			Name: "tunnel",
		})
	c.Assert(err, IsNil)

	hostname, err = os.Hostname()
	c.Assert(err, IsNil)

	s.somehostTrackingDefinition = &api.TrackingDeclaration{
		Hostname:       "somehost",
		StreamServer:   hostname + ".local",
		ExperimentName: "TEST-MODE",
	}

	s.somehostTracking, err = s.o.RegisterTracking(context.Background(), s.somehostTrackingDefinition)
	c.Assert(err, IsNil)

	s.fifouTracking, err = s.o.RegisterTracking(context.Background(),
		&api.TrackingDeclaration{
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
	ctx := context.Background()
	c.Check(s.o.UnregisterClimate(ctx, "somehost", "box", true), IsNil)
	c.Check(s.o.UnregisterClimate(ctx, "another", "box", true), IsNil)
	c.Check(s.o.UnregisterClimate(ctx, "another", "tunnel", true), IsNil)

	c.Check(s.o.UnregisterTracker(ctx, "somehost", true), IsNil)
	c.Check(s.o.UnregisterTracker(ctx, "fifou", true), IsNil)

	c.Check(s.o.Close(), IsNil)
}

func (s *OlympusSuite) TestTrackingAndClimateShareAlarmLogger(c *C) {
	c.Check(s.somehostTracking.alarmLogger, Equals, s.somehostBox.alarmLogger)
}

func (s *OlympusSuite) TestReportClimate(c *C) {
	reports := make([]*api.ClimateReport, 300)
	for i := 0; i < 300; i++ {
		reports[i] = &api.ClimateReport{
			Time:         timestamppb.New(time.Time{}.Add(time.Duration(500*i) * time.Millisecond)),
			Humidity:     newInitialized[float32](55.0),
			Temperatures: []float32{21},
		}
	}
	s.somehostBox.object.PushReports(reports)

	windows := []string{"10m", "1h", "1d", "1w", "10-minutes", "10-minute", "hour", "day", "week", "will default to 10 minutes if window is not a valid one"}
	size := []int{300, 150, 8, 2, 300, 300, 150, 8, 2, 300}
	for i, w := range windows {
		series, err := s.o.GetClimateTimeSerie("somehost", "box", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, size[i], Commentf("for window %s", w))
		c.Check(series.Temperature, HasLen, size[i], Commentf("for window %s", w))

		series, err = s.o.GetClimateTimeSerie("another", "box", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, 0, Commentf("for window %s", w))
		c.Check(series.Temperature, HasLen, 0, Commentf("for window %s", w))

		series, err = s.o.GetClimateTimeSerie("another", "tunnel", w)
		c.Check(err, IsNil, Commentf("for window %s", w))
		c.Check(series.Humidity, HasLen, 0, Commentf("for window %s", w))
		c.Check(series.Temperature, HasLen, 0, Commentf("for window %s", w))
	}

	_, err := s.o.GetClimateTimeSerie("fifou", "bar", "10m")
	c.Check(err, ErrorMatches, "olympus: unknown zone 'fifou.bar'")
	r, err := s.o.GetZoneReport("fifou", "bar")
	c.Check(err, IsNil)
	c.Check(r.Climate, IsNil)
	if c.Check(r.Tracking, NotNil) == true {
		if c.Check(r.Tracking.Stream, NotNil) == true {
			c.Check(r.Tracking.Stream.StreamURL, Matches, "/olympus/fifou/index.m3u8")
		}
	}

	report, err := s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Check(*report.Climate.Humidity, Equals, float32(55.0))
	c.Check(*report.Climate.Temperature, Equals, float32(21.0))

	report, err = s.o.GetZoneReport("another", "box")
	c.Check(err, IsNil)
	c.Check(report.Climate.Humidity, IsNil)
	c.Check(report.Climate.Temperature, IsNil)
}

func (s *OlympusSuite) TestZoneSummary(c *C) {
	summary := s.o.GetZones()
	//we check that it is sorted
	c.Check(sort.SliceIsSorted(summary, func(i, j int) bool {
		if summary[i].Host == summary[j].Host {
			return summary[i].Name < summary[j].Name
		}
		return summary[i].Host < summary[j].Host
	}), Equals, true)
	c.Assert(summary, HasLen, 4)
	c.Check(summary[0].Host, Matches, "another")
	c.Check(summary[0].Name, Matches, "box")
	c.Check(summary[0].Climate, NotNil)
	c.Check(summary[0].Tracking, IsNil)

	c.Check(summary[1].Host, Matches, "another")
	c.Check(summary[1].Name, Matches, "tunnel")
	c.Check(summary[1].Climate, NotNil)
	c.Check(summary[1].Tracking, IsNil)

	c.Check(summary[2].Host, Matches, "fifou")
	c.Check(summary[2].Name, Matches, "box")
	c.Check(summary[2].Climate, IsNil)
	if c.Check(summary[2].Tracking, NotNil) == true {
		c.Check(summary[2].Tracking.Stream.StreamURL, Matches, "/olympus/fifou/index.m3u8")
		c.Check(summary[2].Tracking.Stream.ThumbnailURL, Matches, "/thumbnails/olympus/fifou.jpg")
	}

	c.Check(summary[3].Host, Matches, "somehost")
	c.Check(summary[3].Name, Matches, "box")
	c.Check(summary[3].Climate, NotNil)
	if c.Check(summary[3].Tracking.Stream, NotNil) == true {
		c.Check(summary[3].Tracking.Stream.StreamURL, Matches, "/olympus/somehost/index.m3u8")
		c.Check(summary[3].Tracking.Stream.ThumbnailURL, Matches, "/thumbnails/olympus/somehost.jpg")
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
	s.o.setRoutes(router)

	for _, d := range testdata {
		comment := Commentf("%s: '%s'", d.Method, d.URL)
		var err error
		req, err := http.NewRequest(d.Method, d.URL, nil)
		c.Assert(err, IsNil, comment)
		match := mux.RouteMatch{}

		c.Assert(router.Match(req, &match), Equals, true, comment)
		req = mux.SetURLVars(req, match.Vars)
		w := httptest.NewRecorder()
		match.Handler.ServeHTTP(w, req)
		if len(d.Error) == 0 {
			rerr, _ := ioutil.ReadAll(w.Result().Body)
			c.Check(w.Result().StatusCode,
				Equals,
				http.StatusOK,
				Commentf("%s: '%s' returned: %s", d.Method, d.URL, string(rerr)))
		} else {
			c.Check(w.Result().StatusCode, Equals, http.StatusNotFound, comment)
			res, err := ioutil.ReadAll(w.Result().Body)
			c.Check(err, IsNil, comment)
			c.Check(string(res), Matches, d.Error, comment)
		}

	}
}

func (s *OlympusSuite) TestUnregisterTracking(c *C) {
	report, err := s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Assert(report, Not(IsNil))
	c.Check(report.Climate, Not(IsNil))
	c.Check(report.Tracking, Not(IsNil))

	c.Assert(s.o.UnregisterTracker(context.Background(), "somehost", true), IsNil)
	defer func() {
		s.somehostTracking, _ = s.o.RegisterTracking(context.Background(),
			s.somehostTrackingDefinition)
	}()
	report, err = s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Assert(report, Not(IsNil))
	c.Check(report.Climate, Not(IsNil))
	c.Check(report.Tracking, IsNil)

}

func (s *OlympusSuite) TestUnregisterClimate(c *C) {
	report, err := s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Assert(report, Not(IsNil))
	c.Check(report.Climate, Not(IsNil))
	c.Check(report.Tracking, Not(IsNil))

	c.Assert(s.o.UnregisterClimate(context.Background(), "somehost", "box", true), IsNil)
	defer func() {
		s.somehostBox, _ = s.o.RegisterClimate(context.Background(),
			s.somehostClimateDefinition)
	}()
	report, err = s.o.GetZoneReport("somehost", "box")
	c.Check(err, IsNil)
	c.Assert(report, Not(IsNil))
	c.Check(report.Climate, IsNil)
	c.Check(report.Tracking, Not(IsNil))

}
