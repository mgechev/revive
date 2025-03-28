package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUncheckedDynamicCast(t *testing.T) {
	testRule(t, "unchecked_type_assertion", &rule.UncheckedTypeAssertionRule{})
}

func TestUncheckedDynamicCastWithAcceptIgnored(t *testing.T) {
	testRule(t, "unchecked_type_assertion_accept_ignored", &rule.UncheckedTypeAssertionRule{},
		&lint.RuleConfig{
			Arguments: []any{
				map[string]any{"acceptIgnoredAssertionResult": true},
			},
		},
	)
	testRule(t, "unchecked_type_assertion_accept_ignored", &rule.UncheckedTypeAssertionRule{},
		&lint.RuleConfig{
			Arguments: []any{
				map[string]any{"accept-ignored-assertion-result": true},
			},
		},
	)
}
