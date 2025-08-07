package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUselessFallTrhough(t *testing.T) {
	testRule(t, "useless_fallthrough", &rule.UselessFallthroughRule{})
}
