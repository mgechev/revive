package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestErrorStringsWithCustomFunctions(t *testing.T) {
	args := []any{"pkgErrors.Wrap"}
	testRule(t, "error_strings_with_custom_functions", &rule.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}

func TestErrorStringsIssue1243(t *testing.T) {
	args := []any{"errors.Wrap"}
	testRule(t, "error_strings_issue_1243", &rule.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}
