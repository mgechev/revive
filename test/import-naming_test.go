package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportNaming(t *testing.T) {
	testRule(t, "import-naming", &rule.ImportNamingRule{})
	testRule(t, "import-naming-custom-config", &rule.ImportNamingRule{}, &lint.RuleConfig{
		Arguments: []any{`^[a-z]+$`},
	})
}
