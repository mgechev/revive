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
	r.ignoredDirs = nil
	if len(arguments) < 1 {
		return r.buildIgnoreRegex([]string{defaultIgnoredDirs})
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument type: expected map[string]any, got %T", arguments[0])
	}

	for k, v := range args {
		if !isRuleOption(k, "ignore-directories") {
			return fmt.Errorf("unknown argument %s for %s rule", k, r.Name())
		}
		ignoredDirs, ok := v.([]string)
		if !ok {
			return fmt.Errorf("invalid value %v for argument %s of rule %s, expected []string value got %T", v, k, r.Name(), v)
		}
		return r.buildIgnoreRegex(ignoredDirs)
	}

	return nil
}

func (r *PackageDirectoryMismatchRule) buildIgnoreRegex(ignoredDirs []string) error {
	if len(ignoredDirs) == 0 {
		r.ignoredDirs = nil
		return nil
	}

	patterns := make([]string, len(ignoredDirs))
	for i, dir := range ignoredDirs {
		patterns[i] = regexp.QuoteMeta(dir)
	}
	pattern := "(" + strings.Join(patterns, "|") + ")"

	var err error
	r.ignoredDirs, err = regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex for ignored directories: %w", err)
	}

	return nil
}

// skipDirs contains directory names that should be unconditionally ignored when checking.
// These entries handle edge cases where filepath.Base might return these values.
var skipDirs = map[string]struct{}{
	".": {}, // Current directory
	"/": {}, // Root directory
	"":  {}, // Empty path
}

// Apply applies the rule to the given file.
func (r *PackageDirectoryMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsMain() {
		return nil
	}

	dirPath := filepath.Dir(file.Name)
	dirName := filepath.Base(dirPath)

	if r.ignoredDirs != nil && r.ignoredDirs.MatchString(dirPath) {
		return nil
	}

	packageName := file.AST.Name.Name
	normalizedDirName := normalizePath(dirName)
	normalizedPackageName := normalizePath(packageName)

	// Check if we got an invalid directory.
	if _, skipDir := skipDirs[dirName]; skipDir {
		return nil
	}

	// Files directly in 'internal/' (like 'internal/abcd.go') should not be checked.
	// But files in subdirectories of 'internal/' (like 'internal/foo/abcd.go') should be checked.
	if dirName == "internal" {
		return nil
	}

	if normalizedDirName == normalizedPackageName {
		return nil
	}

	if file.IsTest() {
		// External test package (directory + '_test' suffix)
		if packageName == normalizedDirName+"_test" {
			return nil
		}
	}

	failure := ""

	// For version directories (v1, v2, etc.), we need to check also the parent directory
	if isVersionPath(dirName) {
		parentDirName := filepath.Base(filepath.Dir(dirPath))
		normalizedParentDirName := normalizePath(parentDirName)

		if normalizedPackageName == normalizedParentDirName {
			return nil
		}

		failure = fmt.Sprintf("package name %q does not match directory name %q or parent directory name %q", packageName, dirName, parentDirName)
	}

	// When no specific message was place, use default one
	if failure == "" {
		failure = fmt.Sprintf("package name %q does not match directory name %q", packageName, dirName)
	}

	return []lint.Failure{
		{
			Failure:    failure,
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
