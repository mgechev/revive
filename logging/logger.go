// Package logging provides a logger and related methods.
package logging

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

// GetLogger retrieves an instance of an application logger which outputs
// to a file if the debug flag is enabled
func GetLogger() (*log.Logger, error) {
	if logger != nil {
		return logger, nil
	}

	var writer io.Writer
	var err error

	debugModeEnabled := os.Getenv("DEBUG") == "1"
	if debugModeEnabled {
		writer, err = os.Create("revive.log")
		if err != nil {
			return nil, err
		}
	} else {
		// Suppress all logging output if debug mode is disabled
		writer = io.Discard
	}

	logger = log.New(writer, "", log.LstdFlags)

	if !debugModeEnabled {
		// Clear all flags to skip log output formatting step to increase
		// performance somewhat if we're not logging anything
		logger.SetFlags(0)
	}

	logger.Println("Logger initialised")

	return logger, nil
}
