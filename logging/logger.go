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

	debugModeEnabled := os.Getenv("DEBUG") != ""
	if !debugModeEnabled {
		// by default, suppress all logging output
		return slog.New(slog.DiscardHandler), nil
	}

	var err error
	loggerFile, err = os.Create(logFile)
	if err != nil {
		return nil, err
	}

	logger = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stderr, loggerFile), nil))

	logger.Info("Logger initialized", "logFile", logFile)

	return logger, nil
}

// Close closes the logger file if it was opened.
func Close() error {
	if loggerFile == nil {
		return nil
	}
	return loggerFile.Close()
}
