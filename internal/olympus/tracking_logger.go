package olympus

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/formicidae-tracker/olympus/pkg/tm"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
)

type TrackingLogger interface {
	TrackingInfo() *api.TrackingInfo
	PushDiskStatus(*api.DiskStatus)
}

type trackingLogger struct {
	mx    sync.RWMutex
	infos *api.TrackingInfo

	logger *logrus.Entry
}

func NewTrackingLogger(ctx context.Context, declaration *api.TrackingDeclaration) TrackingLogger {
	since := time.Now()
	if declaration.Since != nil {
		since = declaration.Since.AsTime()
	}

	logger := tm.NewLogger("tracking").WithContext(ctx)
	logger.WithField("declaration", declaration).
		WithField("since", since).
		Debug("new declaration")

	return &trackingLogger{
		infos: &api.TrackingInfo{
			Since: since,
			Stream: &api.StreamInfo{
				ExperimentName: declaration.ExperimentName,
				StreamURL:      path.Join("/olympus/", declaration.Hostname, "index.m3u8"),
				ThumbnailURL:   path.Join("/thumbnails/olympus/", declaration.Hostname+".jpg"),
			},
		},
		logger: logger,
	}
}

func (l *trackingLogger) TrackingInfo() (res *api.TrackingInfo) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.infos.Clone()
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

	l.logger.WithField("diskStatus", proto.MarshalTextString(s)).Trace("new disk status")
	l.infos.TotalBytes = Max(s.FreeBytes, s.TotalBytes)
	l.infos.FreeBytes = Max(0, s.FreeBytes)
	l.infos.BytesPerSecond = Max(0, s.BytesPerSecond)
}
