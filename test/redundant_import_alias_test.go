package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestRedundantImportAlias(t *testing.T) {
	testRule(t, "redundant_import_alias", &rule.RedundantImportAlias{})
}
