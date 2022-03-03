package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// TestCallToGC test call-to-gc rule
func TestCallToGC(t *testing.T) {
	testRule(t, "call-to-gc", &rule.CallToGCRule{})
}
