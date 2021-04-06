package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringFormat(t *testing.T) {
	testRule(t, "string-regex-capitalized", &rule.StringRegexRule{}, &lint.RuleConfig{
		Arguments: []interface{}{
			[]string{
				"/^[A-Z]/",
				"must be capitalized"}}})
}
