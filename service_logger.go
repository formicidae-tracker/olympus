package main

import (
	"sort"
	"time"
)

type ServiceLogger interface {
	Log(identifier string, on, graceful bool)
	Logs() [][]ServiceEvent
	OnServices() []string
	OffServices() []string
}

type serviceLogger struct {
	indexes map[string]int
	logs    [][]ServiceEvent
}

func (l *serviceLogger) Log(identifier string, on, graceful bool) {
	index := l.getOrNew(identifier)
	if on == true {
		graceful = true
	}
	l.logs[index] = append(l.logs[index], ServiceEvent{
		Identifier: identifier,
		Time:       time.Now(),
		On:         on,
		Graceful:   graceful,
	})
}

func (l *serviceLogger) Logs() [][]ServiceEvent {
	return l.logs
}

func (l *serviceLogger) find(on bool) []string {
	res := make([]string, 0, len(l.logs))
	for _, logs := range l.logs {
		if len(logs) == 0 || logs[len(logs)-1].On != on {
			continue
		}
		res = append(res, logs[0].Identifier)
	}
	return res
}

func (l *serviceLogger) OnServices() []string {
	return l.find(true)
}

func (l *serviceLogger) OffServices() []string {
	return l.find(false)
}

func (l *serviceLogger) getOrNew(identifier string) int {
	if index, ok := l.indexes[identifier]; ok == true {
		return index
	}
	l.indexes[identifier] = len(l.logs)
	l.logs = append(l.logs, nil)
	l.sort()
	return l.indexes[identifier]
}

func (l *serviceLogger) sort() {
	oldIndexes := make(map[string]int)
	oldData := make([][]ServiceEvent, len(l.logs))
	copy(oldData, l.logs)
	identifiers := make([]string, 0, len(oldIndexes))
	for identifier, index := range l.indexes {
		oldIndexes[identifier] = index
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)
	for index, identifier := range identifiers {
		l.indexes[identifier] = index
		l.logs[index] = oldData[oldIndexes[identifier]]
	}
}

func NewServiceLogger() ServiceLogger {
	return &serviceLogger{
		indexes: make(map[string]int),
	}
}
