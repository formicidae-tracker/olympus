package main

import "time"

type DeadlinePusher struct {
	when   time.Time
	period time.Duration
	timer  *time.Timer
}

func (p *DeadlinePusher) Push() <-chan time.Time {
	now := time.Now()
	p.when = now.Add(p.period)

	if !p.timer.Stop() {
		select {
		case <-p.timer.C:
		}
	}
	p.timer.Reset(p.period)
	return p.timer.C
}

func (p *DeadlinePusher) Check(now time.Time) <-chan time.Time {
	if now.After(p.when) == true {
		return nil
	}
	return p.timer.C
}

func NewDeadlinePusher(period time.Duration) (*DeadlinePusher, <-chan time.Time) {
	now := time.Now()
	res := &DeadlinePusher{
		when:   now.Add(period),
		period: period,
		timer:  time.NewTimer(period),
	}
	return res, res.timer.C
}
