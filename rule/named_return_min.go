package rule

import (
	"fmt"
	"github.com/mgechev/revive/lint"
	"go/ast"
)

const (
	message            = "when a function has more than %d return values, only one should be named"
	defaultMinNonNamed = 2
)

// NamedReturnMinRule lints functions that have more than two return values and more than one named return value.
type NamedReturnMinRule struct {
	minNonNamed int
}

func (r *NamedReturnMinRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.minNonNamed = defaultMinNonNamed
		return nil
	}

	err := checkNumberOfArguments(1, arguments, r.Name())
	if err != nil {
		return err
	}

	maxNonNamed, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		return nil
	}

	if maxNonNamed < 1 {
		return fmt.Errorf("the argument of rule %q must be a positive integer, got %d", r.Name(), maxNonNamed)
	}

	r.minNonNamed = int(maxNonNamed)
	return nil
}

func (r *NamedReturnMinRule) Name() string {
	return "named-return-min"
}

// Apply applies the rule to given file.
func (r *NamedReturnMinRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsMain() || file.IsTest() {
		return nil
	}

	var failures []lint.Failure
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Type == nil || fn.Type.Results == nil || fn.Type.Results.List == nil || len(fn.Type.Results.List) < r.minNonNamed {
			continue // not a function or not enough return values
		}

		fnArgs := fn.Type.Results.List
		for _, arg := range fnArgs {
			if arg.Names == nil || len(arg.Names) == 0 {
				failures = append(failures, lint.Failure{
					Node:       arg,
					Category:   lint.FailureCategoryStyle,
					Failure:    fmt.Sprintf(message, r.minNonNamed),
					Confidence: 0.9,
				})
				break
			}
		}
	}

	return failures
}
