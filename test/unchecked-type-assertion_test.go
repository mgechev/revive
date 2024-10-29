package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUncheckedDynamicCast(t *testing.T) {
	testRule(t, "unchecked-type-assertion", &rule.UncheckedTypeAssertionRule{})
}

func TestUncheckedDynamicCastWithAcceptIgnored(t *testing.T) {
	args := []any{map[string]any{
		"acceptIgnoredAssertionResult": true,
	}}

	testRule(t, "unchecked-type-assertion-accept-ignored", &rule.UncheckedTypeAssertionRule{}, &lint.RuleConfig{Arguments: args})
}
