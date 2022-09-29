package main

import (
	"sync"
	"time"

	"github.com/atuleu/go-lttb"
)

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
	values                TimedValues
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
	d.values.Push(time, values, d.minimumPeriod)
	d.values.RollOutOfWindow(d.window)
	d.computeSeries(mx)
}

func (d *rollingDownsampler) AddBatch(time []time.Time, values [][]*float32, mx *sync.RWMutex) {
	for i, t := range time {
		d.values.Push(t, values[i], d.minimumPeriod)
	}
	d.values.RollOutOfWindow(d.window)
	d.computeSeries(mx)
}

func (d *rollingDownsampler) computeSeries(mx *sync.RWMutex) {
	if mx == nil {
		d.series = d.values.Downsample(d.samples, d.values.times[len(d.values.times)-1])
		return
	}

	go func() {
		mx.RLock()
		values := d.values.DeepCopy()
		mx.RUnlock()

		series := values.Downsample(d.samples, values.times[len(values.times)-1])
		mx.Lock()
		defer mx.Unlock()
		d.series = series
	}()
}

func (d *rollingDownsampler) TimeSerie() [][]lttb.Point[float32] {
	return d.series
}
