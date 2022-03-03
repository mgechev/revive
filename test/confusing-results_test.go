package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestConfusingResults(t *testing.T) {
	testRule(t, "confusing-results", &rule.ConfusingResultsRule{})
}
