package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnhandledError(t *testing.T) {
	testRule(t, "unhandled-error", &rule.UnhandledErrorRule{})
}

func TestUnhandledErrorWithIgnoreList(t *testing.T) {
	args := []interface{}{
		`unhandledError1`,
		`fmt\.Print`,
		`os\.(Create|WriteFile|Chmod)`,
		`net\..*`,
		`bytes\.Buffer\.Write`,
		`fixtures\.unhandledErrorStruct2\.reterr`,
	}

	testRule(t, "unhandled-error-w-ignorelist", &rule.UnhandledErrorRule{}, &lint.RuleConfig{Arguments: args})
}
