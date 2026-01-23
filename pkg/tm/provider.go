package tm

import (
	"context"
	"log/slog"
)

// A Provider provides Logger in regards to a given telemetry scheme
type Provider interface {
	// Shutdown the telemetry engine
	Shutdown(context.Context) error
	// Creates a new logger associated with domain.
	NewLogger(domain string) *slog.Logger
	// Indicates if telemetry is enabled. If not, do not
	// instrumentalize thrid party libs.
	Enabled() bool
}
