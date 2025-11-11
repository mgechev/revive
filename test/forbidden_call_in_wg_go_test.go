package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestForbiddenCallInWgGo(t *testing.T) {
	testRule(t, "forbidden_call_in_wg_go", &rule.ForbiddenCallInWgGoRule{}, &lint.RuleConfig{})
	testRule(t, "go1.25/forbidden_call_in_wg_go", &rule.ForbiddenCallInWgGoRule{}, &lint.RuleConfig{})
}
