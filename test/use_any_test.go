package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseAny(t *testing.T) {
	testRule(t, "go1.18/use_any", &rule.UseAnyRule{})
	testRule(t, "use_any", &rule.UseAnyRule{})
}
