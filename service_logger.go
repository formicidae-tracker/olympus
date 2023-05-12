package main

import (
	"sort"
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceLogger interface {
	Log(identifier string, on, graceful bool)
	Logs() []*api.ServiceEventLogs
	OnServices() []string
	OffServices() []string
}

type serviceLogger struct {
	mx sync.RWMutex

	logs map[string]*api.ServiceEventLogs
}

func (l *serviceLogger) Log(identifier string, on, graceful bool) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if _, ok := l.logs[identifier]; ok == false {
		l.logs[identifier] = &api.ServiceEventLogs{
			Identifier: identifier,
		}
	}
	if on == true {
		graceful = true
	}
	l.logs[identifier].Events = append(l.logs[identifier].Events, &api.ServiceEvent{
		Time:     timestamppb.New(time.Now()),
		On:       on,
		Graceful: graceful,
	})
}

func (l *serviceLogger) Logs() []*api.ServiceEventLogs {
	l.mx.RLock()
	defer l.mx.RUnlock()

	services := make([]string, 0, len(l.logs))
	logs := make([]*api.ServiceEventLogs, 0, len(l.logs))

	for idt := range l.logs {
		services = append(services, idt)
	}
	sort.Strings(services)

	for _, idt := range services {
		logs = append(logs, deepcopy.MustAnything(l.logs[idt]).(*api.ServiceEventLogs))
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
		res = append(res, log.Identifier)
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
		logs: make(map[string]*api.ServiceEventLogs),
	}
}
