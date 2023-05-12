package api

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/atuleu/go-lttb"
)

func (r *AlarmReport) On() bool {
	if len(r.Events) == 0 {
		return false
	}
	return r.Events[len(r.Events)-1].Status == AlarmStatus_ON
}

type Point lttb.Point[float32]

func (p Point) MarshalJSON() ([]byte, error) {
	x, err := json.Marshal(p.X)
	if err != nil {
		return nil, err
	}
	y, err := json.Marshal(p.Y)
	if err != nil {
		return nil, err
	}

	res := `{"x":` + string(x) + `,"y":` + string(y) + `}`
	return []byte(res), nil
}

type PointSeries []lttb.Point[float32]

func (series PointSeries) MarshalJSON() ([]byte, error) {
	res := []byte{'['}
	for i, p := range series {
		if i == 0 {
			res = append(res, `{"x":`...)
		} else {
			res = append(res, `,{"x":`...)
		}
		X := float64(p.X)
		Y := float64(p.Y)
		if math.IsNaN(X) || math.IsNaN(Y) {
			return nil, errors.New("json: unsupported value: NaN")
		}
		res = strconv.AppendFloat(res, X, 'g', 5, 32)
		res = append(res, `,"y":`...)
		res = strconv.AppendFloat(res, Y, 'g', 5, 32)
		res = append(res, '}')
	}
	res = append(res, ']')
	return res, nil
}

type ClimateTimeSeries struct {
	Units          string        `json:"units,omitempty"`
	Reference      time.Time     `json:"reference,omitempty"`
	Humidity       PointSeries   `json:"humidity,omitempty"`
	Temperature    PointSeries   `json:"temperature,omitempty"`
	TemperatureAux []PointSeries `json:"temperatureAux,omitempty"`
}
