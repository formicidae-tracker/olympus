package main

import (
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/zeus"
)

type AlarmEvent struct {
	Reason string
	Level  int
	On     bool
	Time   time.Time
}

type ServiceEvent struct {
	Identifier string
	Time       time.Time
	On         bool
	Graceful   bool
}

type Bounds struct {
	Min *float64
	Max *float64
}

type ClimateReportTimeSerie struct {
	Humidity       []lttb.Point
	TemperatureAnt []lttb.Point
	TemperatureAux [][]lttb.Point
}

type ZoneClimateStatus struct {
	Temperature       float64
	Humidity          float64
	TemperatureBounds Bounds
	HumidityBounds    Bounds
	ActiveWarnings    int
	ActiveEmergencies int
}

type ZoneReportSummary struct {
	Host string
	Name string

	Climate *ZoneClimateStatus

	StreamURL string
}

type ZoneClimateReport struct {
	ZoneClimateStatus
	NumAux     int
	Current    *zeus.State
	CurrentEnd *zeus.State
	Next       *zeus.State
	NextEnd    *zeus.State
	NextTime   *time.Time
}

type LetoTrackingRegister struct {
	Host, URL string
}
