package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestBoolComp rule.
func TestBoolComp(t *testing.T) {
	testRule(t, "bool-comp", &rule.BoolCompRule{})
}
