package main

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

type TrackingWatcher interface {
	Close() error
	Done() <-chan struct{}
	Timeouted() <-chan struct{}
	Stream() *StreamInfo
	Hostname() string
}

type trackingWatcher struct {
	done, timeouted, stop chan struct{}
	host                  string
	stream                StreamInfo
	period                time.Duration
}

func newStreamInfo(URL string) StreamInfo {
	parent := path.Dir(path.Dir(URL))
	thumbnail := strings.Replace(path.Base(URL), ".m3u8", ".png", 1)
	return StreamInfo{
		StreamURL:    URL,
		ThumbnailURL: path.Join(parent, thumbnail),
	}
}

func newTrackingWatcher(host, URL string, period time.Duration) TrackingWatcher {
	res := &trackingWatcher{
		done:      make(chan struct{}),
		stop:      make(chan struct{}),
		timeouted: make(chan struct{}),
		stream:    newStreamInfo(URL),
		host:      host,
		period:    period,
	}
	go res.watch()
	return res
}

func NewTrackingWatcher(host, URL string) TrackingWatcher {
	return newTrackingWatcher(host, URL, 20*time.Second)
}

func (l *trackingWatcher) Done() <-chan struct{} {
	return l.done
}

func (l *trackingWatcher) Timeouted() <-chan struct{} {
	return l.timeouted
}

func (l *trackingWatcher) Stream() *StreamInfo {
	return &l.stream
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

	ticker := time.NewTicker(l.period)
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
	resp, err := http.Get(l.stream.StreamURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
