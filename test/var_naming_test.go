package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNaming(t *testing.T) {
	testRule(t, "var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"ID"}, []any{"VM"}},
	})

	testRule(t, "var_naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})

	testRule(t, "var_naming_upper_case_const_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
	})

	testRule(t, "var_naming_skip_package_name_checks_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_skip_package_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"skipPackageNameChecks": true}}},
	})
}
