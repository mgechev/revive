package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUselessFallthrough(t *testing.T) {
	testRule(t, "useless_fallthrough", &rule.UselessFallthroughRule{})
}
