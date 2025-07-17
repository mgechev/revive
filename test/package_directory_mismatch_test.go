package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestPackageDirectoryMismatch(t *testing.T) {
	testRule(t, "package_directory_mismatch/incorrect/mismatch", &rule.PackageDirectoryMismatchRule{})
	testRule(t, "package_directory_mismatch/correct/match", &rule.PackageDirectoryMismatchRule{})

	testRule(t, "package_directory_mismatch/cmd/main", &rule.PackageDirectoryMismatchRule{})

	testRule(t, "package_directory_mismatch/mixed/correct", &rule.PackageDirectoryMismatchRule{})
	testRule(t, "package_directory_mismatch/mixed/incorrect", &rule.PackageDirectoryMismatchRule{})

	testRule(t, "package_directory_mismatch/utils/file1_test", &rule.PackageDirectoryMismatchRule{})

	testRule(t, "package_directory_mismatch/internal/utils", &rule.PackageDirectoryMismatchRule{})
	testRule(t, "package_directory_mismatch/testdata/fixtures", &rule.PackageDirectoryMismatchRule{})
}
