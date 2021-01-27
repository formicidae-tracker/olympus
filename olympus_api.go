package main

import (
	"time"

	"github.com/formicidae-tracker/libarke/src-go/arke"
	"github.com/formicidae-tracker/zeus"
)

type RegisteredAlarm struct {
	Reason     string
	On         bool
	Level      int
	LastChange *time.Time
	Triggers   int
}

type Bounds struct {
	Min *float64
	Max *float64
}

type RegisteredZone struct {
	Host              string
	Name              string
	NumAux            int
	Temperature       float64
	TemperatureBounds Bounds
	Humidity          float64
	HumidityBounds    Bounds
	Alarms            []RegisteredAlarm
	Current           *zeus.State
	CurrentEnd        *zeus.State
	Next              *zeus.State
	NextEnd           *zeus.State
	NextTime          *time.Time
}

var stubZone RegisteredZone

func init() {
	stubZone = RegisteredZone{
		Temperature: 21.4,
		TemperatureBounds: Bounds{
			Min: new(float64),
			Max: new(float64),
		},
		Humidity: 62.0,
		HumidityBounds: Bounds{
			Min: new(float64),
			Max: new(float64),
		},
		Alarms: []RegisteredAlarm{},
	}
	*stubZone.TemperatureBounds.Min = 22.0
	*stubZone.TemperatureBounds.Max = 28.0
	*stubZone.HumidityBounds.Min = 40.0
	*stubZone.HumidityBounds.Max = 75.0

	alarms := []zeus.Alarm{
		zeus.WaterLevelWarning,
		zeus.WaterLevelCritical,
		zeus.TemperatureOutOfBound,
		zeus.HumidityOutOfBound,
		zeus.TemperatureUnreachable,
		zeus.HumidityUnreachable,
		zeus.NewMissingDeviceAlarm("slcan0", arke.ZeusClass, 1),
		zeus.NewMissingDeviceAlarm("slcan0", arke.CelaenoClass, 1),
		zeus.NewMissingDeviceAlarm("slcan0", arke.HeliosClass, 1),
	}
	for _, a := range alarms {
		aa := RegisteredAlarm{
			Reason:   a.Reason(),
			On:       false,
			Triggers: 0,
		}
		if a.Flags()&zeus.Emergency != 0 {
			aa.Level = 2
		} else {
			aa.Level = 1
		}

		stubZone.Alarms = append(stubZone.Alarms, aa)
	}

}
