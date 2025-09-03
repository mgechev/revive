package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestWaitGroupDoneInWaitGroupGo(t *testing.T) {
	testRule(t, "waitgroup_done_in_waitgroup_go", &rule.WaitGroupDoneInWaitGroupGoRule{}, &lint.RuleConfig{})
	testRule(t, "go1.25/waitgroup_done_in_waitgroup_go", &rule.WaitGroupDoneInWaitGroupGoRule{}, &lint.RuleConfig{})
}
