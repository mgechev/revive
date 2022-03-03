package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// Atomic rule.
func TestAtomic(t *testing.T) {
	testRule(t, "atomic", &rule.AtomicRule{})
}
