package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestAddConstantWithDefaultArguments(t *testing.T) {
	testRule(t, "add_constant_default", &rule.AddConstantRule{}, &lint.RuleConfig{})
}

func TestAddConstantWithArguments(t *testing.T) {
	testRule(t, "add_constant", &rule.AddConstantRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{
			"maxLitCount": "2",
			"allowStrs":   "\"\"",
			"allowInts":   "0,1,2",
			"allowFloats": "0.0,1.0",
			"ignoreFuncs": "os\\.(CreateFile|WriteFile|Chmod|FindProcess),\\.Println,ignoredFunc,\\.Info",
		}},
	})
}
