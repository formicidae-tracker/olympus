package api

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/atuleu/go-lttb"
)

//go:generate go run generate_client_utils.go
//go:generate go fmt client_utils.pb.go
//go:generate protoc --experimental_allow_proto3_optional  --go_out=. --go-grpc_out=. ./olympus_service.proto
//go:generate go run ./examples/generate.go

type AlarmEvent struct {
	Start time.Time  `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type AlarmReport struct {
	Identification string       `json:"identification,omitempty"`
	Level          AlarmLevel   `json:"level"`
	Events         []AlarmEvent `json:"events"`
	Description    string       `json:"description"`
}

func (r *AlarmReport) On() bool {
	if len(r.Events) == 0 {
		return false
	}
	return r.Events[len(r.Events)-1].End != nil
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
			res = append(res, `[`...)
		} else {
			res = append(res, `,[`...)
		}
		X := float64(p.X)
		Y := float64(p.Y)
		if math.IsNaN(X) || math.IsNaN(Y) {
			return nil, errors.New("json: unsupported value: NaN")
		}
		res = strconv.AppendFloat(res, X, 'g', 5, 32)
		res = append(res, `,`...)
		res = strconv.AppendFloat(res, Y, 'g', 5, 32)
		res = append(res, `]`...)
	}
	res = append(res, `]`...)
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
	Since             time.Time     `json:"since,omitempty"`
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
	Since          time.Time   `json:"since,omitempty"`
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
	Start    time.Time  `json:"start,omitempty"`
	End      *time.Time `json:"end,omitempty"`
	Graceful bool       `json:"graceful"`
}

type ServiceLog struct {
	Zone   string          `json:"zone,omitempty"`
	Events []*ServiceEvent `json:"events,omitempty"`
}

func (l *ServiceLog) On() bool {
	if len(l.Events) == 0 {
		return false
	}
	lastEvent := l.Events[len(l.Events)-1]
	return lastEvent.End == nil
}

func (l *ServiceLog) SetOn(time time.Time) {
	if l.On() == true {
		return
	}
	l.Events = append(l.Events, &ServiceEvent{Start: time})
}

func (l *ServiceLog) SetOff(t time.Time, graceful bool) {
	if l.On() == false {
		return
	}
	lastEvent := l.Events[len(l.Events)-1]
	lastEvent.End = new(time.Time)
	*lastEvent.End = t
	lastEvent.Graceful = graceful
}

type ZoneReport struct {
	Host     string             `json:"host,omitempty"`
	Name     string             `json:"name,omitempty"`
	Climate  *ZoneClimateReport `json:"climate,omitempty"`
	Tracking *TrackingInfo      `json:"tracking,omitempty"`
	Alarms   []AlarmReport      `json:"alarms,omitempty"`
}

type NotificationSettings struct {
	NotifyOnWarning   bool     `json:"notifyOnWarning,omitempty"`
	NotifyNonGraceful bool     `json:"notifyNonGraceful,omitempty"`
	SubscribeToAll    bool     `json:"subscribeToAll,omitempty"`
	Subscriptions     []string `json:"subscriptions,omitempty"`
}

func (s NotificationSettings) SubscribedTo(zone string) bool {
	if s.SubscribeToAll == true {
		return true
	}
	for _, z := range s.Subscriptions {
		if z == zone {
			return true
		}
	}
	return false
}

type NotificationSettingsUpdate struct {
	Endpoint string               `json:"endpoint,omitempty"`
	Settings NotificationSettings `json:"settings,omitempty"`
}
