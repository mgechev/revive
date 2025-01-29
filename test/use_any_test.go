package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseAny(t *testing.T) {
	testRule(t, "use_any", &rule.UseAnyRule{})
}
