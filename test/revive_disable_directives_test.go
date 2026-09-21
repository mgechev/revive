package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReviveDisableDirectives(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives", &rule.ExportedRule{}, lint.Config{})
}

func TestReviveDisableDirectives_Modified(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives_modified", &rule.VarNamingRule{}, lint.Config{})
}

func TestReviveDisableDirectives_NextLine(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives_next_line", &rule.VarNamingRule{}, lint.Config{})
}

func TestReviveDisableDirectives_SpecifyDisableReason(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives_specify_disable_reason", &rule.ExportedRule{}, lint.Config{
		Directives: lint.DirectivesConfig{
			"specify-disable-reason": {},
		},
	})
}

func TestReviveDisableDirectives_SpecifyDisableRule(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives_specify_disable_rule", &rule.ExportedRule{}, lint.Config{
		Directives: lint.DirectivesConfig{
			"specify-disable-rule": {},
		},
	})
}

func TestReviveDisableDirectives_SpecifyDisableReasonSpecifyDisableRule(t *testing.T) {
	testRuleWithLintConfig(t, "revive_disable_directives_specify_disable_reason_specify_disable_rule", &rule.ExportedRule{}, lint.Config{
		Directives: lint.DirectivesConfig{
			"specify-disable-reason": {},
			"specify-disable-rule":   {},
		},
	})
}
