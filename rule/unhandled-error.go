package rule

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/mgechev/revive/lint"
)

// UnhandledErrorRule lints given else constructs.
type UnhandledErrorRule struct{}

// Apply applies the rule to given file.
func (r *UnhandledErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := &lintUnhandledErrors{
		pkg: file.Pkg,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck()
	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *UnhandledErrorRule) Name() string {
	return "unhandled-error"
}

type lintUnhandledErrors struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// Visit looks for statements that are function calls.
// If the called function returns a value of type error a failure will be created.
func (w *lintUnhandledErrors) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ExprStmt:
		fCall, ok := n.X.(*ast.CallExpr)
		if !ok {
			return nil // not a function call
		}

		funcType := w.pkg.TypeOf(fCall)
		if funcType == nil {
			return nil // skip, type info not available
		}

		switch t := funcType.(type) {
		case *types.Named:
			if !w.isTypeError(t) {
				return nil // func call does not return an error
			}

			w.addFailure(fCall)
		default:
			retTypes, ok := funcType.Underlying().(*types.Tuple)
			if !ok {
				return nil // skip, unable to retrieve return type of the called function
			}

			if w.returnsAnError(retTypes) {
				w.addFailure(fCall)
			}
		}
	}
	return w
}

func (w *lintUnhandledErrors) addFailure(n *ast.CallExpr) {
	w.onFailure(lint.Failure{
		Category:   "bad practice",
		Confidence: 1,
		Node:       n,
		Failure:    fmt.Sprintf("Unhandled error in call to function %v", gofmt(n.Fun)),
	})
}

func (*lintUnhandledErrors) isTypeError(t *types.Named) bool {
	const errorTypeName = "_.error"

	return t.Obj().Id() == errorTypeName
}

func (w *lintUnhandledErrors) returnsAnError(tt *types.Tuple) bool {
	for i := 0; i < tt.Len(); i++ {
		nt, ok := tt.At(i).Type().(*types.Named)
		if ok && w.isTypeError(nt) {
			return true
		}
	}
	return false
}
