package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// TestEarlyReturn tests early-return rule.
func TestEarlyReturn(t *testing.T) {
	testRule(t, "early-return", &rule.EarlyReturnRule{})
}
