package test

import (
"testing"

"github.com/jchen267/revive/rule"
)

// TestErrorPackageNaming rule.
func TestErrorPackageNaming(t *testing.T) {
	testRule(t, "error-package-naming", &rule.ErrorPackageNamingRule{})
}
