package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// TestEmptyBlock rule.
func TestEmptyBlock(t *testing.T) {
	testRule(t, "empty-block", &rule.EmptyBlockRule{})
}
