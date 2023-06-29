package tm

import (
	"context"

	"github.com/sirupsen/logrus"
)

type localProvider struct {
}

func newLocalProvider(l VerboseLevel) Provider {
	logrus.SetLevel(MapVerboseLevel(l))
	return localProvider{}
}

func (l localProvider) NewLogger(domain string) *logrus.Entry {
	return logrus.WithField("group", domain)
}

func (l localProvider) Shutdown(context.Context) error {
	return nil
}

func (l localProvider) Enabled() bool {
	return false
}
