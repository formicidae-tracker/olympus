package main

import (
	"fmt"
	"net/rpc"
	"path"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type TrackingWatcher interface {
	Close() error
	Done() <-chan struct{}
	Timeouted() <-chan struct{}
	Stream() *StreamInfo
	Hostname() string
	Experiment() string
}

type trackingWatcher struct {
	done, timeouted, stop chan struct{}
	host                  string
	experiment            string
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

type TrackingWatcherArgs struct {
	Host       string
	Experiment string
	URL        string
}

func newTrackingWatcher(o TrackingWatcherArgs, period time.Duration) TrackingWatcher {
	res := &trackingWatcher{
		done:       make(chan struct{}),
		stop:       make(chan struct{}),
		timeouted:  make(chan struct{}),
		stream:     newStreamInfo(o.URL),
		host:       o.Host,
		experiment: o.Experiment,
		period:     period,
	}
	go res.watch()
	return res
}

func NewTrackingWatcher(o TrackingWatcherArgs) TrackingWatcher {
	return newTrackingWatcher(o, 20*time.Second)
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

func (l *trackingWatcher) Experiment() string {
	return l.experiment
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
	c, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s.local:%d", l.host, leto.LETO_PORT))
	if err != nil {
		return false
	}
	defer c.Close()
	status := leto.Status{}
	err = c.Call("Leto.Status", &leto.NoArgs{}, &status)
	if err != nil || status.Experiment == nil {
		return false
	}
	config := leto.TrackingConfiguration{}

	err = yaml.Unmarshal([]byte(status.Experiment.YamlConfiguration), &config)
	if err != nil {
		return false
	}
	return config.ExperimentName == l.experiment
}
