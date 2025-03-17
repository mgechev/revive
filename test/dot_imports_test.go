package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDotImportsDefault(t *testing.T) {
	testRule(t, "dot_imports_default", &rule.DotImportsRule{}, &lint.RuleConfig{})
}

func TestDotImports(t *testing.T) {
	testRule(t, "dot_imports", &rule.DotImportsRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{
			"allowedPackages": []any{"errors", "context", "github.com/BurntSushi/toml"},
		}},
	})
	testRule(t, "dot_imports", &rule.DotImportsRule{}, &lint.RuleConfig{
		Arguments: []any{map[string]any{
			"allowed-packages": []any{"errors", "context", "github.com/BurntSushi/toml"},
		}},
	})
}
