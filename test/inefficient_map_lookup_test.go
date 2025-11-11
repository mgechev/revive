package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestInefficientMapLookup(t *testing.T) {
	testRule(t, "inefficient_map_lookup", &rule.InefficientMapLookupRule{}, &lint.RuleConfig{})
}
