package main

import (
	"time"

	"github.com/atuleu/go-lttb"
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

func lessTime(a, b time.Time) bool {
	return a.Before(b)
}

// pushOne adds a list of values for a single point in time. It is
// optimized to add a point at the end of the TimedValues, but is
// robust for any given time.
func (d *TimedValues) pushOne(t time.Time, values [][]float32, minimumPeriod time.Duration) {
	index := BackLinearSearch(d.times, t, lessTime)

	if index > 0 && t.Sub(d.times[index-1]) < minimumPeriod {
		return
	}

	if index < len(d.times) && d.times[index].Sub(t) < minimumPeriod {
		return
	}

	d.times = Insert(d.times, t, index)
	if d.values == nil {
		d.values = make([][]float32, len(values))
	}

	for i, v := range values {
		if len(v) == 0 {
			continue
		}
		d.values[i] = Insert(d.values[i], v[0], index)
	}
}

func (d *TimedValues) pushBatch(times []time.Time, values [][]float32, minimumPeriod time.Duration) {
	if len(times) == 0 {
		return
	}

	// first we ensure we do not want to add too much values.
	times, values = frequencyCutoff(times, values, minimumPeriod)

	// find the insertion indexes
	insertionStart := BackLinearSearch(d.times, times[0], lessTime)
	insertionEnd := BackLinearSearch(d.times, times[len(times)-1], lessTime)

	if insertionStart > 0 && times[0].Sub(d.times[insertionStart-1]) < minimumPeriod {
		times = times[1:]
		for i, v := range values {
			if len(v) == 0 {
				continue
			}
			values[i] = v[1:]
		}
	}
	if len(times) == 0 {
		return
	}

	if insertionEnd < len(d.times) && d.times[insertionEnd].Sub(times[len(times)-1]) < minimumPeriod {
		times = times[:(len(times) - 1)]
		for i, v := range values {
			if len(v) == 0 {
				continue
			}
			values[i] = v[:(len(v) - 1)]
		}

	}
	if len(times) == 0 {
		return
	}

	d.times = InsertSlice(d.times, times, insertionStart, insertionEnd)
	if len(d.values) == 0 {
		d.values = make([][]float32, len(values))
	}
	for i, v := range values {
		if len(v) == 0 || len(d.values[i]) < (len(d.times)-len(times)) {
			continue
		}
		d.values[i] = InsertSlice(d.values[i], v, insertionStart, insertionEnd)
	}
}

// Push adds another TimedValues to the TimedValues. values may be
// dropped if it would make the series instantaneous frequency exceed
// the frequency 1/minimumPeriod. This method is optimized for two
// scenarios:
//   - Adding a single value at the end of the TimedValues.
//   - Adding a large chunk of values to the TimedValues.
func (d *TimedValues) Push(other TimedValues, minimumPeriod time.Duration) {
	if len(other.times) == 1 {
		d.pushOne(other.times[0], other.values, minimumPeriod)
	} else {
		d.pushBatch(other.times, other.values, minimumPeriod)
	}
}

// RollOutOfWindow prunes values which are older than window from the
// last time in the series.
func (d *TimedValues) RollOutOfWindow(window time.Duration) {
	if len(d.times) == 0 {
		return
	}
	minTime := d.times[len(d.times)-1].Add(-window)
	index := LinearSearch(d.times, minTime, lessTime)
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
func (d *TimedValues) TimeVector(reference time.Time, unit time.Duration) []float32 {
	if len(d.times) == 0 {
		return nil
	}
	factor := 1.0 / unit.Seconds()
	res := make([]float32, len(d.times))
	for i, t := range d.times {
		res[i] = float32(t.Sub(reference).Seconds() * factor)
	}
	return res
}

// Downsample downsamples the TimedValues to not exceed a given number
// of samples. The list of resulting time series may contains empty
// values. After downsampling, the time vector of the series will most
// likely not correspong to each other, as the LTTB algorithm retains
// time reference to keep the most representative shape of the
// original serie.
func (d *TimedValues) Downsample(samples int, reference time.Time, unit time.Duration) [][]lttb.Point[float32] {
	times := d.TimeVector(reference, unit)
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
	if len(d.times) == 0 {
		return TimedValues{}
	}
	times := make([]time.Time, len(d.times))
	copy(times, d.times)
	values := make([][]float32, len(d.values))
	for i, v := range d.values {
		if len(v) == 0 {
			continue
		}
		values[i] = make([]float32, len(v))
		copy(values[i], v)
	}
	return TimedValues{
		times:  times,
		values: values,
	}
}

// FrequencyCutOff creates a new TimedValues from the original one,
// which frequency does not exceed 1/minimumPeriod.
func (d *TimedValues) FrequencyCutoff(minimumPeriod time.Duration) TimedValues {
	times, values := frequencyCutoff(d.times, d.values, minimumPeriod)
	return TimedValues{times: times, values: values}
}

func frequencyCutoff(times []time.Time, values [][]float32, minimumPeriod time.Duration) ([]time.Time, [][]float32) {
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
	for i, v := range values {
		if len(v) == 0 {
			continue
		}
		values[i] = v[:newSize]
	}
	return times, values
}
