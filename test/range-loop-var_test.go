package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRangeLoopVar(t *testing.T) {
	testRule(t, "range-loop-var", &rule.RangeLoopVarRule{}, &lint.RuleConfig{})
}
