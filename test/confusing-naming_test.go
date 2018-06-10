package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestConfusingNaming rule.
func TestConfusingNaming(t *testing.T) {
	testRule(t, "confusing-naming", &rule.ConfusingNamingRule{})
}
