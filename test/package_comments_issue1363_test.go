package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageCommentsIssue1363(t *testing.T) {
	files := []string{
		"../testdata/package_comments_issue1363/myfile1.go",
		"../testdata/package_comments_issue1363/myfile2.go",
	}

	testRuleOnFiles(t, files, &rule.PackageCommentsRule{}, &lint.RuleConfig{})
}
