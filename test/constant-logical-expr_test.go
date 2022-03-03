package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// ConstantLogicalExpr rule.
func TestConstantLogicalExpr(t *testing.T) {
	testRule(t, "constant-logical-expr", &rule.ConstantLogicalExprRule{})
}
