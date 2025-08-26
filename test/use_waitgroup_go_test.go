package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUseWaitGroupGo(t *testing.T) {
	testRule(t, "use_waitgroup_go", &rule.UseWaitGroupGoRule{}, &lint.RuleConfig{})
	testRule(t, "go1.25/use_waitgroup_go", &rule.UseWaitGroupGoRule{}, &lint.RuleConfig{})
}
