package main

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	. "gopkg.in/check.v1"
)

type DataRollingSamplerSuite struct {
}

func Test(t *testing.T) { TestingT(t) }

var _ = oSuite(&DataRollingSamplerSuite{})

func (s *DataRollingSamplerSuite) TestSamplesData(c *C) {
	tick := 1 * time.Second
	jitter := 30 * time.Millisecond

	sampler := NewRollingSampler(time.Minute, 10)

	jr := 0.0
	for current := time.Duration(0); current < 3*time.Minute; current += tick + time.Duration(jr*float64(jitter)) {
		jr = 2*rand.Float64() - 1.0
		sampler.Add(current, jr)
		points := sampler.TimeSerie()
		if c.Check(len(points) > 0, Equals, true) == false {
			continue
		}
		c.Check((points[len(points)-1].X-points[0].X) <= 60.0, Equals, true)
		c.Check(len(points) <= 10.0, Equals, true)
	}

}

func (s *DataRollingSamplerSuite) TestAlwaysXSortData(c *C) {
	sampler := NewRollingSampler(time.Minute, 60)
	for i := 0; i < 120; i++ {
		t := time.Duration(rand.Intn(60000)) * time.Millisecond
		sampler.Add(t, 1.0)
	}
	result := sampler.TimeSerie()
	for i := 1; i < len(result); i++ {
		c.Check(result[i-1].X <= result[i].X, Equals, true, Commentf(" at index %i in %v", i, result))
	}

}

func (s *DataRollingSamplerSuite) TestAlwaysReturnACopy(c *C) {
	sampler := NewRollingSampler(time.Minute, 5)
	sampler.Add(time.Duration(0), 23.0)

	a := sampler.TimeSerie()
	b := sampler.TimeSerie()

	c.Check(unsafe.Pointer(a), Not(Equals), unsafe.Pointer(b))
}
