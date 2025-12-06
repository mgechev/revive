package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseSlicesSort spots calls to sort.* that can be replaced by slices.Sort
type UseSlicesSort struct{}

// Apply applies the rule to given file.
func (*UseSlicesSort) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

	isSortSort := astutils.IsPkgDotName(funcCall.Fun, "sort", "Sort")
	if isSortSort {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryMaintenance,
			Node:       n,
			Confidence: 1,
			Failure:    "replace sort.Sort by slices.SortFunc",
		})
		return nil
	}

	isSortType, basicType := isSortType(funcCall.Fun)
	if isSortType {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryMaintenance,
			Node:       n,
			Confidence: 1,
			Failure:    fmt.Sprintf("replace sort.%s by slices.Sort", basicType),
		})
		return nil
	}

	return w
}

func isSortType(expr ast.Expr) (bool, string) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false, ""
	}

	if !astutils.IsIdent(sel.X, "sort") {
		return false, ""
	}

	switch sel.Sel.Name {
	case "Float64s", "Ints", "Strings":
		return true, sel.Sel.Name
	default:
		return false, ""
	}
}
