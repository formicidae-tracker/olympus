package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TrackingWatcher interface {
	Close() error
	Done() <-chan struct{}
	Timeouted() <-chan struct{}
	StreamURL() string
	Hostname() string
}

type trackingWatcher struct {
	done, timeouted, stop chan struct{}
	URL, host             string
}

func NewTrackingWatcher(host, URL string) TrackingWatcher {
	res := &trackingWatcher{
		done:      make(chan struct{}),
		stop:      make(chan struct{}),
		timeouted: make(chan struct{}),
		URL:       URL,
		host:      host,
	}
	go res.watch()
	return res
}

func (l *trackingWatcher) Done() <-chan struct{} {
	return l.done
}

func (l *trackingWatcher) Timeouted() <-chan struct{} {
	return l.timeouted
}

func (l *trackingWatcher) StreamURL() string {
	return l.URL
}

func (l *trackingWatcher) Hostname() string {
	return l.host
}

func (l *trackingWatcher) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("already closed")
		}
		<-l.done
	}()
	close(l.stop)
	return nil
}

func (l *trackingWatcher) watch() {
	defer close(l.done)

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	once := sync.Once{}
	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			if l.check() == false {
				once.Do(func() { close(l.timeouted) })
			}
		}
	}
}

func (l *trackingWatcher) check() bool {
	resp, err := http.Get(l.URL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
