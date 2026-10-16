package lint_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

// typeCheckRule runs type checking on the package of each linted file and records the result.
type typeCheckRule struct {
	err error
}

var _ lint.Rule = (*typeCheckRule)(nil)

func (*typeCheckRule) Name() string { return "type-check" }

func (r *typeCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	r.err = file.Pkg.TypeCheck()
	return nil
}

// TestPackage_TypeCheckDoesNotPanic checks that go/types does not panic on nested indexed
// composite literals with invalid keys and values (see issue #59).
// Package.TypeCheck does not recover from panics, so a panic would crash the linter.
func TestPackage_TypeCheckDoesNotPanic(t *testing.T) {
	const src = `package p

var a = [][]int{
	2: {1: 1, "x": 2, 3: "y"},
	{-1: 0, 1: undefined},
}

var b = map[string][]int{
	"a": {0: 1, 0: 2, 1 << 70: 3},
	1:   {2: 3},
}

var c = [...][2][]struct{ F int }{
	1: {1: {{F: "x"}, {G: 1}}},
	5: {5: {{}}},
}
`
	rule := &typeCheckRule{}
	lintFile(t, "file.go", src, []lint.Rule{rule}, lint.Config{Rules: lint.RulesConfig{rule.Name(): {}}}, nil)

	// The source is invalid, so type checking reports an error, but it must not panic.
	if rule.err == nil {
		t.Error("expected type checking error for invalid source")
	}
}
