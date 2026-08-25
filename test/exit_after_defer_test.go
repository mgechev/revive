package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestExitAfterDefer(t *testing.T) {
	testRule(t, "exit_after_defer", &rule.ExitAfterDeferRule{})
}
