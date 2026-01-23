package tm

import (
	"log/slog"
)

// Indicate the wanted verbose level for logs
type VerboseLevel int

// MapVerboseLevel maps a verbose level to a given logrus.Level
func MapVerboseLevel(l VerboseLevel) slog.Level {
	levels := MapVerboseLevelList(l)
	return levels[len(levels)-1]
}

var LevelTrace slog.Level = -8

// MapVerboseLevelList maps a verbose level to a given list of
// slog.Level that should be enabled.
func MapVerboseLevelList(l VerboseLevel) []slog.Level {
	res := []slog.Level{slog.LevelError, slog.LevelWarn, slog.LevelInfo}

	if l >= 1 {
		res = append(res, slog.LevelDebug)
	}
	if l >= 2 {
		res = append(res, LevelTrace)
	}
	return res
}
