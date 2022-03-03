package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

// Test that left and right side of Binary operators (only AND, OR) are swapable
func TestOptimizeOperandsOrder(t *testing.T) {
	testRule(t, "optimize-operands-order", &rule.OptimizeOperandsOrderRule{}, &lint.RuleConfig{})
}
