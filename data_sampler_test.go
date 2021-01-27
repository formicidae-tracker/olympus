package main

import (
	"math/rand"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type DataRollingSamplerSuite struct {
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&DataRollingSamplerSuite{})

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
