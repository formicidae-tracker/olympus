package main

import (
	"path"

	"github.com/barkimedes/go-deepcopy"
	"github.com/formicidae-tracker/olympus/olympuspb"
)

type TrackingLogger interface {
	StreamInfo() *StreamInfo
}

type trackingLogger struct {
	infos *StreamInfo
}

func NewTrackingLogger(declaration *olympuspb.TrackingDeclaration) TrackingLogger {
	return &trackingLogger{
		infos: &StreamInfo{
			ExperimentName: declaration.ExperimentName,
			StreamURL:      path.Join("/olympus/hls", declaration.Hostname+".m3u8"),
			ThumbnailURL:   path.Join("/olympus", declaration.Hostname+".png"),
		},
	}
}

func (l *trackingLogger) StreamInfo() *StreamInfo {
	return deepcopy.MustAnything(l.infos).(*StreamInfo)
}
