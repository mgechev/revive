package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestRangeValInClosure tests RangeValInClosure rule
func TestRangeValInClosure(t *testing.T) {
	testRule(t, "range-val-in-closure", &rule.RangeValInClosureRule{})
}
