package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRangeValAddress(t *testing.T) {
	testRule(t, "range-val-address", &rule.RangeValAddress{}, &lint.RuleConfig{})
}

func TestRangeValAddressAfterGo1_22(t *testing.T) {
	testRule(t, "go1.22/range-val-address", &rule.RangeValAddress{}, &lint.RuleConfig{})
}
