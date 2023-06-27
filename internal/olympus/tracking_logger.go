package olympus

import (
	"path"
	"sync"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/pkg/api"
	"golang.org/x/exp/constraints"
)

type TrackingLogger interface {
	TrackingInfo() *api.TrackingInfo
	PushDiskStatus(*api.DiskStatus)
}

type trackingLogger struct {
	mx    sync.RWMutex
	infos *api.TrackingInfo
}

func NewTrackingLogger(declaration *api.TrackingDeclaration) TrackingLogger {
	since := time.Now()
	if declaration.Since != nil {
		since = declaration.Since.AsTime()
	}

	return &trackingLogger{
		infos: &api.TrackingInfo{
			Since: since,
			Stream: &api.StreamInfo{
				ExperimentName: declaration.ExperimentName,
				StreamURL:      path.Join("/olympus/", declaration.Hostname, "index.m3u8"),
				ThumbnailURL:   path.Join("/thumbnails/olympus/", declaration.Hostname+".jpg"),
			},
		},
	}
}

func (l *trackingLogger) TrackingInfo() *api.TrackingInfo {
	l.mx.RLock()
	defer l.mx.RUnlock()
	return deepcopy.MustAnything(l.infos).(*api.TrackingInfo)
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func (l *trackingLogger) PushDiskStatus(s *api.DiskStatus) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.infos.TotalBytes = Max(s.FreeBytes, s.TotalBytes)
	l.infos.FreeBytes = Max(0, s.FreeBytes)
	l.infos.BytesPerSecond = Max(0, s.BytesPerSecond)
}
