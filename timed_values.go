package main

import (
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/barkimedes/go-deepcopy"
)

// TimedValues holds a list of time series with a common time
// vector. Values should be added or removed from the time series only
// using this struct interface. It is designed and optimized for live
// data logging, where data is added in real-time as it becomes
// available, but is robust to out-of-order insertion.
type TimedValues struct {
	// A common time vector for all time series, which is always sorted.
	times []time.Time
	// List of points for each time. Please note that Values can
	// contain empty series. However non-empty series always contains
	// the same number of value that in time.
	values [][]float32
}

// Push adds values for a single time. values may not be added if it
// would make the series instantaneous frequency exceed the frequency
// 1/minimumPeriod. This method is to be preffered to add a single
// point in real-time. To add a chunk of data, please consider
// pushBatch().
func (d *TimedValues) Push(t time.Time, values []*float32, minimumPeriod time.Duration) {
	index := BackLinearSearch[time.Time](d.times, t,
		func(a, b time.Time) bool { return a.Before(b) })

	if index > 0 && t.Sub(d.times[index-1]) < minimumPeriod {
		return
	}

	if index < len(d.times) && d.times[index].Sub(t) < minimumPeriod {
		return
	}

	d.times = Insert[time.Time](d.times, t, index)
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

// RollOutOfWindow prunes values which are older than window from the
// last time in the series.
func (d *TimedValues) RollOutOfWindow(window time.Duration) {
	if len(d.times) == 0 {
		return
	}
	minTime := d.times[len(d.times)-1].Add(-window)
	index := LinearSearch(d.times, minTime, func(a, b time.Time) bool { return a.Before(b) })
	d.times = d.times[index:]
	for i, values := range d.values {
		if values == nil {
			continue
		}
		d.values[i] = values[index:]
	}
}

// TimeVector produces a time vector for a series, from a reference
// point. I.e from the first registered time or the last registered
// time.
func (d *TimedValues) TimeVector(reference time.Time) []float32 {
	if len(d.times) == 0 {
		return nil
	}
	res := make([]float32, len(d.times))
	for i, t := range d.times {
		res[i] = float32(t.Sub(reference).Seconds())
	}
	return res
}

// Downsample downsamples the TimedValues to not exceed a given number
// of samples. The list of resulting time series may contains empty
// values. After downsampling, the time vector of the series will most
// likely not correspong to each other, as the LTTB algorithm retains
// time reference to keep the most representative shape of the
// original serie.
func (d *TimedValues) Downsample(samples int, reference time.Time) [][]lttb.Point[float32] {
	times := d.TimeVector(reference)
	if len(times) == 0 {
		return nil
	}
	res := make([][]lttb.Point[float32], len(d.values))
	for i, values := range d.values {
		if values == nil || len(values) != len(times) {
			continue
		}
		vector := make([]lttb.Point[float32], len(times))
		for j, t := range times {
			vector[j] = lttb.Point[float32]{X: t, Y: values[j]}
		}
		res[i] = lttb.LTTB(vector, samples)
	}
	return res
}

// DeepCopy performs a deep copy of the TimedValues.
func (d *TimedValues) DeepCopy() TimedValues {
	return TimedValues{
		times:  deepcopy.MustAnything(d.times).([]time.Time),
		values: deepcopy.MustAnything(d.values).([][]float32),
	}
}

// FrequencyCutOff creates a new TimedValues from the original one,
// which frequency does not exceed 1/minimumPeriod.
func (d *TimedValues) FrequencyCutoff(minimumPeriod time.Duration) TimedValues {
	times, values := frequencyCutoff[float32](d.times, d.values, minimumPeriod)
	return TimedValues{times: times, values: values}
}

func frequencyCutoff[T any](times []time.Time, values [][]T, minimumPeriod time.Duration) ([]time.Time, [][]T) {
	if len(times) == 0 {
		return nil, nil
	}

	copyValues := func(dst, src int) {
		times[dst] = times[src]
		for _, v := range values {
			if v == nil {
				continue
			}
			v[dst] = v[src]
		}
	}

	newSize := 1
	skipped := false
	for i := 1; i < len(times); i++ {
		if times[i].Sub(times[newSize-1]) >= minimumPeriod {
			if skipped == true {
				copyValues(newSize, i)
				skipped = false
			}
			newSize += 1
		} else {
			skipped = true
		}
	}
	times = times[:newSize]
	for _, v := range values {
		if v == nil {
			continue
		}
		v = v[:newSize]
	}
	return times, values
}
