package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIssue1100(t *testing.T) {
	args := []any{map[string]any{
		"maxLitCount": "2",
		"allowStrs":   "\"\"",
		"allowInts":   "0,1,2",
		"allowFloats": "0.0,1.0",
		"ignoreFuncs": "os\\.(CreateFile|WriteFile|Chmod|FindProcess),\\.Println,ignoredFunc,\\.Info",
	}}

	testRule(t, "goUnknown/issue1100", &rule.AddConstantRule{}, &lint.RuleConfig{
		Arguments: args,
	})
}
