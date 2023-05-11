package main

import (
	"time"

	"github.com/formicidae-tracker/olympus/api"
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
	TimeSeries() api.ClimateTimeSeries
}

type climateDataDownsampler struct {
	window, minimumPeriod, unit time.Duration
	samples                     int
	values                      TimedValues
	series                      api.ClimateTimeSeries
}

var supportedUnits = map[time.Duration]string{
	time.Second:    "s",
	time.Minute:    "m",
	time.Hour:      "h",
	24 * time.Hour: "d",
}

func NewClimateDataDownsampler(window, unit time.Duration, samples int) ClimateDataDownsampler {
	targetPeriod := window / time.Duration(samples)
	minimumPeriod := time.Duration(float64(targetPeriod) / CutOfFrequencyRatio)

	if _, ok := supportedUnits[unit]; ok == false {
		unit = time.Minute
	}

	res := &climateDataDownsampler{
		window:        window,
		unit:          unit,
		samples:       samples,
		minimumPeriod: minimumPeriod,
	}
	return res
}

func (d *climateDataDownsampler) Add(values TimedValues) {
	if d.values.Push(values, d.minimumPeriod) == false {
		return
	}
	d.values.RollOutOfWindow(d.window)
	d.computeSeries()
}

func (d *climateDataDownsampler) computeSeries() {
	reference := d.values.times[len(d.values.times)-1]
	series := d.values.Downsample(d.samples, reference, d.unit)
	d.series = api.ClimateTimeSeries{
		Reference: reference,
		Units:     supportedUnits[d.unit],
	}

	if len(series) > 0 {
		d.series.Humidity = series[0]
	}
	if len(series) > 1 {
		d.series.TemperatureAnt = series[1]
	}

	if len(series) > 2 {
		d.series.TemperatureAux = make([]api.PointSeries, len(series)-2)
		for i := range d.series.TemperatureAux {
			d.series.TemperatureAux[i] = series[i+2]
		}
	}
}

func (d *climateDataDownsampler) TimeSeries() api.ClimateTimeSeries {
	return d.series
}
