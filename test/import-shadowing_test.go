package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestImportShadowing(t *testing.T) {
	testRule(t, "import-shadowing", &rule.ImportShadowingRule{})
}
