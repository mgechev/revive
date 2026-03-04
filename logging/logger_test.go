package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mgechev/revive/logging"
)

func TestGetLogger(t *testing.T) {
	t.Run("default WARN level", func(t *testing.T) {
		t.Setenv("REVIVE_LOG_LEVEL", "")
		var buf bytes.Buffer
		logging.InitForTesting(&buf)

		logger, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected logger to be non-nil")
		}

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		got := buf.String()

		if strings.Contains(got, "level=DEBUG") || strings.Contains(got, "level=INFO") {
			t.Errorf("unexpected output: got %q", got)
		}
		if want := `level=WARN msg="warn message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
		if want := `level=ERROR msg="error message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
	})

	t.Run("REVIVE_LOG_LEVEL debug", func(t *testing.T) {
		t.Setenv("REVIVE_LOG_LEVEL", "debug")
		var buf bytes.Buffer
		logging.InitForTesting(&buf)

		logger, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected logger to be non-nil")
		}

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		got := buf.String()

		if want := `level=INFO msg="Logger initialized" logLevel=debug`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
		if want := `level=DEBUG msg="debug message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
		if want := `level=INFO msg="info message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
		if want := `level=WARN msg="warn message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
		if want := `level=ERROR msg="error message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
	})

	t.Run("REVIVE_LOG_LEVEL invalid defaults to warn", func(t *testing.T) {
		t.Setenv("REVIVE_LOG_LEVEL", "invalid")
		var buf bytes.Buffer
		logging.InitForTesting(&buf)

		logger, err := logging.GetLogger()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected logger to be non-nil")
		}

		logger.Warn("warn message")

		got := buf.String()

		if strings.Contains(got, "level=DEBUG") || strings.Contains(got, "level=INFO") {
			t.Errorf("unexpected output: got %q", got)
		}
		if want := `level=WARN msg="warn message"`; !strings.Contains(got, want) {
			t.Errorf("expected output to contains %q, got %q", want, got)
		}
	})

	t.Run("same logger instance returned", func(t *testing.T) {
		t.Setenv("REVIVE_LOG_LEVEL", "info")
		var buf bytes.Buffer
		logging.InitForTesting(&buf)

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
