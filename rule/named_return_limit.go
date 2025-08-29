package rule

import (
	"github.com/mgechev/revive/lint"
	"go/ast"
)

const (
	message = "when a function has more than two return values, only one should be named"
)

// ReturnLimitNamedRule lints functions that have more than two return values and more than one named return value.
type ReturnLimitNamedRule struct{}

func (r *ReturnLimitNamedRule) Name() string {
	return "return-limit-named"
}

// Apply applies the rule to given file.
func (r *ReturnLimitNamedRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsMain() || file.IsTest() {
		return nil
	}

	var failures []lint.Failure
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Type == nil || fn.Type.Results == nil || fn.Type.Results.List == nil || len(fn.Type.Results.List) <= 2 {
			continue // only care about functions with more than 2 return values
		}

		fnArgs := fn.Type.Results.List
		for _, arg := range fnArgs {
			if arg.Names == nil || len(arg.Names) == 0 {
				failures = append(failures, lint.Failure{
					Node:       arg,
					Category:   lint.FailureCategoryStyle,
					Failure:    message,
					Confidence: 0.9,
				})
				break
			}
		}
	}

	return failures
}
