package revivelib

import "github.com/mgechev/revive/lint"

// ExtraRule configures a new rule to be used with revive.
type ExtraRule struct {
	Rule          lint.Rule
	DefaultConfig lint.RuleConfig
}

// NewExtraRule returns a configured extra rule.
func NewExtraRule(rule lint.Rule, defaultConfig lint.RuleConfig) ExtraRule {
	return ExtraRule{
		Rule:          rule,
		DefaultConfig: defaultConfig,
	}
}
