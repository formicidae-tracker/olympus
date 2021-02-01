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

type ServiceLogs struct {
	Climates [][]ServiceEvent
	Tracking [][]ServiceEvent
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

type StreamInfo struct {
	StreamURL    string
	ThumbnailURL string
}

type ZoneReportSummary struct {
	Host string
	Name string

	Climate *ZoneClimateStatus

	Stream *StreamInfo
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

type ZoneReport struct {
	Host    string
	Name    string
	Climate *ZoneClimateReport
	Stream  *StreamInfo
	Alarms  []AlarmEvent
}

type LetoTrackingRegister struct {
	Host, URL string
}

func (r *ZoneClimateReport) makeCopy() *ZoneClimateReport {
	res := &ZoneClimateReport{
		ZoneClimateStatus: r.ZoneClimateStatus,
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
