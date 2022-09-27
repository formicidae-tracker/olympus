package main

import (
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/olympus/proto"
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

type ClimateTimeSerie struct {
	Humidity       []lttb.Point
	TemperatureAnt []lttb.Point
	TemperatureAux [][]lttb.Point
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
	NumAux            int
	Current           *proto.ClimateState
	CurrentEnd        *proto.ClimateState
	Next              *proto.ClimateState
	NextEnd           *proto.ClimateState
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
