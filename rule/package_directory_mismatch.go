package rule

import (
	"fmt"
	"path/filepath"

	"github.com/mgechev/revive/lint"
)

// PackageDirectoryMismatchRule detects when package name doesn't match directory name.
type PackageDirectoryMismatchRule struct{}

var skipDirs = map[string]struct{}{
	".":        {},
	"/":        {},
	"":         {},
	"internal": {},
	"testdata": {},
}

// Apply applies the rule to the given file.
func (*PackageDirectoryMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.IsTest() || file.Pkg.IsMain() {
		return nil
	}

	dirPath := filepath.Dir(file.Name)
	dirName := filepath.Base(dirPath)
	if _, found := skipDirs[dirName]; found {
		return nil
	}

	packageName := file.AST.Name.Name
	if packageName == dirName {
		return nil
	}

	return []lint.Failure{
		{
			Failure:    fmt.Sprintf("package name %q does not match directory name %q", packageName, dirName),
			Confidence: 1,
			Node:       file.AST.Name,
			Category:   lint.FailureCategoryNaming,
		},
	}
}

// Name returns the rule name.
func (*PackageDirectoryMismatchRule) Name() string {
	return "package-directory-mismatch"
}
