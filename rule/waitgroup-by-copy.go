package rule

import (
	"github.com/mgechev/revive/lint"
	"go/ast"
)

// WaitGroupByCopyRule lints sync.WaitGroup passed by copy in functions.
type WaitGroupByCopyRule struct{}

// Apply applies the rule to given file.
func (r *WaitGroupByCopyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintWaitGroupByCopyRule{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *WaitGroupByCopyRule) Name() string {
	return "waitgroup-by-copy"
}

type lintWaitGroupByCopyRule struct {
	onFailure func(lint.Failure)
}

func (w lintWaitGroupByCopyRule) Visit(node ast.Node) ast.Visitor {
	// look for function declarations
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	// Check all function's parameters
	for _, field := range fd.Type.Params.List {
		if !w.isWaitGroup(field.Type) {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       field,
			Failure:    "sync.WaitGroup passed by value, the function will get a copy of the original one",
		})
	}

	return nil
}

func (lintWaitGroupByCopyRule) isWaitGroup(ft ast.Expr) bool {
	se, ok := ft.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, _ := se.X.(*ast.Ident)
	sel := se.Sel.Name
	return x.Name == "sync" && sel == "WaitGroup"
}
