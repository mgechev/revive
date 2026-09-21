package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestPackageComments(t *testing.T) {
	testRule(t, "package_comments/issue607_not_match", &rule.PackageCommentsRule{})
	testRule(t, "package_comments/issue607_match", &rule.PackageCommentsRule{})
	testRule(t, "package_comments/issue607_drift_not_match", &rule.PackageCommentsRule{})
}
