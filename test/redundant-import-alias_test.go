package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestRedundantImportAlias rule.
func TestRedundantImportAlias(t *testing.T) {
	testRule(t, "redundant-import-alias", &rule.RedundantImportAlias{})
}
