package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
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

func generateUnitTestData() map[string][]interface{} {
	return map[string][]interface{}{
		"unit-testdata/AlarmTimePoint.json": {
			api.AlarmTimePoint{},
			api.AlarmTimePoint{
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
				Events:         []api.AlarmTimePoint{{Time: time.Unix(1, 1), On: true}},
			},
		},
		"unit-testdata/ClimateTimeSeries.json": {
			api.ClimateTimeSeries{},
			api.ClimateTimeSeries{
				Units:       "s",
				Reference:   time.Unix(3, 0),
				Humidity:    api.PointSeries{{X: 0.0, Y: 50.0}},
				Temperature: api.PointSeries{{X: 0.0, Y: 15.0}},
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
				TemperatureBounds: api.Bounds{Minimum: newWithValue[float32](1.0)},
				HumidityBounds:    api.Bounds{Maximum: newWithValue[float32](1.0)},
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
				Host:              "somehost",
				Name:              "box",
				Climate:           &api.ZoneClimateReport{Temperature: newWithValue[float32](18.0)},
				Tracking:          &api.TrackingInfo{TotalBytes: 1000*1024 ^ 2},
				ActiveWarnings:    1,
				ActiveEmergencies: 2,
			},
		},
		"unit-testdata/ServiceEvent.json": {
			api.ServiceEvent{},
			api.ServiceEvent{
				Time:     time.Unix(1, 1),
				On:       true,
				Graceful: true,
			},
		},
		"unit-testdata/ServiceEventList.json": {
			api.ServiceEventList{},
			api.ServiceEventList{
				Zone:   "somehost.box",
				Events: []api.ServiceEvent{{Time: time.Unix(1, 1), On: true, Graceful: false}},
			},
		},
		"unit-testdata/ServicesLogs.json": {
			api.ServicesLogs{},
			api.ServicesLogs{
				Climate:  []api.ServiceEventList{{Zone: "foo"}},
				Tracking: []api.ServiceEventList{{Zone: "foo"}},
			},
		},
		"unit-testdata/ZoneReport.json": {
			api.ZoneReport{},
			api.ZoneReport{
				Host: "foo",
				Name: "bar",
				Climate: &api.ZoneClimateReport{
					Temperature: newWithValue[float32](18.0),
				},
				Tracking: &api.TrackingInfo{
					TotalBytes: 2000,
				},
				Alarms: []api.AlarmReport{{Identification: "error"}},
			},
		},
	}
}

func timeMustParse(value string) time.Time {
	res, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return res
}

func generateAlarmReport(identifier, description string, level api.AlarmLevel, number int) api.AlarmReport {
	on := number%2 == 1
	t := timeMustParse("2023-03-31T23:25:34.000Z")
	points := make([]api.AlarmTimePoint, 0, number)
	for i := 0; i < number; i++ {
		t = t.Add(-1 * time.Duration(rand.Intn(120)) * time.Minute)
		points = append([]api.AlarmTimePoint{{Time: t, On: on}}, points...)
		on = !on
	}

	return api.AlarmReport{
		Identification: identifier,
		Level:          level,
		Description:    description,
		Events:         points,
	}
}

func generateMockData() (map[string]interface{}, map[string]string) {
	minervaClimate := &api.ZoneClimateReport{
		Temperature:       newWithValue[float32](19.8743205432),
		Humidity:          newWithValue[float32](53.465679028734),
		TemperatureBounds: api.Bounds{Minimum: newWithValue[float32](17.0), Maximum: newWithValue[float32](22.0)},
		HumidityBounds:    api.Bounds{Minimum: newWithValue[float32](40.0), Maximum: newWithValue[float32](70.0)},
		Current: &api.ClimateState{
			Name:         "night-to-day",
			Temperature:  newWithValue[float32](19.87645),
			Humidity:     newWithValue[float32](60.0),
			Wind:         newWithValue[float32](100),
			VisibleLight: newWithValue[float32](23.47892348),
			UvLight:      newWithValue[float32](0.348978907),
		},
		CurrentEnd: &api.ClimateState{
			Name:         "night-to-day",
			Temperature:  newWithValue[float32](23.0),
			VisibleLight: newWithValue[float32](100.0),
			UvLight:      newWithValue[float32](5.0),
		},
		Next: &api.ClimateState{
			Name:         "day",
			Temperature:  newWithValue[float32](23.0),
			Humidity:     newWithValue[float32](60.0),
			Wind:         newWithValue[float32](100),
			VisibleLight: newWithValue[float32](100.0),
			UvLight:      newWithValue[float32](5.0),
		},
		NextTime: newWithValue(timeMustParse("2023-04-01T06:30:00.000Z")),
	}

	minervaTracking := &api.TrackingInfo{
		TotalBytes:     int64(2.0 * math.Pow(2, 40)),
		FreeBytes:      int64(0.45123980 * math.Pow(2, 40)),
		BytesPerSecond: int64(2.567879 * math.Pow(2, 20)),
		Stream: &api.StreamInfo{
			ExperimentName: "tackling-universe",
			StreamURL:      "https://moctobpltc-i.akamaihd.net/hls/live/571329/eight/playlist.m3u8",
			ThumbnailURL:   "https://picsum.photos/id/42/1024/776?grayscale",
		},
	}

	minervaAlarms := []api.AlarmReport{
		generateAlarmReport("climate.temperature_out_of_bound",
			"Temperature (22.1°C) is out of allowed range (17.0°C - 22.0°C)",
			api.AlarmLevel_EMERGENCY,
			36),
		generateAlarmReport("climate.humidity_is_unrerachable",
			"Target humidity cannot be reached",
			api.AlarmLevel_WARNING,
			2),
		generateAlarmReport("climate.cannot_read_sensor",
			"Cannot read sensor",
			api.AlarmLevel_EMERGENCY,
			44),
		generateAlarmReport("climate.device_missing(slcan0.Zeus.1)",
			"Device slcan0.Zeus.1 cannot be reached. Climate may not be controlled",
			api.AlarmLevel_EMERGENCY,
			44),
		generateAlarmReport("climate.device_internal_error(slcan0.Zeus.1,0x0021)",
			"Device slcan0.Zeus.1 experienced internal error 0x0021",
			api.AlarmLevel_WARNING,
			56),
		generateAlarmReport("climate.right_fan_is_aging",
			"Right Extraction Fan is not operating as expected.",
			api.AlarmLevel_WARNING,
			16),
	}

	jupyterClimate := &api.ZoneClimateReport{
		Temperature:       newWithValue[float32](28.8743205432),
		Humidity:          newWithValue[float32](53.465679028734),
		TemperatureBounds: api.Bounds{Minimum: newWithValue[float32](17.0), Maximum: newWithValue[float32](22.0)},
		HumidityBounds:    api.Bounds{Minimum: newWithValue[float32](40.0), Maximum: newWithValue[float32](70.0)},
		Current: &api.ClimateState{
			Name:         "onlyone",
			Temperature:  newWithValue[float32](20.0),
			Humidity:     newWithValue[float32](60.0),
			Wind:         newWithValue[float32](100),
			VisibleLight: newWithValue[float32](100.0),
			UvLight:      newWithValue[float32](0.0),
		},
	}

	jupyterAlarms := []api.AlarmReport{
		generateAlarmReport("climate.temperature_out_of_bound",
			"Temperature (28.9°C) is out of allowed range (17.0°C - 22.0°C)",
			api.AlarmLevel_EMERGENCY,
			25),
		generateAlarmReport("climate.cannot_read_sensor",
			"Cannot read sensor",
			api.AlarmLevel_EMERGENCY,
			22),
		generateAlarmReport("climate.device_missing(slcan0.Zeus.1)",
			"Device slcan0.Zeus.1 cannot be reached. Climate may not be controlled",
			api.AlarmLevel_EMERGENCY,
			22),
		generateAlarmReport("climate.device_internal_error(slcan0.Zeus.1,0x0021)",
			"Device slcan0.Zeus.1 experienced internal error 0x0021",
			api.AlarmLevel_WARNING,
			28),
	}

	junoTracking := &api.TrackingInfo{
		TotalBytes:     int64(2.0 * math.Pow(2, 40)),
		FreeBytes:      int64(80.231 * math.Pow(2, 30)),
		BytesPerSecond: int64(3.89691 * math.Pow(2, 20)),
		Stream: &api.StreamInfo{
			ExperimentName: "about to fail",
			StreamURL:      "https://moctobpltc-i.akamaihd.net/hls/live/571329/eight/playlist.m3u8",
			ThumbnailURL:   "https://picsum.photos/id/43/1024/776?grayscale",
		},
	}

	junoAlarms := []api.AlarmReport{
		generateAlarmReport("tracking.criticaly_low_disk_space",
			"Available space on disk is critically low. Tracking will soon stop.",
			api.AlarmLevel_EMERGENCY,
			1),
	}

	data := map[string]interface{}{
		"_api_zones": []api.ZoneReportSummary{
			{
				Host:              "jupyter",
				Name:              "desert",
				Climate:           jupyterClimate,
				ActiveWarnings:    0,
				ActiveEmergencies: 1,
			},
			{
				Host:           "juno",
				Name:           "box",
				Tracking:       junoTracking,
				ActiveWarnings: 1,
			},
			{
				Host:     "minerva",
				Name:     "box",
				Climate:  minervaClimate,
				Tracking: minervaTracking,
			},
			{
				Host:              "jupyter",
				Name:              "desert",
				Climate:           jupyterClimate,
				ActiveWarnings:    0,
				ActiveEmergencies: 1,
			},
			{
				Host:           "juno",
				Name:           "box",
				Tracking:       junoTracking,
				ActiveWarnings: 1,
			},
			{
				Host:     "minerva",
				Name:     "box",
				Climate:  minervaClimate,
				Tracking: minervaTracking,
			},
		},
		"_api_host_jupyter_zone_desert": &api.ZoneReport{
			Host:    "jupyter",
			Name:    "desert",
			Climate: jupyterClimate,
			Alarms:  jupyterAlarms,
		},
		"_api_host_juno_zone_box": &api.ZoneReport{
			Host:     "juno",
			Name:     "box",
			Tracking: junoTracking,
			Alarms:   junoAlarms,
		},
		"_api_host_minerva_zone_box": &api.ZoneReport{
			Host:     "minerva",
			Name:     "box",
			Climate:  minervaClimate,
			Tracking: minervaTracking,
			Alarms:   minervaAlarms,
		},
		"_api_host_jupyter_zone_desert_alarms": jupyterAlarms,
		"_api_host_juno_zone_box_alarms":       junoAlarms,
		"_api_host_minerva_zone_box_alarms":    minervaAlarms,
	}

	routes := map[string]string{}
	for source := range data {
		target := strings.Replace(source, "_", "/", -1)
		routes[target] = "/" + source
	}
	return data, routes
}

var BASEDIR string = "./webapp/src/app/olympus-api/"

func writeFileAsJSON(filename string, data interface{}) error {

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
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeExamples(dataPerFilename map[string][]interface{}) error {
	for filename, data := range dataPerFilename {
		if err := writeFileAsJSON(filename, data); err != nil {
			return err
		}
	}
	mockData, routes := generateMockData()

	if err := writeFileAsJSON("fake-backend/db.json", mockData); err != nil {
		return err
	}
	return writeFileAsJSON("fake-backend/routes.json", routes)
}

func execute() error {
	data := generateUnitTestData()
	return writeExamples(data)
}
