package main

import (
	"sync"
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/barkimedes/go-deepcopy"
)

type timedValues struct {
	Time   []time.Time
	Values [][]float32
}

func (d *timedValues) push(t time.Time, values []*float32, minimumPeriod time.Duration) {
	index := BackLinearSearch[time.Time](d.Time, t,
		func(a, b time.Time) bool { return a.Before(b) })

	if index > 0 && t.Sub(d.Time[index-1]) < minimumPeriod {
		return
	}

	if index < len(d.Time) && d.Time[index].Sub(t) < minimumPeriod {
		return
	}

	d.Time = Insert[time.Time](d.Time, t, index)
	if d.Values == nil {
		d.Values = make([][]float32, len(values))
	}

	for i, v := range values {
		if v == nil {
			continue
		}
		d.Values[i] = Insert[float32](d.Values[i], *v, index)
	}
}

func (d *timedValues) rollOutOfWindow(window time.Duration) {
	if len(d.Time) == 0 {
		return
	}
	minTime := d.Time[len(d.Time)-1].Add(-window)
	index := LinearSearch(d.Time, minTime, func(a, b time.Time) bool { return a.Before(b) })
	d.Time = d.Time[index:]
	for i, values := range d.Values {
		if values == nil {
			continue
		}
		d.Values[i] = values[index:]
	}
}

func (d *timedValues) timeVector() []float32 {
	if len(d.Time) == 0 {
		return nil
	}
	last := d.Time[len(d.Time)-1]
	res := make([]float32, len(d.Time))
	for i, t := range d.Time {
		res[i] = float32(t.Sub(last).Seconds())
	}
	return res
}

// A DataRollingSampler is used to keep trace of a set of time series
// over a certain time window.
type DataRollingSampler interface {
	// Adds a new list of value to the sampler for a given time.Time
	// t. If mx is non nil, a asynchronous update will be performed.
	Add(t time.Time, values []*float32, mx *sync.RWMutex)

	// Adds a new list of timed value to the sampler. If mx is
	// non-nil, an asynchronous update will be performed.
	AddBatch(t []time.Time, values [][]*float32, mx *sync.RWMutex)

	// Returns the resulting time serie
	TimeSerie() [][]lttb.Point[float32]
}

type rollingDownsampler struct {
	window, minimumPeriod time.Duration
	samples               int
	values                timedValues
	series                [][]lttb.Point[float32]
}

func NewRollingSampler(window time.Duration, samples int, minimumPeriod time.Duration) DataRollingSampler {
	res := &rollingDownsampler{
		window:        window,
		samples:       samples,
		minimumPeriod: minimumPeriod,
	}
	return res
}

func (d *rollingDownsampler) Add(time time.Time, values []*float32, mx *sync.RWMutex) {
	d.values.push(time, values, d.minimumPeriod)
	d.values.rollOutOfWindow(d.window)
	d.computeSeries(mx)
}

func (d *rollingDownsampler) AddBatch(time []time.Time, values [][]*float32, mx *sync.RWMutex) {
	for i, t := range time {
		d.values.push(t, values[i], d.minimumPeriod)
	}
	d.values.rollOutOfWindow(d.window)
	d.computeSeries(mx)
}

func downsample(values timedValues, samples int) [][]lttb.Point[float32] {
	times := values.timeVector()
	if len(times) == 0 {
		return nil
	}
	res := make([][]lttb.Point[float32], len(values.Values))
	for i, v := range values.Values {
		if v == nil || len(v) != len(times) {
			continue
		}
		vector := make([]lttb.Point[float32], len(times))
		for j, t := range times {
			vector[j] = lttb.Point[float32]{X: t, Y: v[j]}
		}
		res[i] = lttb.LTTB(vector, samples)
	}

	return res
}

func (d *rollingDownsampler) computeSeries(mx *sync.RWMutex) {
	if mx == nil {
		d.series = downsample(d.values, d.samples)
		return
	}

	go func() {
		mx.RLock()
		values := deepcopy.MustAnything(d.values).(timedValues)
		mx.RUnlock()

		series := downsample(values, d.samples)
		mx.Lock()
		defer mx.Unlock()
		d.series = series
	}()
}

func (d *rollingDownsampler) TimeSerie() [][]lttb.Point[float32] {
	return d.series
}
