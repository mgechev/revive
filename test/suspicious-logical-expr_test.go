package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// SuspiciousLogicalExpr rule.
func TestSuspiciousLogicalExpr(t *testing.T) {
	testRule(t, "suspicious-logical-expr", &rule.SuspiciousLogicalExprRule{})
}
