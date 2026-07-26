package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFunctionResultsLimit(t *testing.T) {
	testRule(t, "function_result_limit", &rule.FunctionResultsLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}
