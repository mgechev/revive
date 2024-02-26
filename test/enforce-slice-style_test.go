package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSliceStyle_any(t *testing.T) {
	testRule(t, "enforce-slice-style-any", &rule.EnforceSliceStyleRule{})
}

func TestEnforceSliceStyle_make(t *testing.T) {
	testRule(t, "enforce-slice-style-make", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"make"},
	})
}

func TestEnforceSliceStyle_literal(t *testing.T) {
	testRule(t, "enforce-slice-style-literal", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"literal"},
	})
}

func TestEnforceSliceStyle_nil(t *testing.T) {
	testRule(t, "enforce-slice-style-nil", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"nil"},
	})
}
