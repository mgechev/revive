package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestImportShadowing(t *testing.T) {
	testRule(t, "import_shadowing", &rule.ImportShadowingRule{})
}
