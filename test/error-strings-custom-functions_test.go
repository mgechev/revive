package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestErrorStringsWithCustomFunctions(t *testing.T) {
	args := []interface{}{"pkgErrors.Wrap"}
	testRule(t, "error-strings-with-custom-functions", &rule.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}
