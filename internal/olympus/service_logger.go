package olympus

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/formicidae-tracker/olympus/pkg/tm"
	"github.com/sirupsen/logrus"
)

type ServiceLogger interface {
	Log(ctx context.Context, identifier string, on, graceful bool)
	Logs() []api.ServiceLog
	OnServices() []string
	OffServices() []string
}

type serviceLogger struct {
	mx sync.RWMutex

	logs   *PersistentMap[*api.ServiceLog]
	logger *logrus.Entry
}

func (l *serviceLogger) Log(ctx context.Context, zone string, on, graceful bool) {
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
	l.save(ctx, zone)
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
		logs = append(logs, *l.logs.Map[idt].Clone())
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

func (l *serviceLogger) save(ctx context.Context, zone string) {
	if err := l.logs.SaveKey(zone); err != nil {
		l.logger.WithContext(ctx).
			WithFields(logrus.Fields{
				"zone":  zone,
				"error": err,
			}).Error("could not save to persistent storage")
	}
}

func NewServiceLogger() ServiceLogger {
	res := &serviceLogger{
		logs:   NewPersistentMap[*api.ServiceLog]("services"),
		logger: tm.NewLogger("services"),
	}
	res.mx.Lock()
	defer res.mx.Unlock()
	return res
}
