package tm

import (
	"context"
	"log/slog"
	"os"
)

type localProvider struct {
}

var handler *slog.TextHandler = nil

func newLocalProvider(_ VerboseLevel) Provider {

	return localProvider{}
}

func (l localProvider) NewLogger(domain string) *slog.Logger {
	return slog.With("domain", domain)
}

func (l localProvider) Shutdown(context.Context) error {
	return nil
}

func (l localProvider) Enabled() bool {
	return false
}

func init() {
	handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: false})
	slog.SetDefault(slog.New(handler))
}
