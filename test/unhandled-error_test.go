package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestUnhandledError(t *testing.T) {
	testRule(t, "unhandled-error", &rule.UnhandledErrorRule{})
}

func TestUnhandledErrorWithBlacklist(t *testing.T) {
	args := []interface{}{"os.Chdir", "unhandledError1"}

	testRule(t, "unhandled-error-w-ignorelist", &rule.UnhandledErrorRule{}, &lint.RuleConfig{Arguments: args})
}
