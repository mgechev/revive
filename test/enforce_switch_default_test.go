package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSwitchDefault(t *testing.T) {
	testRule(t, "enforce_switch_default", &rule.EnforceSwitchDefaultRule{})
	testRule(t, "enforce_switch_default_allow_no_default", &rule.EnforceSwitchDefaultRule{}, &lint.RuleConfig{
		Arguments: []any{"allowNoDefault"},
	})
	testRule(t, "enforce_switch_default_allow_not_last", &rule.EnforceSwitchDefaultRule{}, &lint.RuleConfig{
		Arguments: []any{"allowDefaultNotLast"},
	})
}
