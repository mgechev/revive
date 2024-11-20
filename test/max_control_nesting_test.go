package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestMaxControlNestingDefault(t *testing.T) {
	testRule(t, "max_control_nesting_default", &rule.MaxControlNestingRule{}, &lint.RuleConfig{})
}

func TestMaxControlNesting(t *testing.T) {
	testRule(t, "max_control_nesting", &rule.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: []any{int64(2)}},
	)
}
