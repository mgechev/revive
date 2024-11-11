package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// UselessBreak rule.
func TestUselessBreak(t *testing.T) {
	testRule(t, "useless_break", &rule.UselessBreak{})
}
