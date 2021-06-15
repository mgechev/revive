package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestNestedStructs(t *testing.T) {
	testRule(t, "nested-structs", &rule.NestedStructs{}, &lint.RuleConfig{})
}
