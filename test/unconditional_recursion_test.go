package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnconditionalRecursion(t *testing.T) {
	testRule(t, "unconditional_recursion", &rule.UnconditionalRecursionRule{})
}
