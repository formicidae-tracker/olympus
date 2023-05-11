package main

import (
	"path"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
)

type TrackingLogger interface {
	StreamInfo() *api.StreamInfo
}

type trackingLogger struct {
	infos *api.StreamInfo
}

func NewTrackingLogger(declaration *api.TrackingDeclaration) TrackingLogger {
	return &trackingLogger{
		infos: &api.StreamInfo{
			ExperimentName: declaration.ExperimentName,
			StreamURL:      path.Join("/olympus/hls", declaration.Hostname+".m3u8"),
			ThumbnailURL:   path.Join("/olympus", declaration.Hostname+".png"),
		},
	}
}

func (l *trackingLogger) StreamInfo() *api.StreamInfo {
	return deepcopy.MustAnything(l.infos).(*api.StreamInfo)
}
