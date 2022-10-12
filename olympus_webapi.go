package main

import (
	"encoding/json"
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
	var res []byte
	for _, p := range series {
		asJson, err := json.Marshal(Point(p))
		if err != nil {
			return nil, err
		}
		res = append(res, asJson...)
	}
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
