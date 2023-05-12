package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		"Bounds.json": {
			api.Bounds{},
			api.Bounds{Minimum: newWithValue[float32](1.0), Maximum: newWithValue[float32](2.0)},
		},
		"ClimateState.json": {
			api.ClimateState{},
			api.ClimateState{Name: "day",
				Temperature:  newWithValue[float32](18.0),
				Humidity:     newWithValue[float32](60.0),
				Wind:         newWithValue[float32](100.0),
				VisibleLight: newWithValue[float32](80.0),
				UvLight:      newWithValue[float32](10.0),
			},
		},
		"ZoneClimateReport.json": {
			api.ZoneClimateReport{},
			api.ZoneClimateReport{
				Temperature:       newWithValue[float32](18.0),
				Humidity:          newWithValue[float32](60.0),
				TemperatureBounds: &api.Bounds{},
				HumidityBounds:    &api.Bounds{},
				ActiveWarnings:    1,
				ActiveEmergencies: 2,
				Current:           &api.ClimateState{},
				CurrentEnd:        &api.ClimateState{},
				Next:              &api.ClimateState{},
				NextEnd:           &api.ClimateState{},
				NextTime:          timestamppb.New(time.Unix(0, 0)),
			},
		},
		"StreamInfo.json": {
			api.StreamInfo{},
			api.StreamInfo{
				ExperimentName: "foo",
				Stream_URL:     "/olympus/hls/somehost.m3u",
				Thumbnail_URL:  "/olympus/somehost.png",
			},
		},
		"TrackingInfo.json": {
			api.TrackingInfo{},
			api.TrackingInfo{
				TotalBytes:     1000*1024 ^ 2,
				FreeBytes:      800*1024 ^ 2,
				BytesPerSecond: 10*1024 ^ 2,
				Stream:         &api.StreamInfo{},
			},
		},

		"ZoneReportSummary.json": {
			api.ZoneReportSummary{},
			api.ZoneReportSummary{
				Host:     "somehost",
				Name:     "box",
				Climate:  &api.ZoneClimateReport{},
				Tracking: &api.TrackingInfo{},
			},
		},
	}
}

var BASEDIR string = "./webapp/src/app/olympus-api/examples"

func writeFileAsJSON(filename string, data []interface{}) error {
	filepath := path.Join(BASEDIR, filename)
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
	if err := os.MkdirAll(BASEDIR, 0755); err != nil {
		return err
	}
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
