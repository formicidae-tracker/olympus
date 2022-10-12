package main

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/atuleu/go-lttb"
	"github.com/formicidae-tracker/olympus/olympuspb"
)

type AlarmEvent struct {
	On   bool
	Time time.Time
}

type AlarmReport struct {
	Reason string
	Level  int
	Events []AlarmEvent
}

type ServiceEvent struct {
	Identifier string
	Time       time.Time
	On         bool
	Graceful   bool
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
	Unit           string
	Humidity       PointSeries
	TemperatureAnt PointSeries
	TemperatureAux []PointSeries
}

type Bounds struct {
	Min *float32
	Max *float32
}

type ZoneClimateReport struct {
	Temperature       *float32
	Humidity          *float32
	TemperatureBounds Bounds
	HumidityBounds    Bounds
	ActiveWarnings    int
	ActiveEmergencies int
	Current           *olympuspb.ClimateState
	CurrentEnd        *olympuspb.ClimateState
	Next              *olympuspb.ClimateState
	NextEnd           *olympuspb.ClimateState
	NextTime          *time.Time
}

type ZoneReportSummary struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
}

type ZoneReport struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
	Alarms  []AlarmReport
}

type ServiceLogs struct {
	Climates [][]ServiceEvent
	Tracking [][]ServiceEvent
}

type StreamInfo struct {
	ExperimentName string
	StreamURL      string
	ThumbnailURL   string
}
