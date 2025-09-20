package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnnecessaryConditional(t *testing.T) {
	testRule(t, "unnecessary_conditional", &rule.UnnecessaryConditionalRule{}, &lint.RuleConfig{})
}
