package main

import (
	"encoding/json"
	"math"

	. "gopkg.in/check.v1"
)

type WebAPISuite struct {
}

var _ = Suite(&WebAPISuite{})

func (s *WebAPISuite) TestPointFormatting(c *C) {
	testdata := []struct {
		Input    Point
		Expected string
		Error    string
	}{
		{Point{}, `{"x":0,"y":0}`, ""},
		{Point{1.2, 1.4}, `{"x":1.2,"y":1.4}`, ""},
		{Point{X: float32(math.NaN())}, ``, ".*json: unsupported value: NaN"},
		{Point{Y: float32(math.NaN())}, ``, ".*json: unsupported value: NaN"},
	}

	for _, d := range testdata {
		comment := Commentf("Test fixture %+v", d)
		res, err := json.Marshal(d.Input)
		if len(d.Error) == 0 {
			c.Check(err, IsNil, comment)
			c.Check(string(res), Equals, d.Expected, comment)
		} else {
			c.Check(err, ErrorMatches, d.Error, comment)
			c.Check(res, IsNil, comment)
		}
	}
}
