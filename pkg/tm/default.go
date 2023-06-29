// Package tm enables (or not) Open Telemetry at the package level.
package tm

import (
	"context"

	"github.com/sirupsen/logrus"
)

var defaultProvider Provider = newLocalProvider(0)

// NewLogger returns a *logrus.Entry associated with a domain.
func NewLogger(domain string) *logrus.Entry {
	return defaultProvider.NewLogger(domain)
}

// Shutdown shutdowns the telemetry package. It is always safe to call
// it even if the Open Telemetry was not previously enabled at the
// package level. In this case, it will do nothing.
func Shutdown(ctx context.Context) error {
	return defaultProvider.Shutdown(ctx)
}

// Enabled indicates that Open Telemetry global provider is enabled at
// the package level.
func Enabled() bool {
	return defaultProvider.Enabled()
}

// Sets up the package telemetry to use a given OpenTelemetry
// Endpoint. You should call this at the very beginning of your
// program, before instantiating any log.
func SetUpTelemetry(args OtelProviderArgs) {
	defaultProvider.Shutdown(context.Background())
	defaultProvider = newOtelProvider(args)
}

// Sets up the package to only log locally at a given level. Open
// Telemetry global provider will not be enabled. You should call this
// at the very beginning of your program, before instantiating any
// log.
func SetUpLocal(l VerboseLevel) {
	defaultProvider.Shutdown(context.Background())
	defaultProvider = newLocalProvider(l)
}
