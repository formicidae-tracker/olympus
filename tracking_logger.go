package main

type TrackingLogger interface {
	StreamInfo() *StreamInfo
	Close() error
}

type trackingLogger struct {
	infos *StreamInfo
}

func (*trackingLogger) Close() error {
	return nil
}

func (l *trackingLogger) StreamInfo() *StreamInfo {
	return l.infos
}
