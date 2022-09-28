package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type benchmarkData struct {
	time   []time.Time
	values [][]*float32

	asTimedValue timedValues
}

type DataRollingSamplerSuite struct {
	benchmarkData benchmarkData
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&DataRollingSamplerSuite{})

func (s *DataRollingSamplerSuite) SetUpSuite(c *C) {
	var week = 7 * 24 * time.Hour
	var period = 500 * time.Millisecond
	var nbSamples = week / period
	var times = make([]time.Time, nbSamples)
	var values = make([][]*float32, nbSamples)

	for i := range times {
		times[i] = time.Unix(0, 0).Add(time.Duration(i) * period)
		values[i] = []*float32{new(float32), new(float32), new(float32), new(float32)}
	}

	var tValues timedValues
	for i, t := range times {
		tValues.push(t, values[i], 3*time.Minute)
	}

	s.benchmarkData = benchmarkData{times, values, tValues}
}

func (s *DataRollingSamplerSuite) TestSamplesData(c *C) {
	tick := 1 * time.Second
	jitter := 30 * time.Millisecond

	sampler := NewRollingSampler(time.Minute, 10, time.Millisecond)

	jr := 0.0
	for current := time.Unix(0, 0); current.Before(time.Unix(180, 0)); current = current.Add(tick + time.Duration(jr*float64(jitter))) {

		jr = 2*rand.Float64() - 1.0
		sampler.Add(current, []*float32{newInitialized[float32](float32(jr))}, nil)
		points := sampler.TimeSerie()[0]
		if c.Check(len(points) > 0, Equals, true) == false {
			continue
		}

		c.Check((points[len(points)-1].X-points[0].X) <= 60.0, Equals, true)
		c.Check(len(points) <= 10.0, Equals, true)
	}

}

func (s *DataRollingSamplerSuite) TestAlwaysXSortData(c *C) {
	sampler := NewRollingSampler(time.Minute, 60, time.Millisecond)
	for i := 0; i < 120; i++ {
		delta := time.Duration(rand.Intn(60000)) * time.Millisecond
		t := time.Unix(0, 0).Add(delta)
		sampler.Add(t, []*float32{new(float32)}, nil)
	}
	result := sampler.TimeSerie()[0]
	for i := 1; i < len(result); i++ {
		c.Check(result[i-1].X <= result[i].X, Equals, true, Commentf(" at index %i in %v", i, result))
	}

}

func (s *DataRollingSamplerSuite) TestAsyncEventuallyTerminate(c *C) {
	sampler := NewRollingSampler(time.Minute, 60, time.Millisecond)
	mx := sync.RWMutex{}

	for i := 0; i < 60; i++ {
		t := time.Unix(int64(i), 0)
		sampler.Add(t, []*float32{new(float32)}, &mx)
	}

	for j := 0; j < 20; j++ {
		series := sampler.TimeSerie()
		if len(series) == 1 && len(series[0]) == 60 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	series := sampler.TimeSerie()
	c.Assert(series, HasLen, 1)
	c.Check(series[0], HasLen, 60)
}

func (s *DataRollingSamplerSuite) BenchmarkWeek(c *C) {

	for i := 0; i < c.N; i++ {
		sampler := NewRollingSampler(7*24*time.Hour, 400, 3*time.Minute)
		sampler.AddBatch(s.benchmarkData.time, s.benchmarkData.values, nil)
	}

}

func (s *DataRollingSamplerSuite) BenchmarkLTTB(c *C) {
	for i := 0; i < c.N; i++ {
		downsample(s.benchmarkData.asTimedValue, 300)
	}
}
