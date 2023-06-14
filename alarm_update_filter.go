package main

import (
	"path"
	"time"

	"github.com/formicidae-tracker/olympus/api"
)

type ZonedAlarmUpdate struct {
	Zone   string
	Update *api.AlarmUpdate
}

func (u ZonedAlarmUpdate) ID() string {
	return path.Join(u.Zone, u.Update.Identification)
}

type UpdateFilter func(outgoing chan<- ZonedAlarmUpdate, incoming <-chan ZonedAlarmUpdate)

type updateFilter struct {
	minimumOn time.Duration

	staged map[string]ZonedAlarmUpdate
	fired  map[string]api.AlarmLevel
}

func newUpdateFilter(minimumOn time.Duration) *updateFilter {
	return &updateFilter{
		minimumOn: minimumOn,
		staged:    make(map[string]ZonedAlarmUpdate),
		fired:     make(map[string]api.AlarmLevel),
	}
}

func FilterAlarmUpdates(minimumOn time.Duration) UpdateFilter {
	filter := newUpdateFilter(minimumOn)
	return filter.filter
}

func (f *updateFilter) filter(
	outgoing chan<- ZonedAlarmUpdate,
	incoming <-chan ZonedAlarmUpdate) {

	defer close(outgoing)

	var timer <-chan time.Time

	for {
		select {
		case u, ok := <-incoming:
			if ok == false {
				return
			}
			if f.stage(u) == true && timer == nil {
				timer = time.After(f.minimumOn)
			}
		case t := <-timer:
			oldest := f.unstage(outgoing, t)
			if len(f.staged) > 0 {
				timer = time.After(oldest.Add(f.minimumOn).Sub(t))
			}
		}
	}
}

func (f *updateFilter) stage(u ZonedAlarmUpdate) bool {
	id := u.ID()
	if u.Update.Status == api.AlarmStatus_OFF {
		delete(f.fired, id)
		delete(f.staged, id)
		return false
	}

	if firedLevel, ok := f.fired[id]; ok == true &&
		(firedLevel == api.AlarmLevel_EMERGENCY || u.Update.Level == api.AlarmLevel_WARNING) {
		return false
	}

	staged, ok := f.staged[id]
	if ok == false {
		f.staged[id] = u
		return true
	}
	if len(u.Update.Description) > 0 {
		staged.Update.Description = u.Update.Description
	}
	staged.Update.Level = u.Update.Level
	return false
}

func (f *updateFilter) unstage(outgoing chan<- ZonedAlarmUpdate, now time.Time) time.Time {
	toDelete := make([]string, 0, len(f.staged))

	oldest := now
	for idt, u := range f.staged {
		uTime := u.Update.Time.AsTime()
		if now.Sub(uTime) > f.minimumOn {
			outgoing <- u
			f.fired[idt] = u.Update.Level
			toDelete = append(toDelete, idt)
			continue
		}

		if uTime.Before(oldest) == true {
			oldest = uTime
		}
	}

	return oldest
}
