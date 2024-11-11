package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestDeepExit(t *testing.T) {
	testRule(t, "deep_exit", &rule.DeepExitRule{})
	testRule(t, "deep_exit_test", &rule.DeepExitRule{})
}
