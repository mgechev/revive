package test

import (
	"testing"

	"github.com/mgechev/revive/internal/ifelse"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// TestSuperfluousElse rule.
func TestSuperfluousElse(t *testing.T) {
	testRule(t, "superfluous-else", &rule.SuperfluousElseRule{})
	testRule(t, "superfluous-else-scope", &rule.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: []any{ifelse.PreserveScope}})
}
