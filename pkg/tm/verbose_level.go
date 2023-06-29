package tm

import "github.com/sirupsen/logrus"

// Indicate the wanted verbose level for logs
type VerboseLevel int

// MapVerboseLevel maps a verbose level to a given logrus.Level
func MapVerboseLevel(l VerboseLevel) logrus.Level {
	levels := MapVerboseLevelList(l)
	return levels[len(levels)-1]
}

// MapVerboseLevelLiss maps a verbose level to a given list of
// logrus.Level that should be enabled.
func MapVerboseLevelList(l VerboseLevel) []logrus.Level {
	res := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}

	if l >= 1 {
		res = append(res, logrus.InfoLevel)
	}
	if l >= 2 {
		res = append(res, logrus.DebugLevel)
	}
	if l >= 3 {
		res = append(res, logrus.TraceLevel)
	}
	return res
}
