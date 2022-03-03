package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

// UselessBreak rule.
func TestUselessBreak(t *testing.T) {
	testRule(t, "useless-break", &rule.UselessBreak{})
}
