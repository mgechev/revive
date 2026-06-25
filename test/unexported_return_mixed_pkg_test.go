package test_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// TestUnexportedReturnMixedPackages verifies that the unexported-return rule produces
// stable and correct results when a directory contains both source files,
// internal test files (`package pub`) and external test files (`package pub_test`).
//
// This mirrors what dots.ResolvePackages returns at runtime: GoFiles + TestGoFiles + XTestGoFiles
// collapsed into a single []string passed to lint.Linter.Lint.
//
// Without the fix the rule:
//   - reports the wrong package qualifier (*pub_test.impl[T] instead of *pub.impl[T])
//   - occasionally drops the issue entirely due to a race on map iteration in TypeCheck
func TestUnexportedReturnMixedPackages(t *testing.T) {
	baseDir := filepath.Join("..", "testdata", "unexported_return_mixed_pkg")

	files := []string{
		filepath.Join(baseDir, "pub.go"),
		filepath.Join(baseDir, "pub_internal_test.go"),
		filepath.Join(baseDir, "pub_external_test.go"),
	}

	const expected = "exported func New returns unexported type *pub.impl[T], which can be annoying to use"

	// Repeat enough times to hit the non-deterministic path with high probability.
	const iterations = 50

	for i := range iterations {
		l := lint.New(os.ReadFile, 0)

		ps, err := l.Lint([][]string{files}, []lint.Rule{&rule.UnexportedReturnRule{}}, lint.Config{})
		if err != nil {
			t.Fatalf("iteration %d: Lint failed: %v", i, err)
		}

		var failures []lint.Failure
		for p := range ps {
			failures = append(failures, p)
		}

		// Filter out failures from test files — they're not the target of this assertion.
		var prodFailures []lint.Failure
		for _, f := range failures {
			if !strings.HasSuffix(f.Filename(), "_test.go") {
				prodFailures = append(prodFailures, f)
			}
		}

		if len(prodFailures) != 1 {
			t.Fatalf("iteration %d: expected exactly 1 failure on pub.go, got %d: %+v", i, len(prodFailures), prodFailures)
		}

		got := prodFailures[0].Failure
		if got != expected {
			t.Fatalf("iteration %d: wrong failure message\n got: %s\nwant: %s", i, got, expected)
		}
	}
}
