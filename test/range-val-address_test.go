package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRangeValAddress(t *testing.T) {
	testRule(t, "range-val-address", &rule.RangeValAddress{}, &lint.RuleConfig{})
}
