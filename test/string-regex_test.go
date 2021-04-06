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
				"/^[A-Z]/",
				"must start with a capital letter"},

			[]interface{}{
				"/[^\\.]$/"}}}) // must not end with a period
}
