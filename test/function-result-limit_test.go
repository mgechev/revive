package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestFunctionResultsLimit(t *testing.T) {
	testRule(t, "function-result-limit", &rule.FunctionResultsLimitRule{}, &lint.RuleConfig{
		Arguments: []interface{}{int64(3)},
	})
}
