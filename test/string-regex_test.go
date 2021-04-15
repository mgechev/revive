package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringRegex(t *testing.T) {
	testRule(t, "string-regex", &rule.StringRegexRule{}, &lint.RuleConfig{
		Arguments: []interface{}{
			[]interface{}{
				"stringRegexMethod1", // The first argument is checked by default
				"/^[A-Z]/",
				"must start with a capital letter"},

			[]interface{}{
				"stringRegexMethod2[2].d",
				"/[^\\.]$/"}}}) // Must not end with a period
}
