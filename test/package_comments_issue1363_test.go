package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageCommentsIssue1363(t *testing.T) {
	testRuleOnDir(t, "../testdata", "package_comments_issue1363/...", &rule.PackageCommentsRule{}, &lint.RuleConfig{})
}
