package api

import (
	"github.com/mitchellh/copystructure"
)

func (a *AlarmUpdate) Clone() *AlarmUpdate {
	return copystructure.Must(copystructure.Copy(a)).(*AlarmUpdate)
}

func (a *ClimateState) Clone() *ClimateState {
	return copystructure.Must(copystructure.Copy(a)).(*ClimateState)
}

func (r *ZoneClimateReport) Clone() *ZoneClimateReport {
	return copystructure.Must(copystructure.Copy(r)).(*ZoneClimateReport)
}

func (i *TrackingInfo) Clone() *TrackingInfo {
	return copystructure.Must(copystructure.Copy(i)).(*TrackingInfo)
}

func (l *ServiceLog) Clone() *ServiceLog {
	return copystructure.Must(copystructure.Copy(l)).(*ServiceLog)
}
