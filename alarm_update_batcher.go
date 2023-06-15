package main

import "time"

type notificationBuilder struct {
	batchPeriod time.Duration
}

type BatchingFilter func(outgoing chan<- []ZonedAlarmUpdate, incoming <-chan ZonedAlarmUpdate)

// Batches alarms updates so they are no more than a batch every
// silenceWindow
func BatchAlarmUpdate(batchPeriod time.Duration) BatchingFilter {
	filter := &notificationBuilder{
		batchPeriod: batchPeriod,
	}
	return filter.filter
}

func (b *notificationBuilder) filter(outgoing chan<- []ZonedAlarmUpdate, incoming <-chan ZonedAlarmUpdate) {

	defer close(outgoing)

	var timer <-chan time.Time
	var batch []ZonedAlarmUpdate

	for {
		select {
		case u, ok := <-incoming:
			if ok == false {
				return
			}
			if timer == nil {
				outgoing <- []ZonedAlarmUpdate{u}
				timer = time.After(b.batchPeriod)
			} else {
				batch = append(batch, u)
			}
		case <-timer:
			timer = nil
			if len(batch) > 0 {
				outgoing <- batch
				timer = time.After(b.batchPeriod)
			}
			batch = nil
		}
	}
}
