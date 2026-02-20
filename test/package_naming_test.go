package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageNaming(t *testing.T) {
	testRule(t, "package_naming_a_pkg_with_caps", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_top_level_pkg", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_lib_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_lib_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-package-name-collision-with-go-std": true},
		},
	})
	testRule(t, "package_naming_skip_package_name_checks_false", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_skip_package_name_checks_true", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skipPackageNameChecks": true},
		},
	})
	testRule(t, "package_naming_skip_package_name_checks_true", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-package-name-checks": true},
		},
	})
	testRule(t, "package_naming_meaningless_package_name", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_meaningless_package_name", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-package-name-checks": false},
		},
	})
	testRule(t, "package_naming_bad_package_name", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"skip-package-name-checks": false,
				"extra-bad-package-names":  []any{"helpers"},
			},
		},
	})
}
