package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestJsonDataFormat(t *testing.T) {
	testRule(t, "json_data_format_atomic", &rule.AtomicRule{})
}

func TestJsonDataFormatVarNaming(t *testing.T) {
	testRule(t, "json_data_format_var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{})
}
