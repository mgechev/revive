package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestRedundantCanonicalImport(t *testing.T) {
	testRule(t, "go1.11/redundant_canonical_import", &rule.RedundantCanonicalImport{})
	testRule(t, "redundant_canonical_import", &rule.RedundantCanonicalImport{})
}
