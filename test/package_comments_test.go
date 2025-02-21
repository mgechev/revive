package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageCommentsIssue607NotMatch(t *testing.T) {
	testRule(t, "package_comments/issue607_not_match", &rule.PackageCommentsRule{}, &lint.RuleConfig{})
}

func TestPackageCommentsIssue607Match(t *testing.T) {
	testRule(t, "package_comments/issue607_match", &rule.PackageCommentsRule{}, &lint.RuleConfig{})
}
