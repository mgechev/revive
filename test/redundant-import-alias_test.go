package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// TestRedundantImportAlias rule.
func TestRedundantImportIgnoredAliases(t *testing.T) {

	args := []any{map[string]any{
		"ignoreUsed": true,
	}}

	testRule(t, "redundant-import-alias-ignored", &rule.RedundantImportAlias{}, &lint.RuleConfig{
		Arguments: args,
	})

}

func TestRedundantImportAlias(t *testing.T) {

	args := []any{map[string]any{
		"ignoreUsed": false,
	}}

	testRule(t, "redundant-import-alias", &rule.RedundantImportAlias{}, &lint.RuleConfig{
		Arguments: args,
	})

}
