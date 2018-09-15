package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestRangValInClosure(t *testing.T) {
	testRule(t, "range-val-in-closure", &rule.RangValInClosureRule{})
}
