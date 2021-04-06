package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringRegex(t *testing.T) {
	testRule(t, "string-regex", &rule.StringRegexRule{}, &lint.RuleConfig{
		Arguments: []interface{}{
			[]string{
				"/^[A-Z]/",
				"must be capitalized"},

			[]string{
				"/[^\\.]$/"}}}) // must not end with a period
}
