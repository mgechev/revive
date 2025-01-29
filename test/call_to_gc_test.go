package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestCallToGC(t *testing.T) {
	testRule(t, "call_to_gc", &rule.CallToGCRule{})
}
