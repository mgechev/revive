package rule

import (
	"testing"

	"github.com/mgechev/revive/rule"
	"github.com/mgechev/revive/testutil"
)

func TestPackageCommentsRule(t *testing.T) {
	t.Parallel()

	program := `
	/*
	Package foo is pretty sweet.
	*/
	     
	package [@f]foo[/@f]

	func foo(a int, b int, c int) {
		return a + b + c;
	}
	`
	testutil.AssertFailures(t, program, &PackageCommentsRule{}, rule.Arguments{})
}

func TestPackageCommentsRule_Success(t *testing.T) {
	t.Parallel()

	program := `
	// Package foo is awesome
	package foo

	func foo(a int, b int, c int) {
		return a + b + c;
	}
	`
	testutil.AssertSuccess(t, program, &PackageCommentsRule{}, rule.Arguments{})
}
