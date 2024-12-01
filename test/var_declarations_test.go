package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarDeclarations(t *testing.T) {
	testRule(t, "var_declarations_type_inference", &rule.VarDeclarationsRule{}, &lint.RuleConfig{})
	testRule(t, "var_declarations_zero_value", &rule.VarDeclarationsRule{}, &lint.RuleConfig{})
}
