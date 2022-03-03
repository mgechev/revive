package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// TestSuperfluousElse rule.
func TestSuperfluousElse(t *testing.T) {
	testRule(t, "superfluous-else", &rule.SuperfluousElseRule{})
}
