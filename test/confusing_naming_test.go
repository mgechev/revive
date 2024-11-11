package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestConfusingNaming rule.
func TestConfusingNaming(t *testing.T) {
	testRule(t, "confusing_naming1", &rule.ConfusingNamingRule{})
}
