// Package logging provides a logger and related methods.
package logging

import (
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"
)

// GetLogger retrieves an instance of an application logger.
// The log level can be configured via the REVIVE_LOG_LEVEL environment variable.
// If REVIVE_LOG_LEVEL is not set, it defaults to WARN level.
//
//nolint:unparam // err is always nil, but is included in the signature for future extensibility.
func GetLogger() (*slog.Logger, error) {
	return getLogger(), nil
}

var getLogger = sync.OnceValue(initLogger(os.Stderr))

func initLogger(out io.Writer) func() *slog.Logger {
	return func() *slog.Logger {
		leveler := &slog.LevelVar{}
		opts := &slog.HandlerOptions{Level: leveler}

		// Check if REVIVE_LOG_LEVEL is set, otherwise default to WARN
		if logLevel := os.Getenv("REVIVE_LOG_LEVEL"); logLevel != "" {
			level := slog.LevelWarn
			_ = level.UnmarshalText([]byte(logLevel)) // Ignore error and default to WARN if invalid
			leveler.Set(level)
			logger := slog.New(slog.NewTextHandler(out, opts))

			logger.Info("Logger initialized", "logLevel", logLevel)

			return logger
		}

		// Default to WARN level
		leveler.Set(slog.LevelWarn)
		logger := slog.New(slog.NewTextHandler(out, opts))

		logger.Info("Logger initialized", "logLevel", slog.LevelWarn)

		return logger
	}
}

// InitForTesting initializes the logger singleton cache for testing purposes.
// This function should only be called in tests.
func InitForTesting(tb testing.TB, w io.Writer) {
	tb.Helper()
	getLogger = sync.OnceValue(initLogger(w))
}
