package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNaming(t *testing.T) {
	testRule(t, "var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}},
	})
	testRule(t, "var_naming_skip_initialism_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		},
	})
	testRule(t, "var_naming_skip_initialism_name_checks_false", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": false}},
		},
	})
	testRule(t, "var_naming_allowlist_blocklist_skip_initialism_name_checks", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{"ID"},
			[]any{"VM"},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		},
	})

	testRule(t, "var_naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})

	testRule(t, "var_naming_upper_case_const_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
	})
	testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upper-case-const": true}}},
	})
}

func BenchmarkUpperCaseConstTrue(b *testing.B) {
	for b.Loop() {
		testRule(b, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
		})
	}
}
