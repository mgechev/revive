package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestEmptyBlock rule.
func TestEmptyBlock(t *testing.T) {
	testRule(t, "empty_block", &rule.EmptyBlockRule{})
}
