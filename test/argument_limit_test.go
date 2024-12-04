package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestArgumentsLimitDefault(t *testing.T) {
	testRule(t, "argument_limit_default", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{})
}

func TestArgumentsLimit(t *testing.T) {
	testRule(t, "argument_limit", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{
		Arguments: []any{int64(3)},
	})
}

func BenchmarkArgumentsLimit(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "argument_limit_default", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{})
	}
}
