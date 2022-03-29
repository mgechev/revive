package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseAny(t *testing.T) {
	testRule(t, "use-any", &rule.UseAnyRule{})
}
