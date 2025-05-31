package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSwitchStyle(t *testing.T) {
	testRule(t, "enforce_switch_style", &rule.EnforceSwitchStyleRule{})
	testRule(t, "enforce_switch_style_allow_no_default", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"allowNoDefault"},
	})
	testRule(t, "enforce_switch_style_allow_not_last", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"allowDefaultNotLast"},
	})
	testRule(t, "enforce_switch_style_allow_no_default_allow_not_last", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"allowNoDefault", "allowDefaultNotLast"},
	})
}
