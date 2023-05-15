package api

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/atuleu/go-lttb"
)

type AlarmTimePoint struct {
	Time time.Time `json:"time,omitempty"`
	On   bool      `json:"on,omitempty"`
}

type AlarmReport struct {
	Identification string           `json:"identification,omitempty"`
	Level          AlarmLevel       `json:"level"`
	Events         []AlarmTimePoint `json:"events"`
	Description    string           `json:"description"`
}

func (r *AlarmReport) On() bool {
	if len(r.Events) == 0 {
		return false
	}
	return r.Events[len(r.Events)-1].On
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

type Bounds struct {
	Minimum *float32 `json:"minimum,omitempty"`
	Maximum *float32 `json:"maximum,omitempty"`
}

type ZoneClimateReport struct {
	Temperature       *float32      `json:"temperature,omitempty"`
	Humidity          *float32      `json:"humidity,omitempty"`
	TemperatureBounds Bounds        `json:"temperature_bounds,omitempty"`
	HumidityBounds    Bounds        `json:"humidity_bounds,omitempty"`
	Current           *ClimateState `json:"current,omitempty,omitempty"`
	CurrentEnd        *ClimateState `json:"current_end,omitempty,omitempty"`
	Next              *ClimateState `json:"next,omitempty"`
	NextEnd           *ClimateState `json:"next_end,omitempty"`
	NextTime          *time.Time    `json:"next_time,omitempty"`
}

type StreamInfo struct {
	ExperimentName string `json:"experiment_name,omitempty"`
	StreamURL      string `json:"stream_URL,omitempty"`
	ThumbnailURL   string `json:"thumbnail_URL,omitempty"`
}

type TrackingInfo struct {
	TotalBytes     int64       `json:"total_bytes,omitempty"`
	FreeBytes      int64       `json:"free_bytes,omitempty"`
	BytesPerSecond int64       `json:"bytes_per_second,omitempty"`
	Stream         *StreamInfo `json:"stream,omitempty"`
}

type ZoneReportSummary struct {
	Host              string             `json:"host,omitempty"`
	Name              string             `json:"name,omitempty"`
	Climate           *ZoneClimateReport `json:"climate,omitempty"`
	Tracking          *TrackingInfo      `json:"tracking,omitempty"`
	ActiveWarnings    int                `json:"active_warnings,omitempty"`
	ActiveEmergencies int                `json:"active_emergencies,omitempty"`
}

type ServiceEvent struct {
	Time     time.Time `json:"time,omitempty"`
	On       bool      `json:"on,omitempty"`
	Graceful bool      `json:"graceful,omitempty"`
}

type ServiceEventList struct {
	Zone   string         `json:"zone,omitempty"`
	Events []ServiceEvent `json:"events,omitempty"`
}

type ServicesLogs struct {
	Climate  []ServiceEventList `json:"climate,omitempty"`
	Tracking []ServiceEventList `json:"tracking,omitempty"`
}

type ZoneReport struct {
	Host     string             `json:"host,omitempty"`
	Name     string             `json:"name,omitempty"`
	Climate  *ZoneClimateReport `json:"climate,omitempty"`
	Tracking *TrackingInfo      `json:"tracking,omitempty"`
	Alarms   []AlarmReport      `json:"alarms,omitempty"`
}
