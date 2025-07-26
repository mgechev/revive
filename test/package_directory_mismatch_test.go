package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageDirectoryMismatch(t *testing.T) {
	// Configure rule to ignore no directories (empty string) so our tests in 'testdata/' can run
	config := &lint.RuleConfig{Arguments: []any{""}}

	testRule(t, "package_directory_mismatch/good/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/bad/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/maincmd/main", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/mixed/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/mixed/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/test/good_test", &rule.PackageDirectoryMismatchRule{}, config)

	// Test normalization cases
	testRule(t, "package_directory_mismatch/normalization/fo-ob_ar/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/foo_bar/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/foo.b_ar/good", &rule.PackageDirectoryMismatchRule{}, config)

	// Test version directories (v1, v2, etc.)
	testRule(t, "package_directory_mismatch/api/v1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/V1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/v1v/api", &rule.PackageDirectoryMismatchRule{}, config)

	// Test internal directory variations
	testRule(t, "package_directory_mismatch/internal/good/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/bad/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/any", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/api/v1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/api/v2/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/v1/api", &rule.PackageDirectoryMismatchRule{}, config)
}

func TestPackageDirectoryMismatchWithDefaultConfig(t *testing.T) {
	// Test with default configuration (should ignore testdata directories by default)
	// This test verifies that files in testdata directories are ignored by default
	testRule(t, "package_directory_mismatch/testdata/ignored", &rule.PackageDirectoryMismatchRule{})
}

func TestPackageDirectoryMismatchWithTestPrefixRegex(t *testing.T) {
	// Test with regex that excludes everything with "test"
	config := &lint.RuleConfig{Arguments: []any{"test"}}

	testRule(t, "package_directory_mismatch/testinfo/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/testutils/good", &rule.PackageDirectoryMismatchRule{}, config)
}

func TestPackageDirectoryMismatchWithMultipleDirectories(t *testing.T) {
	// Test with regex that excludes both "testcases" and "testinfo" specifically
	config := &lint.RuleConfig{Arguments: []any{"(testutils|testinfo)"}}

	testRule(t, "package_directory_mismatch/testutils/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/testinfo/good", &rule.PackageDirectoryMismatchRule{}, config)
}
