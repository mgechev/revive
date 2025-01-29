package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFuncLengthDefault(t *testing.T) {
	testRule(t, "function_length_default", &rule.FunctionLength{}, &lint.RuleConfig{})
}

func TestFuncLengthLimitsStatements(t *testing.T) {
	testRule(t, "function_length1", &rule.FunctionLength{}, &lint.RuleConfig{
		Arguments: []any{int64(2), int64(100)},
	})
}

func TestFuncLengthLimitsLines(t *testing.T) {
	testRule(t, "function_length2", &rule.FunctionLength{}, &lint.RuleConfig{
		Arguments: []any{int64(100), int64(5)},
	})
}

func TestFuncLengthLimitsDeactivated(t *testing.T) {
	testRule(t, "function_length3", &rule.FunctionLength{}, &lint.RuleConfig{
		Arguments: []any{int64(0), int64(0)},
	})
}
