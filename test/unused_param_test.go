package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnusedParam(t *testing.T) {
	testRule(t, "unused_param", &rule.UnusedParamRule{})
	testRule(t, "unused_param", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: lint.Arguments{}})
	testRule(t, "unused_param", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"a": "^xxx"},
	}})
	testRule(t, "unused_param_custom_regex", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"allowRegex": "^xxx"},
	}})
	testRule(t, "unused_param_custom_regex", &rule.UnusedParamRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"allow-regex": "^xxx"},
	}})
}

func BenchmarkUnusedParam(b *testing.B) {
	for b.Loop() {
		testRule(b, "unused_param", &rule.UnusedParamRule{}, &lint.RuleConfig{})
	}
}
