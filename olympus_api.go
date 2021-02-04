package main

import (
	"time"

	"github.com/dgryski/go-lttb"
	"github.com/formicidae-tracker/zeus"
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
	Min *float64
	Max *float64
}

type ZoneClimateReport struct {
	Temperature       float64
	Humidity          float64
	TemperatureBounds Bounds
	HumidityBounds    Bounds
	ActiveWarnings    int
	ActiveEmergencies int
	NumAux            int
	Current           *zeus.State
	CurrentEnd        *zeus.State
	Next              *zeus.State
	NextEnd           *zeus.State
	NextTime          *time.Time
}

type ZoneReport struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
	Alarms  []AlarmReport
}

type ZoneReportSummary struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
}

type ServiceLogs struct {
	Climates [][]ServiceEvent
	Tracking [][]ServiceEvent
}

type StreamInfo struct {
	StreamURL    string
	ThumbnailURL string
}

type LetoTrackingRegister struct {
	Host, URL string
}

func (r *ZoneClimateReport) makeCopy() *ZoneClimateReport {
	res := &ZoneClimateReport{
		Temperature:       r.Temperature,
		Humidity:          r.Humidity,
		TemperatureBounds: r.TemperatureBounds,
		HumidityBounds:    r.HumidityBounds,
		ActiveWarnings:    r.ActiveWarnings,
		ActiveEmergencies: r.ActiveEmergencies,
		NumAux:            r.NumAux,
	}
	if r.Current != nil {
		res.Current = &zeus.State{}
		*res.Current = zeus.SanitizeState(*r.Current)
	}
	if r.CurrentEnd != nil {
		res.CurrentEnd = &zeus.State{}
		*res.CurrentEnd = zeus.SanitizeState(*r.CurrentEnd)
	}
	if r.Next != nil {
		res.Next = &zeus.State{}
		*res.Next = zeus.SanitizeState(*r.Next)
	}
	if r.NextEnd != nil {
		res.NextEnd = &zeus.State{}
		*res.NextEnd = zeus.SanitizeState(*r.NextEnd)
	}
	if r.NextTime != nil {
		res.NextTime = &time.Time{}
		*res.NextTime = *r.NextTime
	}
	return res
}
