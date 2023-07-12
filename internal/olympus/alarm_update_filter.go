package olympus

import (
	"path"
	"strings"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
)

type ZonedAlarmUpdate struct {
	Zone   string
	Update *api.AlarmUpdate
}

func AppendSuffix(str string, suffix string) string {
	if strings.HasSuffix(str, suffix) == true {
		return str
	}
	return str + suffix
}

func (u ZonedAlarmUpdate) ID() string {
	if u.Update == nil {
		return AppendSuffix(u.Zone, "/")
	}
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
			if u.Update == nil {
				f.cleanUpFired(u.Zone)
			} else if f.stage(u) == true && timer == nil {
				wait := u.Update.Time.AsTime().Add(f.minimumOn).Sub(time.Now())
				timer = time.After(wait)
			}
		case t := <-timer:
			timer = nil
			oldest := f.unstage(outgoing, t)
			if len(f.staged) > 0 {
				wait := oldest.Add(f.minimumOn).Sub(t)
				timer = time.After(wait)
			}
		}
	}
}

func (f *updateFilter) cleanUpFired(zone string) {
	zone = AppendSuffix(zone, "/")
	toDelete := make([]string, 0, len(f.fired))

	for id := range f.fired {
		if strings.HasPrefix(id, zone) == true {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(f.fired, id)
	}
}

func (f *updateFilter) stage(u ZonedAlarmUpdate) bool {
	id := u.ID()
	if u.Update.Status == api.AlarmStatus_OFF {
		delete(f.fired, id)
		delete(f.staged, id)
		return false
	}

	if firedLevel, ok := f.fired[id]; ok == true {
		if firedLevel >= u.Update.Level {
			return false
		}
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

	for _, s := range toDelete {
		delete(f.staged, s)
	}

	return oldest
}
