package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedParam(t *testing.T) {
	testRule(t, "unused_param", &rule.UnusedParamRule{})
	testRule(t, "unused_param", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: []any{}})
	testRule(t, "unused_param", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"a": "^xxx"},
	}})
	testRule(t, "unused_param_custom_regex", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"allowRegex": "^xxx"},
	}})
	testRule(t, "unused_param_custom_regex", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: []any{
		map[string]any{"allow-regex": "^xxx"},
	}})
}

func BenchmarkUnusedParam(b *testing.B) {
	var t *testing.T
	for i := 0; i <= b.N; i++ {
		testRule(t, "unused_param", &rule.UnusedParamRule{})
	}
}
