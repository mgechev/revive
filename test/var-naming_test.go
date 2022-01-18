package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNaming(t *testing.T) {
	testRule(t, "var-naming", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: []interface{}{[]interface{}{"ID"}, []interface{}{"VM"}},
	})

	testRule(t, "var-naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})
}
