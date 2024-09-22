package test

import (
	"testing"

	"github.com/mgechev/revive/internal/typeparams"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestReceiverNamingTypeParams(t *testing.T) {
	if !typeparams.Enabled() {
		t.Skip("type parameters are not enabled in the current build environment")
	}
	testRule(t, "receiver-naming-issue-669", &rule.ReceiverNamingRule{})
}

func TestReceiverNamingMaxLength(t *testing.T) {
	testRule(t, "receiver-naming-issue-1040", &rule.ReceiverNamingRule{},
		&lint.RuleConfig{Arguments: []any{"max-length=2"}})
}
