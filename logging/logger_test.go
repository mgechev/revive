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
		t.Cleanup(func() { os.Remove("revive.log") })

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
	})

	t.Run("reuse logger", func(t *testing.T) {
		t.Setenv("DEBUG", "1")
		t.Cleanup(func() { os.Remove("revive.log") })

		logger1, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		logger2, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if logger1 != logger2 {
			t.Errorf("expected the same logger instance to be returned: logger1=%+v, logger2=%+v", logger1, logger2)
		}
	})
}
