package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceMapStyle_any(t *testing.T) {
	testRule(t, "enforce-map-style-any", &rule.EnforceMapStyleRule{})
}

func TestEnforceMapStyle_make(t *testing.T) {
	testRule(t, "enforce-map-style-make", &rule.EnforceMapStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"make"},
	})
}

func TestEnforceMapStyle_literal(t *testing.T) {
	testRule(t, "enforce-map-style-literal", &rule.EnforceMapStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"literal"},
	})
}
