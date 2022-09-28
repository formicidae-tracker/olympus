package main

import (
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/dgryski/go-lttb"
)

type timedValues struct {
	time   []time.Time
	values [][]float32
}

func (d *timedValues) push(t time.Time, values []*float32) {
	index := BackLinearSearch[time.Time](d.time, t,
		func(a, b time.Time) bool { return a.Before(b) })

	d.time = Insert[time.Time](d.time, t, index)
	if d.values == nil {
		d.values = make([][]float32, len(values))
	}

	for i, v := range values {
		if v == nil {
			continue
		}
		d.values[i] = Insert[float32](d.values[i], *v, index)
	}
}

func (d *timedValues) rollOutOfWindow(window time.Duration) {
	if len(d.time) == 0 {
		return
	}
	minTime := d.time[len(d.time)-1].Add(-window)
	index := LinearSearch(d.time, minTime, func(a, b time.Time) bool { return a.Before(b) })
	d.time = d.time[index:]
	for i, values := range d.values {
		if values == nil {
			continue
		}
		d.values[i] = values[index:]
	}
}

// A DataRollingSampler is used to keep trace of a set of time series
// over a certain time window.
type DataRollingSampler interface {
	// Adds a new point to the sampler
	Add(time.Time, []*float32)
	// Returns the resulting time serie
	TimeSerie() [][]lttb.Point
}

type rollingDownsampler struct {
	window  time.Duration
	samples int
	values  timedValues
	series  [][]lttb.Point
	mx      *sync.RWMutex
	caching bool
}

func NewAsyncRollingSampler(window time.Duration, samples int, mx *sync.RWMutex) DataRollingSampler {
	res := &rollingDownsampler{
		window:  window,
		samples: samples,
		mx:      mx,
	}
	return res
}

func NewRollingSampler(window time.Duration, samples int) DataRollingSampler {
	return NewAsyncRollingSampler(window, samples, nil)
}

func (d *rollingDownsampler) Add(time time.Duration, values []*float32) {
	if d.mx != nil {
		d.mx.Lock()
		defer d.mx.Unlock()
	}

	x := duration.Seconds()
	if d.outdated(x) == true {
		// we may get back outdated data, we can skip it
		return
	}
	d.points = BackInsertionSort(d.points,
		lttb.Point{X: x, Y: v},
		func(a, b lttb.Point) bool {
			return a.X < b.X
		})
	d.rollOut()
	if len(d.points) <= d.samplesThreshold {
		d.sampled = d.points
	} else {
		d.sampled = lttb.LTTB(d.points, d.samplesThreshold)
	}
}

func (d *rollingDownsampler) TimeSerie() [][]lttb.Point {
	return deepcopy.MustAnything(d.sampled).([]lttb.Point)
}
