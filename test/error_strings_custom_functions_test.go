package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestErrorStringsWithCustomFunctions(t *testing.T) {
	testRule(t, "error_strings_with_custom_functions", &rule.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"pkgErrors.Wrap"},
	})
}

func TestErrorStringsIssue1243(t *testing.T) {
	testRule(t, "error_strings_issue_1243", &rule.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"errors.Wrap"},
	})
}
