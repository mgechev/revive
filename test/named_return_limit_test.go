package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestNamedReturnLimitDefault(t *testing.T) {
	testRule(t, "named_return_min", &rule.NamedReturnMinRule{}, &lint.RuleConfig{})
}

func TestNamedReturnLimitConfigured(t *testing.T) {
	testRule(t, "named_return_min_configured", &rule.NamedReturnMinRule{}, &lint.RuleConfig{
		Arguments: []any{int64(1)},
	})
}
