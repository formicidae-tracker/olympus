package main

import (
	"time"

	"github.com/dgryski/go-lttb"
)

// A DataRollingSampler is used to keep trace of a time series over a
// certain time window, keeping the number of sample as low
// as possible.
type DataRollingSampler interface {
	// Adds a new point to the sampler
	Add(t time.Duration, v float64)
	// Returns the resulting time serie
	TimeSerie() []lttb.Point
}

type rollingDownsampler struct {
	xPeriod   float64
	threshold int
	points    []lttb.Point
}

func NewRollingSampler(window time.Duration, nbSamples int) DataRollingSampler {
	res := &rollingDownsampler{
		xPeriod:   window.Seconds(),
		threshold: nbSamples,
		points:    make([]lttb.Point, 0, nbSamples),
	}
	return res
}

func (d *rollingDownsampler) Add(t time.Duration, v float64) {
	d.points = append(d.points, lttb.Point{t.Seconds(), v})
	idx := 0
	last := d.points[len(d.points)-1].X
	for {
		if (last - d.points[idx].X) <= d.xPeriod {
			break
		}
		idx += 1
	}

	if idx != 0 {
		d.points = d.points[idx:]
	}
}

func (d *rollingDownsampler) TimeSerie() []lttb.Point {
	return lttb.LTTB(d.points, d.threshold)
}
