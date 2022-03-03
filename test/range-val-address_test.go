package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestRangeValAddress(t *testing.T) {
	testRule(t, "range-val-address", &rule.RangeValAddress{}, &lint.RuleConfig{})
}
