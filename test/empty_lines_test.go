package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestEmptyLines rule.
func TestEmptyLines(t *testing.T) {
	testRule(t, "empty_lines", &rule.EmptyLinesRule{})
}
