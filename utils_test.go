package main

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
