package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSliceStyle_any(t *testing.T) {
	testRule(t, "enforce_slice_style_any", &rule.EnforceSliceStyleRule{})
}

func TestEnforceSliceStyle_make(t *testing.T) {
	testRule(t, "enforce_slice_style_make", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"make"},
	})
}

func TestEnforceSliceStyle_literal(t *testing.T) {
	testRule(t, "enforce_slice_style_literal", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"literal"},
	})
}

func TestEnforceSliceStyle_nil(t *testing.T) {
	testRule(t, "enforce_slice_style_nil", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"nil"},
	})
}
