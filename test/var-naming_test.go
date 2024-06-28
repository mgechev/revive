package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNaming(t *testing.T) {
	testRule(t, "var-naming", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"ID"}, []any{"VM"}},
	})

	testRule(t, "var-naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})

	testRule(t, "var-naming_upperCaseConst-false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var-naming_upperCaseConst-true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
	})

	testRule(t, "var-naming_skipPackageNameChecks-false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var-naming_skipPackageNameChecks-true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"skipPackageNameChecks": true}}},
	})

}
