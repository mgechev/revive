package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageNaming_convention(t *testing.T) {
	testRule(t, "package_naming_mixed_caps", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_mixed_caps", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-checks": false},
		},
	})
	testRule(t, "package_naming_mixed_caps_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-checks": true},
		},
	})
	testRule(t, "package_naming_underscore", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_underscore", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-checks": false},
		},
	})
	testRule(t, "package_naming_underscore_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-checks": true},
		},
	})
}

func TestPackageNaming_topLevel(t *testing.T) {
	testRule(t, "package_naming_top_level_pkg", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_top_level_pkg", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-checks": false},
		},
	})
	testRule(t, "package_naming_top_level_pkg_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-checks": true},
		},
	})
}

func TestPackageNaming_badNames(t *testing.T) {
	testRule(t, "package_naming_bad_default", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_default", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-checks": false},
		},
	})
	testRule(t, "package_naming_bad_default_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-checks": true},
		},
	})

	testRule(t, "package_naming_bad_extra", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-names": true},
		},
	})
	testRule(t, "package_naming_bad_extra_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_extra_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-names": false},
		},
	})

	testRule(t, "package_naming_bad_user_defined", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{"data"}},
		},
	})
	testRule(t, "package_naming_bad_user_defined_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_user_defined_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{}},
		},
	})
}

func TestPackageNaming_stdLibConflict(t *testing.T) {
	testRule(t, "package_naming_std_common_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_common_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": false},
		},
	})
	testRule(t, "package_naming_std_common_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": true},
		},
	})
	testRule(t, "package_naming_std_all_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": true},
		},
	})
	testRule(t, "package_naming_std_all_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_all_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": false},
		},
	})
}
