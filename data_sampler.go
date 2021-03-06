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
	window           float64
	samples          int
	samplesThreshold int
	points           []lttb.Point
}

func NewRollingSampler(window time.Duration, nbSamples int) DataRollingSampler {
	res := &rollingDownsampler{
		window:           window.Seconds(),
		samples:          int(0.9 * float64(nbSamples)),
		samplesThreshold: nbSamples,
		points:           make([]lttb.Point, 0, nbSamples),
	}
	return res
}

func (d *rollingDownsampler) outdated(x float64) bool {
	if len(d.points) != 0 && (d.points[len(d.points)-1].X-x) > d.window {
		return true
	}
	return false
}

func (d *rollingDownsampler) insertionSort(x, v float64) {
	i := BackLinearSearch(len(d.points), func(i int) bool { return d.points[i].X <= x }) + 1
	d.points = append(d.points, lttb.Point{})
	copy(d.points[i+1:], d.points[i:])
	d.points[i] = lttb.Point{x, v}
}

func (d *rollingDownsampler) rollOut() {
	last := d.points[len(d.points)-1].X
	start := LinearSearch(len(d.points), func(i int) bool { return (last - d.points[i].X) <= d.window })
	d.points = d.points[start:]
}

func (d *rollingDownsampler) Add(t time.Duration, v float64) {
	x := t.Seconds()
	if d.outdated(x) == true {
		// we may get back outdated data, we can skip it
		return
	}
	d.insertionSort(x, v)
	d.rollOut()
}

func (d *rollingDownsampler) TimeSerie() []lttb.Point {
	if len(d.points) >= d.samplesThreshold {
		d.points = lttb.LTTB(d.points, d.samplesThreshold)
	}
	return append([]lttb.Point(nil), d.points...)
}
