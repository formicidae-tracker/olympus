package main

import (
	"time"
)

var CutOfFrequencyRatio float64 = 15.0

// A ClimateDataDownsample is used to keep trace of a set of time
// series over a certain time window. It optimize system performance
// by greatly reducing the amount of data.
type ClimateDataDownsampler interface {
	// Adds a new list of value to the sampler for a given time.Time
	// t. If mx is non nil, an asynchronous update will be performed.
	Add(values TimedValues)

	// Returns the resulting time series
	TimeSeries() ClimateTimeSeries
}

type climateDataDownsampler struct {
	window, minimumPeriod time.Duration
	samples               int
	values                TimedValues
	series                ClimateTimeSeries
}

func NewClimateDataDownsampler(window time.Duration, samples int) ClimateDataDownsampler {
	targetPeriod := window / time.Duration(samples)
	minimumPeriod := time.Duration(float64(targetPeriod) / CutOfFrequencyRatio)

	res := &climateDataDownsampler{
		window:        window,
		samples:       samples,
		minimumPeriod: minimumPeriod,
	}
	return res
}

func (d *climateDataDownsampler) Add(values TimedValues) {
	d.values.Push(values, d.minimumPeriod)
	d.values.RollOutOfWindow(d.window)
	d.computeSeries()
}

func (d *climateDataDownsampler) computeSeries() {
	series := d.values.Downsample(d.samples, d.values.times[len(d.values.times)-1])
	d.series = ClimateTimeSeries{}

	if len(series) > 0 {
		d.series.Humidity = series[0]
	}
	if len(series) > 1 {
		d.series.TemperatureAnt = series[1]
	}
	if len(series) > 2 {
		d.series.TemperatureAux = series[2:]
	}
}

func (d *climateDataDownsampler) TimeSeries() ClimateTimeSeries {
	return d.series
}
