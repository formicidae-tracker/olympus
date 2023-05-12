package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/formicidae-tracker/olympus/api"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("could not generate JSON examples: %s", err)
	}
}

func newWithValue[T any](v T) *T {
	res := new(T)
	*res = v
	return res
}

func generateData() map[string][]interface{} {
	return map[string][]interface{}{
		"unit-testdata/AlarmTimepoint.json": {
			api.AlarmTimepPoint{},
			api.AlarmTimepPoint{
				Time: time.Unix(1, 1),
				On:   true,
			},
		},
		"unit-testdata/AlarmReport.json": {
			api.AlarmReport{},
			api.AlarmReport{
				Identification: "climate.temperature.out-of-bound",
				Description:    "Current Temperature (34.2°C) is outside [15°C,25°C]",
				Level:          api.AlarmLevel_EMERGENCY,
				Events:         []api.AlarmTimepPoint{{time.Unix(1, 1), true}},
			},
		},
		"unit-testdata/Bounds.json": {
			api.Bounds{},
			api.Bounds{Minimum: newWithValue[float32](1.0), Maximum: newWithValue[float32](2.0)},
		},
		"unit-testdata/ClimateState.json": {
			api.ClimateState{},
			api.ClimateState{Name: "day",
				Temperature:  newWithValue[float32](18.0),
				Humidity:     newWithValue[float32](60.0),
				Wind:         newWithValue[float32](100.0),
				VisibleLight: newWithValue[float32](80.0),
				UvLight:      newWithValue[float32](10.0),
			},
		},
		"unit-testdata/ZoneClimateReport.json": {
			api.ZoneClimateReport{},
			api.ZoneClimateReport{
				Temperature:       newWithValue[float32](18.0),
				Humidity:          newWithValue[float32](60.0),
				TemperatureBounds: api.Bounds{Minimum: newWithValue[float32](0.0)},
				HumidityBounds:    api.Bounds{Maximum: newWithValue[float32](0.0)},
				Current:           &api.ClimateState{Name: "day"},
				CurrentEnd:        nil,
				Next:              &api.ClimateState{Name: "day-to-night"},
				NextEnd:           &api.ClimateState{Name: "day-to-night"},
				NextTime:          newWithValue(time.Unix(1, 1)),
			},
		},
		"unit-testdata/StreamInfo.json": {
			api.StreamInfo{},
			api.StreamInfo{
				ExperimentName: "foo",
				StreamURL:      "/olympus/hls/somehost.m3u",
				ThumbnailURL:   "/olympus/somehost.png",
			},
		},
		"unit-testdata/TrackingInfo.json": {
			api.TrackingInfo{},
			api.TrackingInfo{
				TotalBytes:     1000*1024 ^ 2,
				FreeBytes:      800*1024 ^ 2,
				BytesPerSecond: 10*1024 ^ 2,
				Stream:         &api.StreamInfo{ExperimentName: "foo"},
			},
		},
		"unit-testdata/ZoneReportSummary.json": {
			api.ZoneReportSummary{},
			api.ZoneReportSummary{
				Host:     "somehost",
				Name:     "box",
				Climate:  &api.ZoneClimateReport{Temperature: newWithValue[float32](18.0)},
				Tracking: &api.TrackingInfo{TotalBytes: 1000*1024 ^ 2},
			},
		},
		"unit-testdata/ServiceEvent.json": {
			api.ServiceEvent{},
			api.ServiceEvent{
				Time:     time.Unix(1, 1),
				On:       false,
				Graceful: true,
			},
		},
		"unit-testdata/ServeEventLogs.json": {
			api.ServiceEventList{},
			api.ServiceEventList{
				Zone:   "somehost.box",
				Events: []api.ServiceEvent{{Time: time.Unix(1, 1), On: true, Graceful: false}},
			},
		},
		"unit-testdata/ServiceLogs.json": {
			api.ServicesLogs{},
			api.ServicesLogs{
				Climate:  []api.ServiceEventList{{Zone: "foo"}},
				Tracking: []api.ServiceEventList{{Zone: "foo"}},
			},
		},
	}
}

var BASEDIR string = "./webapp/src/app/olympus-api/"

func writeFileAsJSON(filename string, data []interface{}) error {

	filepath := path.Join(BASEDIR, filename)

	dirpath := path.Dir(filepath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	fmt.Printf("Writing '%s'\n", filepath)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(data)
}

func writeExamples(dataPerFilename map[string][]interface{}) error {
	for filename, data := range dataPerFilename {
		if err := writeFileAsJSON(filename, data); err != nil {
			return err
		}
	}
	return nil
}

func execute() error {
	data := generateData()
	return writeExamples(data)
}
