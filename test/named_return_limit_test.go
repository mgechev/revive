package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestNamedReturnLimit(t *testing.T) {
	testRule(t, "named_return_limit", &rule.ReturnLimitNamedRule{}, &lint.RuleConfig{})
}
