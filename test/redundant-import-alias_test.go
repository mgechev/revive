package test

import (
	"github.com/mgechev/revive/rule"
	"testing"
)

// TestRedundantImportAlias rule.
func TestRedundantImportAlias(t *testing.T) {
	testRule(t, "redundant-import-alias", &rule.RedundantImportAlias{})
}
