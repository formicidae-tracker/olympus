package tm

import (
	"context"

	"github.com/sirupsen/logrus"
)

// A Provider provides Logger in regards to a given telemetry scheme
type Provider interface {
	// Shutdown the telemetry engine
	Shutdown(context.Context) error
	// Creates a new logger associated with domain.
	NewLogger(domain string) *logrus.Entry
	// Indicates if telemetry is enabled. If not, do not
	// instrumentalize thrid party libs.
	Enabled() bool
}
