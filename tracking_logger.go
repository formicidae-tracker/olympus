package main

import (
	"path"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/api"
)

type TrackingLogger interface {
	TrackingInfo() *api.TrackingInfo
}

type trackingLogger struct {
	infos *api.TrackingInfo
}

func NewTrackingLogger(declaration *api.TrackingDeclaration) TrackingLogger {
	return &trackingLogger{
		infos: &api.TrackingInfo{
			Stream: &api.StreamInfo{
				ExperimentName: declaration.ExperimentName,
				StreamURL:      path.Join("/olympus/hls", declaration.Hostname+".m3u8"),
				ThumbnailURL:   path.Join("/olympus", declaration.Hostname+".png"),
			},
		},
	}
}

func (l *trackingLogger) TrackingInfo() *api.TrackingInfo {
	return deepcopy.MustAnything(l.infos).(*api.TrackingInfo)
}
