package logging_test

import (
	"os"
	"testing"

	"github.com/mgechev/revive/logging"
)

func TestGetLogger(t *testing.T) {
	t.Run("no debug", func(t *testing.T) {
		t.Setenv("DEBUG", "")

		logger, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected logger to be non-nil")
		}
		logger.Info("msg") // no output
	})

	t.Run("debug", func(t *testing.T) {
		t.Setenv("DEBUG", "1")
		t.Cleanup(func() {
			if err := logging.Close(); err != nil {
				t.Error(err)
			}
			if err := os.Remove("revive.log"); err != nil {
				t.Error(err)
			}
		})

		logger, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected logger to be non-nil")
		}
		if _, err := os.Stat("revive.log"); os.IsNotExist(err) {
			t.Error("expected revive.log file to be created")
		}

		logger2, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if logger != logger2 {
			t.Errorf("expected the same logger instance to be returned: logger1=%+v, logger2=%+v", logger, logger2)
		}
	})
}
