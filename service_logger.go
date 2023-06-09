package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/adrg/xdg"
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

	logs   map[string]*api.ServiceEventList
	logger *log.Logger
}

func (l *serviceLogger) Log(zone string, on, graceful bool) {
	l.mx.Lock()
	defer l.mx.Unlock()
	now := time.Now()
	list, ok := l.logs[zone]
	if ok == false {
		list = &api.ServiceEventList{
			Zone: zone,
		}
		l.logs[zone] = list
	}
	if on == true {
		list.SetOn(now)
	} else {
		list.SetOff(now, graceful)
	}
	l.save(zone)
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
		if on != log.On() {
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

func (l *serviceLogger) saveUnsafe(zone string) error {
	filename := l.zoneFilePath(zone)
	dirpath := filepath.Dir(filename)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return fmt.Errorf("could not make %s: %w", dirpath, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", filename, err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(l.logs[zone]); err != nil {
		return fmt.Errorf("could not write %s: %w", filename, err)
	}
	return nil
}

func (l *serviceLogger) save(zone string) {
	if err := l.saveUnsafe(zone); err != nil {
		l.logger.Printf("%s", err)
	}
}

func (l *serviceLogger) loadUnsafe(name string) error {
	if filepath.Ext(name) != ".json" {
		return nil
	}
	filename := filepath.Join(l.dataPath(), filepath.Base(name))
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", filename, err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	res := &api.ServiceEventList{}
	if err := dec.Decode(res); err != nil {
		return fmt.Errorf("could not read %s: %w", filename, err)
	}
	l.logs[res.Zone] = res
	return nil
}

func (l *serviceLogger) dataPath() string {
	return filepath.Join(xdg.DataHome, "fort", "olympus", "services")
}

func (l *serviceLogger) zoneFilePath(zone string) string {
	return filepath.Join(l.dataPath(), "zone"+".json")
}

func (l *serviceLogger) reload() {
	entries, err := os.ReadDir(l.dataPath())
	if err != nil {
		l.logger.Printf("could not read %s: %s", l.dataPath(), err)
		return
	}

	for _, e := range entries {
		if err := l.loadUnsafe(e.Name()); err != nil {
			l.logger.Println(err.Error())
		}
	}
}

func NewServiceLogger() ServiceLogger {
	res := &serviceLogger{
		logs:   make(map[string]*api.ServiceEventList),
		logger: log.New(os.Stderr, "[services]: ", log.LstdFlags),
	}
	res.mx.Lock()
	defer res.mx.Unlock()
	res.reload()
	return res
}
