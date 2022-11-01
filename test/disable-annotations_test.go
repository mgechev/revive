package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDisabledAnnotations(t *testing.T) {
	testRule(t, "disable-annotations", &rule.ExportedRule{}, &lint.RuleConfig{})
}

func TestModifiedAnnotations(t *testing.T) {
	testRule(t, "disable-annotations2", &rule.VarNamingRule{}, &lint.RuleConfig{})
}

func TestDisableNextLineAnnotations(t *testing.T) {
	testRule(t, "disable-annotations3", &rule.VarNamingRule{}, &lint.RuleConfig{})
}
