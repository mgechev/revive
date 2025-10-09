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
	testRule(t, "var_naming_skip_initialism_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		}})
	testRule(t, "var_naming_skip_initialism_name_checks_false", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": false}},
		}})
	testRule(t, "var_naming_allowlist_blocklist_skip_initialism_name_checks", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{
			[]any{"ID"},
			[]any{"VM"},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		}})

	testRule(t, "var_naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})

	testRule(t, "var_naming_upper_case_const_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
	})
	testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"upper-case-const": true}}},
	})

	testRule(t, "var_naming_skip_package_name_checks_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_skip_package_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"skipPackageNameChecks": true}}},
	})
	testRule(t, "var_naming_skip_package_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{}, []any{map[string]any{"skip-package-name-checks": true}}},
	})
	testRule(t, "var_naming_meaningless_package_name", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_meaningless_package_name", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{},
			[]any{map[string]any{"skip-package-name-checks": false}},
		},
	})
	testRule(t, "var_naming_bad_package_name", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{},
			[]any{map[string]any{"skip-package-name-checks": false, "extra-bad-package-names": []any{"helpers"}}},
		},
	})
	testRule(t, "var_naming_top_level_pkg", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_std_lib_conflict", &rule.VarNamingRule{}, &lint.RuleConfig{})
	testRule(t, "var_naming_std_lib_conflict_skip", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{}, []any{},
			[]any{map[string]any{"skip-package-name-collision-with-go-std": true}},
		},
	})
}

func BenchmarkUpperCaseConstTrue(b *testing.B) {
	var t *testing.T
	for i := 0; i < b.N; i++ {
		testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
			Arguments: []any{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
		})
	}
}
