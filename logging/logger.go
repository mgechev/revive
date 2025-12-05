// Package logging provides a logger and related methods.
package logging

import (
	"io"
	"log/slog"
	"os"
)

const logFile = "revive.log"

var (
	logger     *slog.Logger
	loggerFile *os.File
)

// GetLogger retrieves an instance of an application logger which outputs
// to a file if the debug flag is enabled.
func GetLogger() (*slog.Logger, error) {
	if logger != nil {
		return logger, nil
	}

	// By default, suppress all logging output below level WARN,
	// and only log to stderr.
	leveler := new(slog.LevelVar)
	leveler.Set(slog.LevelWarn)
	opts := &slog.HandlerOptions{Level: leveler}

	var out io.Writer = os.Stderr

	debugModeEnabled := os.Getenv("DEBUG") != ""
	if debugModeEnabled {
		// In DEBUG mode, log all levels of output at level DEBUG and higher,
		// to both stderr and the logFile.
		leveler.Set(slog.LevelDebug)

		var err error
		loggerFile, err = os.Create(logFile)
		if err != nil {
			return nil, err
		}

		out = io.MultiWriter(os.Stderr, loggerFile)
	}

	logger = slog.New(slog.NewTextHandler(out, opts))

	logger.Info("Logger initialized", "logFile", logFile)

	return logger, nil
}

// Close closes the logger file if it was opened.
func Close() error {
	logger = nil
	if loggerFile == nil {
		return nil
	}
	return loggerFile.Close()
}
