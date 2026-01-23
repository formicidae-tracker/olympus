package tm

import (
	"context"
	"log/slog"
	"os"
)

type localProvider struct {
}

var handler slog.Handler = nil
var logger *slog.Logger = nil

func newLocalProvider(_ VerboseLevel) Provider {

	return localProvider{}
}

func (l localProvider) NewLogger(domain string) *slog.Logger {
	return logger.With("domain", domain)
}

func (l localProvider) Shutdown(context.Context) error {
	return nil
}

func (l localProvider) Enabled() bool {
	return false
}

func init() {
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: false})
	logger = slog.New(handler)
}
