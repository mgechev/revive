package rule

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgechev/revive/lint"
)

// PackageDirectoryMismatchRule detects when package name doesn't match directory name.
type PackageDirectoryMismatchRule struct {
	ignoredDirs *regexp.Regexp
	workDir     string // working directory at configuration time, ignored directories are matched relative to it
}

const defaultIgnoredDirs = "testdata"

var _ lint.ConfigurableRule = (*PackageDirectoryMismatchRule)(nil)

// Configure the rule to exclude certain directories.
func (r *PackageDirectoryMismatchRule) Configure(arguments lint.Arguments) error {
	r.ignoredDirs = nil

	var err error
	r.workDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory for rule %s: %w", r.Name(), err)
	}

	if len(arguments) < 1 {
		r.ignoredDirs, err = r.buildIgnoreRegex([]string{defaultIgnoredDirs})
		return err
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument type: expected map[string]any, got %T", arguments[0])
	}

	for k, v := range args {
		if !isRuleOption(k, "ignoreDirectories") {
			return fmt.Errorf("unknown argument %s for %s rule", k, r.Name())
		}

		ignoredAny, ok := v.([]any)
		if !ok {
			return fmt.Errorf("invalid value %v for argument %s of rule %s, expected []string got %T", v, k, r.Name(), v)
		}

		ignoredDirs := make([]string, len(ignoredAny))
		for i, item := range ignoredAny {
			str, ok := item.(string)
			if !ok {
				return fmt.Errorf("invalid value in %s argument of rule %s: expected string, got %T", k, r.Name(), item)
			}
			ignoredDirs[i] = str
		}

		if len(ignoredDirs) == 0 {
			r.ignoredDirs = nil
			return nil
		}

		r.ignoredDirs, err = r.buildIgnoreRegex(ignoredDirs)
		return err
	}

	return nil
}

// buildIgnoreRegex builds a regular expression matching the given directories as sequences of whole path components.
// Entries are relative to the working directory and normalized: "testdata/", "./testdata" and "testdata" are equivalent,
// and `internal\gen` is "internal/gen". Absolute entries and entries with ".." are rejected.
func (r *PackageDirectoryMismatchRule) buildIgnoreRegex(ignoredDirs []string) (*regexp.Regexp, error) {
	patterns := make([]string, 0, len(ignoredDirs))
	for _, orig := range ignoredDirs {
		slashed := strings.ReplaceAll(orig, `\`, "/")
		if strings.HasPrefix(slashed, "/") || filepath.IsAbs(orig) {
			return nil, fmt.Errorf("invalid ignored directory %q of rule %s: absolute paths are not supported", orig, r.Name())
		}

		var components []string
		for component := range strings.SplitSeq(slashed, "/") {
			switch component {
			case "", ".":
			case "..":
				return nil, fmt.Errorf("invalid ignored directory %q of rule %s: %q is not supported", orig, r.Name(), "..")
			default:
				components = append(components, regexp.QuoteMeta(component))
			}
		}
		if len(components) == 0 {
			return nil, fmt.Errorf("invalid empty ignored directory %q of rule %s", orig, r.Name())
		}
		patterns = append(patterns, strings.Join(components, "/"))
	}
	// For ["internal/gen"] matches "internal/gen", "a/internal/gen/b" but not "internal/generated" or "myinternal/gen".
	pattern := "(^|/)(" + strings.Join(patterns, "|") + ")(/|$)"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex for ignored directories: %w", err)
	}

	return regex, nil
}

// isIgnoredDir reports whether the directory is, or is below, an ignored directory.
// The path is matched relative to the working directory without leading "..",
// so names of the directories above the working directory never match.
func (r *PackageDirectoryMismatchRule) isIgnoredDir(absDirPath string) bool {
	if r.ignoredDirs == nil {
		return false
	}

	rel, err := filepath.Rel(r.workDir, absDirPath)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	for rel == ".." || strings.HasPrefix(rel, "../") {
		rel = strings.TrimPrefix(strings.TrimPrefix(rel, ".."), "/")
	}

	return r.ignoredDirs.MatchString(rel)
}

// skipDirs contains directory names that should be unconditionally ignored when checking.
// These entries handle edge cases where [filepath.Base] might return these values.
var skipDirs = map[string]struct{}{
	".": {}, // Current directory
	"/": {}, // Root directory
	"":  {}, // Empty path
}

// semanticallyEqual checks if package and directory names are semantically equal to each other.
func (*PackageDirectoryMismatchRule) semanticallyEqual(packageName, dirName string) bool {
	normDir := normalizePath(dirName)
	normPkg := normalizePath(packageName)
	return normDir == normPkg || normDir == "go"+normPkg
}

// Apply applies the rule to the given file.
func (r *PackageDirectoryMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsMain() {
		return nil
	}

	absPath, err := filepath.Abs(file.Name)
	if err != nil {
		return nil
	}

	dirPath := filepath.Dir(absPath)
	dirName := filepath.Base(dirPath)

	if r.isIgnoredDir(dirPath) {
		return nil
	}

	// Check if we got an invalid directory.
	if _, skipDir := skipDirs[dirName]; skipDir {
		return nil
	}

	// Files directly in 'internal/' (like 'internal/abcd.go') should not be checked.
	// But files in subdirectories of 'internal/' (like 'internal/foo/abcd.go') should be checked.
	if dirName == "internal" {
		return nil
	}

	packageName := file.AST.Name.Name

	if r.semanticallyEqual(packageName, dirName) {
		return nil
	}

	if isRootDir(dirPath) {
		return nil
	}

	if file.IsTest() {
		// treat main_test differently because it's a common package name for tests
		if packageName == "main_test" {
			return nil
		}
		// External test package (directory + '_test' suffix)
		if r.semanticallyEqual(packageName, dirName+"_test") {
			return nil
		}
	}

	// define a default failure message
	failure := fmt.Sprintf("package name %q does not match directory name %q", packageName, dirName)

	// For version directories (v1, v2, etc.), we need to check also the parent directory
	if isVersionPath(dirName) {
		parentDirName := filepath.Base(filepath.Dir(dirPath))
		if r.semanticallyEqual(packageName, parentDirName) {
			return nil
		}

		if file.IsTest() {
			// External test package (directory + '_test' suffix)
			if r.semanticallyEqual(packageName, parentDirName+"_test") {
				return nil
			}
		}

		failure = fmt.Sprintf("package name %q does not match directory name %q or parent directory name %q", packageName, dirName, parentDirName)
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

// isRootDir checks if the given directory contains go.mod or .git, indicating it's a root directory.
func isRootDir(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, e := range entries {
		switch e.Name() {
		case "go.mod", ".git":
			return true
		}
	}

	return false
}

// Name returns the rule name.
func (*PackageDirectoryMismatchRule) Name() string {
	return "package-directory-mismatch"
}
