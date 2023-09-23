package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestExportedWithDisableStutteringCheck(t *testing.T) {
	args := []any{"disableStutteringCheck"}

	testRule(t, "exported-issue-555", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}

func TestExportedWithChecksOnMethodsOfPrivateTypes(t *testing.T) {
	args := []any{"checkPrivateReceivers"}

	testRule(t, "exported-issue-552", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}

func TestExportedReplacingStuttersByRepetitive(t *testing.T) {
	args := []any{"sayRepetitiveInsteadOfStutters"}

	testRule(t, "exported-issue-519", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}
