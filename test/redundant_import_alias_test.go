package test

import (
	"github.com/mgechev/revive/rule"
	"testing"
)

func TestRedundantImportAlias(t *testing.T) {
	testRule(t, "redundant_import_alias", &rule.RedundantImportAlias{})
}
