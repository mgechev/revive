package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnnecessaryIf(t *testing.T) {
	testRule(t, "unnecessary_if", &rule.UnnecessaryIfRule{}, &lint.RuleConfig{})
}
