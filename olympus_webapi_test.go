package main

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/atuleu/go-lttb"
	. "gopkg.in/check.v1"
)

type WebAPISuite struct {
	benchmarkData []lttb.Point[float32]
}

var _ = Suite(&WebAPISuite{})

func (s *WebAPISuite) SetUpSuite(c *C) {
	s.benchmarkData = make([]lttb.Point[float32], 6000)
	for i := range s.benchmarkData {
		s.benchmarkData[i] = lttb.Point[float32]{X: rand.Float32(),
			Y: rand.Float32()}
	}
}

func (s *WebAPISuite) TestPointJSONEncoding(c *C) {
	testdata := []struct {
		Input    Point
		Expected string
		Error    string
	}{
		{Point{}, `{"x":0,"y":0}`, ""},
		{Point{1.2, 1.4}, `{"x":1.2,"y":1.4}`, ""},
		{Point{X: float32(math.NaN())}, ``, ".*json: unsupported value: NaN"},
		{Point{Y: float32(math.NaN())}, ``, ".*json: unsupported value: NaN"},
	}

	for _, d := range testdata {
		comment := Commentf("Test fixture %+v", d)
		res, err := json.Marshal(d.Input)
		if len(d.Error) == 0 {
			c.Check(err, IsNil, comment)
			c.Check(string(res), Equals, d.Expected, comment)
		} else {
			c.Check(err, ErrorMatches, d.Error, comment)
			c.Check(res, IsNil, comment)
		}
	}
}

func (s *WebAPISuite) TestPointSeriesJSONEncoding(c *C) {
	testdata := []struct {
		Input    PointSeries
		Expected string
		Error    string
	}{
		{PointSeries{}, `[]`, ""},
		{PointSeries{{X: 1.2, Y: 1.4}}, `[{"x":1.2,"y":1.4}]`, ""},
		{PointSeries{{X: float32(math.NaN())}}, ``, ".*json: unsupported value: NaN"},
		{PointSeries{{Y: float32(math.NaN())}}, ``, ".*json: unsupported value: NaN"},
	}

	for _, d := range testdata {
		comment := Commentf("Test fixture %+v", d)
		res, err := json.Marshal(d.Input)
		if len(d.Error) == 0 {
			c.Check(err, IsNil, comment)
			c.Check(string(res), Equals, d.Expected, comment)
		} else {
			c.Check(err, ErrorMatches, d.Error, comment)
			c.Check(res, IsNil, comment)
		}
	}
}

func (s *WebAPISuite) BenchmarkPointSeriesJSONEncoding(c *C) {
	for i := 0; i < c.N; i++ {
		json.Marshal(PointSeries(s.benchmarkData))
	}
}

func (s *WebAPISuite) BenchmarkPointSliceJSONEncoding(c *C) {
	for i := 0; i < c.N; i++ {
		json.Marshal(s.benchmarkData)
	}
}
