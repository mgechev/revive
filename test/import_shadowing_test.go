package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestImportShadowing(t *testing.T) {
	testRule(t, "import_shadowing", &rule.ImportShadowingRule{})
	testRule(t, "import_shadowing_issue_1435", &rule.ImportShadowingRule{})
	testRule(t, "import_shadowing_issue_1435_v1", &rule.ImportShadowingRule{})
}
