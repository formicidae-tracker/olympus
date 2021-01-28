package main

import . "gopkg.in/check.v1"

type UtilsSuite struct {
}

var _ = Suite(&UtilsSuite{})

func (s *UtilsSuite) TestBackLinearSearch(c *C) {
	testdata := []struct {
		Slice    []int
		Value    int
		Expected int
	}{
		{nil, 2, -1},
		{[]int{3, 4, 5}, 2, -1},
		{[]int{1, 3, 4}, 2, 0},
		{[]int{1, 2, 3, 4}, 2, 1},
		{[]int{1, 2, 3, 4}, 5, 3},
	}

	for _, d := range testdata {
		r := BackLinearSearch(len(d.Slice), func(i int) bool { return d.Slice[i] <= d.Value })
		c.Check(r, Equals, d.Expected, Commentf("Searching <=%d in %v", d.Value, d.Slice))
	}

}

func (s *UtilsSuite) TestLinearSearch(c *C) {
	testdata := []struct {
		Slice    []int
		Value    int
		Expected int
	}{
		{nil, 2, 0},
		{[]int{3, 4, 5}, 2, 3},
		{[]int{1, 3, 4}, 2, 0},
		{[]int{3, 2, 3, 4}, 2, 1},
		{[]int{1, 2, 3, 4}, 0, 4},
	}

	for _, d := range testdata {
		r := LinearSearch(len(d.Slice), func(i int) bool { return d.Slice[i] <= d.Value })
		c.Check(r, Equals, d.Expected, Commentf("Searching <=%d in %v", d.Value, d.Slice))
	}

}
