package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportAliasNaming(t *testing.T) {
	testRule(t, "import-alias-naming", &rule.ImportAliasNamingRule{})
	testRule(t, "import-alias-naming", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{},
		},
	})
}

func TestImportAliasNaming_CustomConfig(t *testing.T) {
	testRule(t, "import-alias-naming-custom-config", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: []any{`^[a-z]+$`},
	})
}

func TestImportAliasNaming_CustomConfigWithMultipleRules(t *testing.T) {
	testRule(t, "import-alias-naming-custom-config-with-multiple-values", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"allowRegex": `^[a-z][a-z0-9]*$`,
				"denyRegex":  `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}

func TestImportAliasNaming_CustomConfigWithOnlyDeny(t *testing.T) {
	testRule(t, "import-alias-naming-custom-config-with-only-deny", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"denyRegex": `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}
