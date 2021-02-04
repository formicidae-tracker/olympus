package main

import (
	"github.com/formicidae-tracker/leto"
	"github.com/formicidae-tracker/zeus"
)

type OlympusRPCWrapper Olympus

func (o *OlympusRPCWrapper) RegisterZone(zr zeus.ZoneRegistration, unused *int) error {
	return (*Olympus)(o).RegisterZone(zr)
}

func (o *OlympusRPCWrapper) UnregisterZone(zr zeus.ZoneUnregistration, unused *int) error {
	return (*Olympus)(o).UnregisterZone(zr)
}

func (o *OlympusRPCWrapper) ReportClimate(cr zeus.NamedClimateReport, unused *int) error {
	return (*Olympus)(o).ReportClimate(cr)
}

func (o *OlympusRPCWrapper) ReportAlarm(ae zeus.AlarmEvent, unused *int) error {
	return (*Olympus)(o).ReportAlarm(ae)
}

func (o *OlympusRPCWrapper) ReportState(sr zeus.StateReport, unused *int) error {
	return (*Olympus)(o).ReportState(sr)
}

func (o *OlympusRPCWrapper) RegisterTracker(args leto.RegisterTrackerArgs, unused *int) error {
	return (*Olympus)(o).RegisterTracker(args)
}

func (o *OlympusRPCWrapper) UnregisterTracker(hostname string, unused *int) error {
	return (*Olympus)(o).UnregisterTracker(hostname)
}

func (o *OlympusRPCWrapper) ZoneIsRegistered(args zeus.ZoneUnregistration, reply *bool) error {
	*reply = (*Olympus)(o).ZoneIsRegistered(args.Host, args.Name)
	return nil
}
