package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportAliasNaming(t *testing.T) {
	testRule(t, "import-alias-naming", &rule.ImportAliasNamingRule{})
	testRule(t, "import-alias-naming-custom-config", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: []any{`^[a-z]+$`},
	})
}
