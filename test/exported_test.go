package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/lint"
	"github.com/deepsourcelabs/revive/rule"
)

func TestExportedWithDisableStutteringCheck(t *testing.T) {
	args := []interface{}{"disableStutteringCheck"}

	testRule(t, "exported-issue-555", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}

func TestExportedWithChecksOnMethodsOfPrivateTypes(t *testing.T) {
	args := []interface{}{"checkPrivateReceivers"}

	testRule(t, "exported-issue-552", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}

func TestExportedReplacingStuttersByRepetitive(t *testing.T) {
	args := []interface{}{"sayRepetitiveInsteadOfStutters"}

	testRule(t, "exported-issue-519", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: args})
}
