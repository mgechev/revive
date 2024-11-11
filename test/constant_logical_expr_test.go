package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// ConstantLogicalExpr rule.
func TestConstantLogicalExpr(t *testing.T) {
	testRule(t, "constant_logical_expr", &rule.ConstantLogicalExprRule{})
}
