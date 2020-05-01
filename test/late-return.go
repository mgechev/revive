package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// late-return rule.
func TestLateReturn(t *testing.T) {
	testRule(t, "late-return", &rule.LateReturnRule{})
}
