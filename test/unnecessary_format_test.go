package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnnecessaryFormat(t *testing.T) {
	testRule(t, "unnecessary_format", &rule.UnnecessaryFormatRule{}, &lint.RuleConfig{})
}
