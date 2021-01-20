package main

import (
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ZoneNotFoundError string

func NewZoneNotFoundError(fullname string) ZoneNotFoundError {
	return ZoneNotFoundError("hermes: unknwon zone '" + fullname + "'")
}

func (z ZoneNotFoundError) Error() string {
	return string(z)
}

type ZoneData struct {
	zone RegisteredZone

	climate  ClimateReportManager
	alarmMap map[string]int
}

type Olympus struct {
	mutex *sync.RWMutex
	zones map[string]*ZoneData
	log   *log.Logger
}

func buildRegisteredAlarm(ae *zeus.AlarmEvent) RegisteredAlarm {
	res := RegisteredAlarm{
		Reason:     ae.Reason,
		Level:      zeus.MapPriority(ae.Priority),
		LastChange: &time.Time{},
		Triggers:   0,
		On:         false,
	}
	*res.LastChange = ae.Time
	return res
}

func (z *ZoneData) registerAlarm(ae *zeus.AlarmEvent) {
	if _, ok := z.alarmMap[ae.Reason]; ok == true {
		return
	}

	z.alarmMap[ae.Reason] = len(z.zone.Alarms)

	z.zone.Alarms = append(z.zone.Alarms, buildRegisteredAlarm(ae))
}

func (h *Olympus) RegisterZone(reg *zeus.ZoneRegistration, unused *int) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if z, ok := h.zones[reg.Fullname()]; ok == true {
		//close everything
		close(z.climate.Inbound())
		delete(h.zones, reg.Fullname())
	}
	h.log.Printf("Registering %s", reg.Fullname())
	res := &ZoneData{
		zone: RegisteredZone{
			Host:        reg.Host,
			Name:        reg.Name,
			Temperature: 0.0,
			TemperatureBounds: Bounds{
				nil, nil,
			},
			Humidity: 0.0,
			HumidityBounds: Bounds{
				nil, nil,
			},
		},
		climate:  NewClimateReportManager(),
		alarmMap: make(map[string]int),
	}
	go func() {
		res.climate.Sample()
	}()

	h.zones[reg.Fullname()] = res

	return nil
}

func (h *Olympus) ZoneIsRegistered(reg *zeus.ZoneUnregistration, ok *bool) error {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, *ok = h.zones[reg.Fullname()]
	return nil
}

func (h *Olympus) UnregisterZone(reg *zeus.ZoneUnregistration, unused *int) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	z, ok := h.zones[reg.Fullname()]
	if ok == false {
		return ZoneNotFoundError(reg.Fullname())
	}
	h.log.Printf("Unregistering  %s", reg.Fullname())
	//it will close Sample go routine
	close(z.climate.Inbound())
	delete(h.zones, reg.Fullname())

	return nil
}

func (h *Olympus) ReportClimate(cr *zeus.NamedClimateReport, unused *int) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	z, ok := h.zones[cr.ZoneIdentifier]
	if ok == false {
		return ZoneNotFoundError(cr.ZoneIdentifier)
	}

	z.zone.Temperature = float64((*cr).Temperatures[0])
	z.zone.Humidity = float64((*cr).Humidity)
	z.climate.Inbound() <- zeus.ClimateReport{
		Time:         cr.Time,
		Humidity:     cr.Humidity,
		Temperatures: [4]zeus.Temperature{cr.Temperatures[0], cr.Temperatures[1], cr.Temperatures[2], cr.Temperatures[3]},
	}

	return nil
}

func (h *Olympus) ReportAlarm(ae *zeus.AlarmEvent, unused *int) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	z, ok := h.zones[ae.Zone]
	if ok == false {
		return ZoneNotFoundError(ae.Zone)
	}
	aIdx, ok := z.alarmMap[ae.Reason]
	if ok == false {
		z.registerAlarm(ae)
		aIdx = z.alarmMap[ae.Reason]
	}

	h.log.Printf("New alarm event %+v", ae)
	if ae.Status == zeus.AlarmOn {
		if z.zone.Alarms[aIdx].On == false {
			z.zone.Alarms[aIdx].Triggers += 1
		}
		z.zone.Alarms[aIdx].On = true
	} else {
		z.zone.Alarms[aIdx].On = false
	}
	z.zone.Alarms[aIdx].LastChange = &time.Time{}
	*z.zone.Alarms[aIdx].LastChange = ae.Time
	//TODO: notify

	return nil
}

func (h *Olympus) ReportState(sr *zeus.StateReport, unused *int) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	z, ok := h.zones[sr.Zone]
	if ok == false {
		return ZoneNotFoundError(sr.Zone)
	}

	if z.zone.Current == nil {
		z.zone.Current = &zeus.State{}
	}
	*z.zone.Current = sr.Current
	z.zone.CurrentEnd = sr.CurrentEnd
	if sr.Next != nil && sr.NextTime != nil {
		z.zone.Next = sr.Next
		z.zone.NextTime = sr.NextTime
		z.zone.NextEnd = sr.NextEnd
	} else {
		z.zone.Next = nil
		z.zone.NextEnd = nil
		z.zone.NextTime = nil
	}

	return nil
}

func (h *Olympus) getZones() []RegisteredZone {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	res := make([]RegisteredZone, 0, len(h.zones))
	for _, z := range h.zones {
		toAppend := RegisteredZone{
			Host: z.zone.Host,
			Name: z.zone.Name,
		}
		res = append(res, toAppend)
	}

	return res
}

func (h *Olympus) getZone(host, name string) (*RegisteredZone, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	fname := path.Join(host, "zone", name)
	z, ok := h.zones[fname]
	if ok == false {
		return nil, NewZoneNotFoundError(fname)
	}
	res := &RegisteredZone{}
	*res = z.zone
	return res, nil
}

func (h *Olympus) getClimateReport(host, name, window string) (ClimateReportTimeSerie, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	fname := path.Join(host, "zone", name)
	z, ok := h.zones[fname]
	if ok == false {
		return ClimateReportTimeSerie{}, NewZoneNotFoundError(fname)
	}

	switch window {
	case "hour":
		return z.climate.LastHour(), nil
	case "day":
		return z.climate.LastDay(), nil
	case "week":
		return z.climate.LastWeek(), nil
	case "ten-minutes":
		return z.climate.LastTenMinutes(), nil
	default:
		return z.climate.LastHour(), nil
	}

}

func NewOlympus() *Olympus {
	return &Olympus{
		mutex: &sync.RWMutex{},
		zones: make(map[string]*ZoneData),
		log:   log.New(os.Stderr, "[rpc] :", log.LstdFlags),
	}
}
