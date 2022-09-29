package main

import (
	"fmt"
	"math"
	"time"

	"github.com/atuleu/go-lttb"
	. "gopkg.in/check.v1"
)

type TimedValuesSuite struct {
	benchmarkData TimedValues
}

var _ = Suite(&TimedValuesSuite{})

func (s *TimedValuesSuite) SetUpSuite(c *C) {
	duration := 7 * 24 * time.Hour
	period := 500 * time.Millisecond
	s.benchmarkData.values = [][]float32{nil, nil}
	for d := time.Duration(0); d < duration; d += period {
		s.benchmarkData.times = append(s.benchmarkData.times, time.Time{}.Add(d))
		s.benchmarkData.values[1] = append(s.benchmarkData.values[1], float32(d.Seconds()))
	}
}

func (s *TimedValuesSuite) TestPush(c *C) {
	t := time.Time{}

	testdata := []struct {
		Name           string
		A, B, Expected TimedValues
		window         time.Duration
	}{
		{
			Name: "empty",
		},
		{
			Name: "no additions",
			A: TimedValues{
				times:  []time.Time{t},
				values: [][]float32{nil, {0}},
			},
			Expected: TimedValues{
				times:  []time.Time{t},
				values: [][]float32{nil, {0}},
			},
		},
		{
			Name: "single addition to empty",
			B: TimedValues{
				times:  []time.Time{t.Add(1)},
				values: [][]float32{nil, {1}},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(1)},
				values: [][]float32{nil, {1}},
			},
		},
		{
			Name: "multiple addtions to empty",
			B: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2)},
				values: [][]float32{nil, {1, 2}},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2)},
				values: [][]float32{nil, {1, 2}},
			},
		},
		{
			Name: "single addtion to non-empty",
			A: TimedValues{
				times:  []time.Time{t},
				values: [][]float32{nil, {0}},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(1)},
				values: [][]float32{nil, {1}},
			},
			Expected: TimedValues{
				times:  []time.Time{t, t.Add(1)},
				values: [][]float32{nil, {0, 1}},
			},
		},
		{
			Name: "multiple prepend",
			A: TimedValues{
				times:  []time.Time{t.Add(5)},
				values: [][]float32{{5}},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2), t.Add(3), t.Add(4)},
				values: [][]float32{{1, 2, 3, 4}},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2), t.Add(3), t.Add(4), t.Add(5)},
				values: [][]float32{{1, 2, 3, 4, 5}},
			},
		},
		{
			Name: "multiple insertion",
			A: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(3), t.Add(4), t.Add(5)},
				values: [][]float32{{1, 3, 4, 5}, nil},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(2), t.Add(3), t.Add(4)},
				values: [][]float32{{12, 13, 14}, nil},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2), t.Add(3), t.Add(4), t.Add(5)},
				values: [][]float32{{1, 12, 13, 14, 5}, nil},
			},
		},
		{
			Name:   "single frequency pruning - after",
			window: 2,
			A: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(3)},
				values: [][]float32{{13}},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}},
			},
		},
		{
			Name:   "single frequency pruning - before",
			window: 2,
			A: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(1)},
				values: [][]float32{{11}},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}},
			},
		},
		{
			Name:   "multiple frequency pruning - after",
			window: 2,
			A: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}, nil},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(3), t.Add(4)},
				values: [][]float32{{13, 14}, nil},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}, nil},
			},
		},
		{
			Name:   "multiple frequency pruning - before",
			window: 2,
			A: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}, nil},
			},
			B: TimedValues{
				times:  []time.Time{t.Add(1), t.Add(2)},
				values: [][]float32{{11, 12}, nil},
			},
			Expected: TimedValues{
				times:  []time.Time{t.Add(2)},
				values: [][]float32{{2}, nil},
			},
		},
	}

	for _, d := range testdata {
		res := d.A.DeepCopy()
		res.Push(d.B, d.window)
		c.Check(res, DeepEquals, d.Expected, Commentf("testing %s", d.Name))
	}

}

func appendValue(values *TimedValues, i int) {
	values.times = append(values.times, time.Time{}.Add(time.Duration(i)))
	if values.values == nil {
		values.values = [][]float32{{float32(i)}}
	} else {
		values.values[0] = append(values.values[0], float32(i))
	}
}

func (s *TimedValuesSuite) TestFrequencyCutoff(c *C) {
	var input, expected TimedValues
	for i := 0; i < 100; i++ {
		appendValue(&input, i)
		if i%2 == 0 {
			appendValue(&expected, i)
		}
	}

	times, values := frequencyCutoff(input.times, input.values, 1)
	c.Check(times, DeepEquals, input.times)
	c.Check(values, DeepEquals, input.values)
	times, values = frequencyCutoff(input.times, input.values, 2)
	c.Check(times, DeepEquals, expected.times)
	c.Check(values, DeepEquals, expected.values)
	c.Assert(input.times, HasLen, 100)
	c.Assert(expected.times, HasLen, 50)

}

func (s *TimedValuesSuite) TestRollOutOfWindow(c *C) {
	t := time.Time{}
	data := TimedValues{
		times:  []time.Time{t.Add(1), t.Add(2), t.Add(3), t.Add(4), t.Add(5)},
		values: [][]float32{{1, 2, 3, 4, 5}, nil},
	}

	var empty TimedValues
	empty.RollOutOfWindow(0)
	c.Check(empty, DeepEquals, TimedValues{})

	unchanged := data.DeepCopy()
	unchanged.RollOutOfWindow(100)
	c.Check(unchanged, DeepEquals, data)

	truncated := data.DeepCopy()
	truncated.RollOutOfWindow(3)
	expected := TimedValues{times: data.times[2:], values: [][]float32{data.values[0][2:], nil}}
	c.Check(truncated, DeepEquals, expected)
}

type allPointClose float64

var AllPointClose Checker = allPointClose(1e-12)

func (e allPointClose) Info() *CheckerInfo {
	return &CheckerInfo{Name: "DeepEquals", Params: []string{"obtained", "expected"}}
}

func (epsilon allPointClose) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()

	obtained := params[0].([][]lttb.Point[float32])
	expected := params[1].([][]lttb.Point[float32])
	if len(obtained) != len(expected) {
		return false, fmt.Sprintf("len(%s) == %d and len(%s) == %d mismatch", names[0], len(obtained), names[1], len(expected))
	}
	for i := range obtained {
		if len(obtained[i]) != len(expected[i]) {
			return false, fmt.Sprintf("len(%s[%d]) == %d and len(%s[%d]) == %d mismatch",
				names[0], i, len(obtained[i]),
				names[1], i, len(expected[i]),
			)
		}

		for j := range obtained[i] {
			o := obtained[i][j]
			e := expected[i][j]
			if math.Abs(float64(o.X-e.X)) > float64(epsilon) || math.Abs(float64(o.Y-e.Y)) > float64(epsilon) {
				return false, fmt.Sprintf("Points %s[%d][%d] (X:%f,Y:%f) and %s[%d][%d] (X:%f,Y:%f) are not close",
					names[0], i, j, o.X, o.Y,
					names[1], i, j, e.X, e.Y)
			}
		}
	}

	return true, ""
}

func (s *TimedValuesSuite) TestDownsample(c *C) {
	var data TimedValues
	expected := [][]lttb.Point[float32]{nil}
	for i := 0; i < 100; i++ {
		appendValue(&data, i)
		expected[0] = append(expected[0],
			lttb.Point[float32]{
				X: float32(i) * 1e-9,
				Y: float32(i),
			})
	}

	c.Check(data.Downsample(100, time.Time{}), AllPointClose, expected)
	downsampled := data.Downsample(50, time.Time{})
	c.Assert(downsampled, HasLen, 1)
	c.Check(downsampled[0], HasLen, 50)

	fromLast := data.Downsample(100, data.times[len(data.times)-1])
	c.Assert(fromLast, HasLen, 1)
	for i, p := range fromLast[0] {
		x := float64(p.X)
		y := float64(p.Y)
		expectedTime := float64(i-99) * 1e-9
		expectedY := float64(i)
		comment := Commentf("%d:(%fns,%f) expected:(%fns,%f)", i, p.X*1e9, p.Y, expectedTime*1e9, expectedY)
		c.Check(math.Abs(x-expectedTime) < 1e-14, Equals, true, comment)
		c.Check(math.Abs(y-float64(i)) < 1e-14, Equals, true, comment)
	}

	empty := (&TimedValues{}).Downsample(100, time.Time{})
	c.Check(empty, IsNil)
}

func (s *TimedValuesSuite) BenchmarkInsertionTime(c *C) {
	minimumPeriod := 3 * time.Minute
	nbSamples := len(s.benchmarkData.times)
	pageSize := 500
	for i := 0; i < c.N; i++ {
		var v TimedValues
		for j := 0; j < nbSamples; j += pageSize {
			end := j + pageSize
			if end > nbSamples {
				end = nbSamples
			}
			v.Push(TimedValues{
				times:  s.benchmarkData.times[j:end],
				values: [][]float32{nil, s.benchmarkData.values[1][j:end]},
			},
				minimumPeriod)
		}
	}
}
