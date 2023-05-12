package main

import (
	"sort"
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
)

type ServiceLogger interface {
	Log(identifier string, on, graceful bool)
	Logs() []api.ServiceEventList
	OnServices() []string
	OffServices() []string
}

type serviceLogger struct {
	mx sync.RWMutex

	logs map[string]*api.ServiceEventList
}

func (l *serviceLogger) Log(zone string, on, graceful bool) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, ok := l.logs[zone]; ok == false {
		l.logs[zone] = &api.ServiceEventList{
			Zone: zone,
		}
	}
	if on == true {
		graceful = true
	}
	l.logs[zone].Events = append(l.logs[zone].Events, api.ServiceEvent{
		Time:     time.Now(),
		On:       on,
		Graceful: graceful,
	})
}

func (l *serviceLogger) Logs() []api.ServiceEventList {
	l.mx.RLock()
	defer l.mx.RUnlock()

	services := make([]string, 0, len(l.logs))
	logs := make([]api.ServiceEventList, 0, len(l.logs))

	for idt := range l.logs {
		services = append(services, idt)
	}
	sort.Strings(services)

	for _, idt := range services {
		logs = append(logs, *deepcopy.MustAnything(l.logs[idt]).(*api.ServiceEventList))
	}

	return logs
}

func (l *serviceLogger) find(on bool) []string {
	res := make([]string, 0, len(l.logs))
	for _, log := range l.logs {
		events := log.Events
		if len(events) == 0 || events[len(events)-1].On != on {
			continue
		}
		res = append(res, log.Zone)
	}
	sort.Strings(res)
	return res
}

func (l *serviceLogger) OnServices() []string {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.find(true)
}

func (l *serviceLogger) OffServices() []string {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.find(false)
}

func NewServiceLogger() ServiceLogger {
	return &serviceLogger{
		logs: make(map[string]*api.ServiceEventList),
	}
}
