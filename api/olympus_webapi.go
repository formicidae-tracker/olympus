package api

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/atuleu/go-lttb"
)

type WebAlarmEvent struct {
	On   bool
	Time time.Time
}

type WebAlarmReport struct {
	Reason string
	Level  int
	Events []WebAlarmEvent
}

func (r *WebAlarmReport) On() bool {
	if len(r.Events) == 0 {
		return false
	}
	return r.Events[len(r.Events)-1].On
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
	Units          string
	Reference      time.Time
	Humidity       PointSeries
	TemperatureAnt PointSeries
	TemperatureAux []PointSeries
}

type Bounds struct {
	Min *float32 `json:"min,omitempty"`
	Max *float32 `json:"max,omitempty"`
}

type ZoneClimateReport struct {
	Temperature       *float32      `json:"temperature,omitempty"`
	Humidity          *float32      `json:"humidity,omitempty"`
	TemperatureBounds Bounds        `json:"temperature_bounds"`
	HumidityBounds    Bounds        `json:"humidity_bounds"`
	ActiveWarnings    int           `json:"active_warnings"`
	ActiveEmergencies int           `json:"active_emergencies"`
	Current           *ClimateState `json:"current,omitempty"`
	CurrentEnd        *ClimateState `json:"current_end,omitempty"`
	Next              *ClimateState `json:"next,omitempty"`
	NextEnd           *ClimateState `json:"next_empty,omitempty"`
	NextTime          *time.Time    `json:"next_time,omitempty"`
}

type ZoneReportSummary struct {
	Host    string             `json:"host"`
	Name    string             `json:"name"`
	Climate *ZoneClimateReport `json:"climate,omitempty"`
	Stream  *StreamInfo        `json:"stream,omitempty"`
}

type ZoneReport struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
	Alarms  []WebAlarmReport
}

type ServiceLogs struct {
	Climates [][]ServiceEvent
	Tracking [][]ServiceEvent
}

type StreamInfo struct {
	ExperimentName string `json:"experiment_name"`
	StreamURL      string `json:"stream_URL"`
	ThumbnailURL   string `json:"thumbnail_URL"`
}
