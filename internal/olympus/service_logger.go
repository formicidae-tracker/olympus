package olympus

import (
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/pkg/api"
)

type ServiceLogger interface {
	Log(identifier string, on, graceful bool)
	Logs() []api.ServiceLog
	OnServices() []string
	OffServices() []string
}

type serviceLogger struct {
	mx sync.RWMutex

	logs   *PersistentMap[*api.ServiceLog]
	logger *log.Logger
}

func (l *serviceLogger) Log(zone string, on, graceful bool) {
	now := time.Now()

	l.mx.Lock()
	defer l.mx.Unlock()

	log, ok := l.logs.Map[zone]
	if ok == false {
		log = &api.ServiceLog{
			Zone: zone,
		}
		l.logs.Map[zone] = log
	}
	if on == true {
		log.SetOn(now)
	} else {
		log.SetOff(now, graceful)
	}
	l.save(zone)
}

func (l *serviceLogger) Logs() []api.ServiceLog {
	l.mx.RLock()
	defer l.mx.RUnlock()

	services := make([]string, 0, len(l.logs.Map))
	logs := make([]api.ServiceLog, 0, len(l.logs.Map))

	for idt := range l.logs.Map {
		services = append(services, idt)
	}
	sort.Strings(services)

	for _, idt := range services {
		logs = append(logs, *deepcopy.MustAnything(l.logs.Map[idt]).(*api.ServiceLog))
	}

	return logs
}

func (l *serviceLogger) find(on bool) []string {
	res := make([]string, 0, len(l.logs.Map))
	for _, log := range l.logs.Map {
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

func (l *serviceLogger) save(zone string) {
	if err := l.logs.SaveKey(zone); err != nil {
		l.logger.Printf("%s", err)
	}
}

func NewServiceLogger() ServiceLogger {
	res := &serviceLogger{
		logs:   NewPersistentMap[*api.ServiceLog]("services"),
		logger: log.New(os.Stderr, "[services]: ", log.LstdFlags),
	}
	res.mx.Lock()
	defer res.mx.Unlock()
	return res
}
