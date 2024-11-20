package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestArgumentLimitDefault(t *testing.T) {
	testRule(t, "argument_limit_default", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{})
}

func TestArgumentLimit(t *testing.T) {
	testRule(t, "argument_limit", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{
		Arguments: []any{int64(3)},
	})
}
