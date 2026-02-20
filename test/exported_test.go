package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestExportedWithDisableStutteringCheck(t *testing.T) {
	testRule(t, "exported_issue_555", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"disableStutteringCheck"}})
}

func TestExportedWithChecksOnMethodsOfPrivateTypes(t *testing.T) {
	testRule(t, "exported_issue_552", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"checkPrivateReceivers"}})
}

func TestExportedReplacingStuttersByRepetitive(t *testing.T) {
	testRule(t, "exported_issue_519", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"sayRepetitiveInsteadOfStutters"}})
}

func TestCheckPublicInterfaceOption(t *testing.T) {
	testRule(t, "exported_issue_1002", &rule.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"checkPublicInterface"}})
}

func TestCheckDisablingOnDeclarationTypes(t *testing.T) {
	testRule(t, "exported_issue_1045", &rule.ExportedRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"disableChecksOnConstants", "disableChecksOnFunctions", "disableChecksOnMethods", "disableChecksOnTypes", "disableChecksOnVariables"},
	})
	testRule(t, "exported_issue_1045", &rule.ExportedRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"disable-checks-on-constants", "disable-checks-on-functions", "disable-checks-on-methods", "disable-checks-on-types", "disable-checks-on-variables"},
	})
}

func TestCheckDirectiveComment(t *testing.T) {
	testRule(t, "exported_issue_1202", &rule.ExportedRule{}, &lint.RuleConfig{})
}

func TestCheckDeprecationComment(t *testing.T) {
	testRule(t, "exported_issue_1231", &rule.ExportedRule{}, &lint.RuleConfig{})
}

func TestExportedMainPackage(t *testing.T) {
	testRule(t, "exported_main", &rule.ExportedRule{}, &lint.RuleConfig{})
}

func TestCommentVariations(t *testing.T) {
	testRule(t, "exported_issue_1235", &rule.ExportedRule{}, &lint.RuleConfig{})
}
