package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLineLengthLimitDefault(t *testing.T) {
	testRule(t, "line_length_limit_default", &rule.LineLengthLimitRule{}, &lint.RuleConfig{})
}

func TestLineLengthLimit(t *testing.T) {
	testRule(t, "line_length_limit", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: []any{int64(100)},
	})
}
