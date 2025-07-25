package rule

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgechev/revive/lint"
)

// PackageDirectoryMismatchRule detects when package name doesn't match directory name.
type PackageDirectoryMismatchRule struct {
	ignoredDirs *regexp.Regexp
}

const defaultIgnoredDirs = "testdata"

// Configure validates the rule configuration, and configures the rule accordingly.
func (r *PackageDirectoryMismatchRule) Configure(arguments lint.Arguments) error {
	ignoredDirs := defaultIgnoredDirs

	if len(arguments) > 0 {
		var ok bool
		ignoredDirs, ok = arguments[0].(string)
		if !ok {
			return fmt.Errorf("invalid argument type for ignored directories: expected string, got %T", arguments[0])
		}
	}

	if ignoredDirs == "" {
		r.ignoredDirs = nil
		return nil
	}

	var err error
	r.ignoredDirs, err = regexp.Compile(ignoredDirs)
	if err != nil {
		return fmt.Errorf("invalid regex for ignored directories: %w", err)
	}

	return nil
}

// normalizeName removes hyphens, underscores, and dots from the name
// to allow matching between directory names like "foo-bar.buz" and package names like "foobarbuz".
func normalizeName(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, ".", "")
	return name
}

// skipDirs contains directory names that should be unconditionally ignored when checking.
var skipDirs = map[string]struct{}{
	".": {},
	"/": {},
	"":  {},
}

// isVersionDirectory checks if a directory name is a version directory (v1, v2, etc.)
func isVersionDirectory(name string) bool {
	if len(name) < 2 || (name[0] != 'v' && name[0] != 'V') {
		return false
	}

	for i := 1; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}

	return true
}

// Apply applies the rule to the given file.
func (r *PackageDirectoryMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.IsTest() || file.Pkg.IsMain() {
		return nil
	}

	dirPath := filepath.Dir(file.Name)
	dirName := filepath.Base(dirPath)

	if r.ignoredDirs != nil && r.ignoredDirs.MatchString(dirPath) {
		return nil
	}

	// Use the parent directory for comparison if the immediate directory is a version directory (v1, v2, etc.).
	if isVersionDirectory(dirName) {
		dirName = filepath.Base(filepath.Dir(dirPath))
	}

	// Files directly in 'internal/' (like 'internal/abcd.go') should not be checked.
	// But files in subdirectories of 'internal/' (like 'internal/foo/abcd.go') should be checked.
	if dirName == "internal" {
		return nil
	}

	// Check if we got an invalid directory.
	if _, skipDir := skipDirs[dirName]; skipDir {
		return nil
	}

	packageName := file.AST.Name.Name
	if normalizeName(dirName) == normalizeName(packageName) {
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
