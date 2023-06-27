package olympus

import (
	. "gopkg.in/check.v1"
)

type UtilsSuite struct {
}

var _ = Suite(&UtilsSuite{})

var less = func(i, j int) bool { return i < j }

func (s *UtilsSuite) TestBackLinearSearch(c *C) {
	testdata := []struct {
		Slice          []int
		Value          int
		ExpectedIndice int
		ExpectedSlice  []int
	}{
		{nil, 2, 0, []int{2}},
		{[]int{3, 4, 5}, 2, 0, []int{2, 3, 4, 5}},
		{[]int{1, 3, 4}, 2, 1, []int{1, 2, 3, 4}},
		{[]int{1, 2, 3, 4}, 2, 2, []int{1, 2, 2, 3, 4}},
		{[]int{1, 2, 3, 4}, 5, 4, []int{1, 2, 3, 4, 5}},
	}

	for _, d := range testdata {
		r := BackLinearSearch(d.Slice, d.Value, less)
		c.Check(r, Equals, d.ExpectedIndice,
			Commentf("Searching index <%d in %v", d.Value, d.Slice))

		sorted := BackInsertionSort(d.Slice, d.Value, less)
		c.Check(sorted, DeepEquals, d.ExpectedSlice,
			Commentf("Inserting %d in %v", d.Value, d.Slice))

	}

}

func (s *UtilsSuite) TestLinearSearch(c *C) {
	testdata := []struct {
		Slice    []int
		Value    int
		Expected int
	}{
		{nil, 2, 0},
		{[]int{3, 4, 5}, 2, 0},
		{[]int{1, 3, 4}, 2, 1},
		{[]int{3, 2, 3, 4}, 2, 0},
		{[]int{1, 2, 3, 4}, 0, 0},
	}

	for _, d := range testdata {
		r := LinearSearch(d.Slice, d.Value, less)
		c.Check(r, Equals, d.Expected, Commentf("Searching <%d in %v", d.Value, d.Slice))
	}

}

func (s *UtilsSuite) TestInsertSlice(c *C) {
	testdata := []struct {
		A, B, Expected []int
		Begin, End     int
	}{
		{A: nil, B: nil, Expected: nil, Begin: 0, End: 0},
		{A: nil, B: []int{1, 2, 3}, Expected: []int{1, 2, 3}, Begin: 0, End: 0},
		{A: []int{1, 5}, B: []int{2, 3, 4}, Expected: []int{1, 2, 3, 4, 5}, Begin: 1, End: 1},
		{A: []int{1, 4, 5}, B: []int{2, 3, 4}, Expected: []int{1, 2, 3, 4, 5}, Begin: 1, End: 2},
		{A: []int{1, 2, 3}, B: []int{4, 5}, Expected: []int{1, 2, 3, 4, 5}, Begin: 3, End: 3},
		{A: []int{5}, B: []int{1, 2, 3, 4}, Expected: []int{1, 2, 3, 4, 5}, Begin: 0, End: 0},
	}
	for _, d := range testdata {
		if len(d.B) > 0 {
			if c.Check(BackLinearSearch(d.A, d.B[0], less), Equals, d.Begin, Commentf("Malformed test %+v", d)) == false {
				continue
			}

			if c.Check(BackLinearSearch(d.A, d.B[len(d.B)-1], less), Equals, d.End, Commentf("Malformed test %+v", d)) == false {
				continue
			}
		}
		res := InsertSlice(d.A, d.B, d.Begin, d.End)
		c.Check(res, DeepEquals, d.Expected,
			Commentf("Inserting %v in %v at [%d,%d]", d.B, d.A, d.Begin, d.End))
	}
}

func (s *UtilsSuite) TestZoneIdentifier(c *C) {
	testdata := []struct {
		Host, Zone, Expected string
	}{
		{"", "", "."},
		{"foo", "", "foo."},
		{"", "bar", ".bar"},
		{"foo", "bar", "foo.bar"},
	}

	for _, d := range testdata {
		c.Check(ZoneIdentifier(d.Host, d.Zone), Equals, d.Expected)
	}

}
