package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseSlicesSort spots calls to sort.* that can be replaced by [slices] package methods.
type UseSlicesSort struct{}

// Apply applies the rule to given file.
func (*UseSlicesSort) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go121) {
		return nil // nothing to do, the package slices was added in version 1.21
	}

	var failures []lint.Failure

	walker := lintSort{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseSlicesSort) Name() string {
	return "use-slices-sort"
}

type lintSort struct {
	onFailure func(lint.Failure)
}

func (w lintSort) Visit(n ast.Node) ast.Visitor {
	funcCall, ok := n.(*ast.CallExpr)
	if !ok {
		return w // not a function call
	}

	replacement, found := sortReplacement(funcCall.Fun)
	if !found {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryMaintenance,
		Node:       n,
		Confidence: 1,
		Failure:    fmt.Sprintf("replace sort.%s by slices.%s", replacement.sortMethod, replacement.slicesMethod),
	})

	return nil
}

// sortReplacement returns true if the given function call is a call to a sort method that
// can be replaced with a call to a slices package method, false otherwise.
// Alongside with the boolean, the function returns the sort method name and the name of its
// replacement method from the slices package.
func sortReplacement(expr ast.Expr) (struct{ sortMethod, slicesMethod string }, bool) {
	var result struct{ sortMethod, slicesMethod string }
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return struct{ sortMethod, slicesMethod string }{}, false
	}

	if !astutils.IsIdent(sel.X, "sort") {
		return struct{ sortMethod, slicesMethod string }{}, false
	}

	result.sortMethod = sel.Sel.Name
	switch result.sortMethod {
	case "Float64s", "Ints", "Strings":
		result.slicesMethod = "Sort"
	case "Slice", "Sort":
		result.slicesMethod = "SortFunc"
	case "SliceStable", "Stable":
		result.slicesMethod = "SortStableFunc"
	case "Float64sAreSorted", "IntsAreSorted", "StringsAreSorted":
		result.slicesMethod = "IsSorted"
	case "IsSorted", "SliceIsSorted":
		result.slicesMethod = "IsSortedFunc"
	default:
		return struct{ sortMethod, slicesMethod string }{}, false
	}

	return result, true
}
