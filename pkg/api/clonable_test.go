package api

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	. "gopkg.in/check.v1"
)

type ClonableSuite struct{}

var _ = Suite(&ClonableSuite{})

func (s *ClonableSuite) TestAlarmUpdate(c *C) {
	testdata := []*AlarmUpdate{
		{},
		{
			Identification: "foo",
			Level:          AlarmLevel_EMERGENCY,
			Status:         AlarmStatus_OFF,
			Time:           timestamppb.Now(),
			Description:    "something",
		},
	}

	for _, d := range testdata {
		comment := Commentf("%+v", d)
		c.Check(d.Clone(), DeepEquals, d, comment)
	}

}

func newValue[T any](v T) *T {
	res := new(T)
	*res = v
	return res
}

func (s *ClonableSuite) TestClimateState(c *C) {
	testdata := []*ClimateState{
		{},
		{
			Name:         "foo",
			Temperature:  newValue[float32](20.0),
			Humidity:     newValue[float32](60.0),
			Wind:         newValue[float32](70.0),
			VisibleLight: newValue[float32](80.0),
			UvLight:      newValue[float32](10.0),
		},
	}

	for _, d := range testdata {
		comment := Commentf("%+v", d)
		c.Check(d.Clone(), DeepEquals, d, comment)
	}

}

func (s *ClonableSuite) TestZoneClimateReport(c *C) {
	testdata := []*ZoneClimateReport{
		{},
		{
			Since:       time.Now(),
			Temperature: newValue[float32](20.0),
			Humidity:    newValue[float32](60.0),
			TemperatureBounds: Bounds{
				Minimum: newValue[float32](15.0), Maximum: newValue[float32](25.0),
			},
			HumidityBounds: Bounds{
				Minimum: newValue[float32](40.0), Maximum: newValue[float32](70.0),
			},
			Current:    &ClimateState{Name: "foo"},
			CurrentEnd: nil,
			Next:       &ClimateState{Name: "foo to bar"},
			NextEnd:    &ClimateState{Name: "bar"},
			NextTime:   newValue(time.Now().Add(10 * time.Minute)),
		},
	}

	for _, d := range testdata {
		comment := Commentf("%+v", d)
		c.Check(d.Clone(), DeepEquals, d, comment)
	}

}

func (s *ClonableSuite) TestTrackingInfo(c *C) {
	testdata := []*TrackingInfo{
		{},
		{
			Since:          time.Now(),
			TotalBytes:     100,
			FreeBytes:      10,
			BytesPerSecond: 1,
			Stream: &StreamInfo{
				ExperimentName: "coucou",
				StreamURL:      "https://example.com",
				ThumbnailURL:   "https://example.com",
			},
		},
	}

	for _, d := range testdata {
		comment := Commentf("%+v", d)
		c.Check(d.Clone(), DeepEquals, d, comment)
	}

}

func (s *ClonableSuite) TestServiceLog(c *C) {
	testdata := []*ServiceLog{
		{},
		{
			Zone: "some.zone",
			Events: []*ServiceEvent{
				{
					Start:    time.Now(),
					End:      newValue(time.Now().Add(10 * time.Second)),
					Graceful: true,
				},
			},
		},
	}

	for _, d := range testdata {
		comment := Commentf("%+v", d)
		c.Check(d.Clone(), DeepEquals, d, comment)
	}

}
