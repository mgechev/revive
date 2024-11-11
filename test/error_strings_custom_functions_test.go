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
