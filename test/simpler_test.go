package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestSimplerNaming rule.
func TestSimplerNaming(t *testing.T) {
	testRule(t, "simpler", &rule.SimplerRule{})
}
