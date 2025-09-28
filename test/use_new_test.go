package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUseNew(t *testing.T) {
	testRule(t, "use_new", &rule.UseNewRule{}, &lint.RuleConfig{})
	testRule(t, "go1.26/use_new", &rule.UseNewRule{}, &lint.RuleConfig{})
}
