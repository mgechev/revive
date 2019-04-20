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

	fileAst := file.AST
	walker := &lintUnhandledErrors{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck()
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *UnhandledErrorRule) Name() string {
	return "unhandled-error"
}

type lintUnhandledErrors struct {
	fileAst   *ast.File
	file      *lint.File
	lastGen   *ast.GenDecl
	onFailure func(lint.Failure)
}

func (w *lintUnhandledErrors) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ExprStmt:
		x := n.X
		fCall, ok := x.(*ast.CallExpr)
		if !ok {
			return nil // not a function call
		}

		funcType := w.file.Pkg.TypeOf(fCall)
		if funcType == nil {
			return nil // skip, type info not available
		}

		switch t := funcType.(type) {
		case *types.Named:
			if t.String() != "error" {
				return nil // func call does not return an error
			}

			w.addFailure(fCall)
		default:
			retTypes, ok := funcType.Underlying().(*types.Tuple)
			if !ok {
				return nil // skip, unable to retrieve return type of the called function
			}

			if returnsAnError(retTypes) {
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
		Failure:    fmt.Sprintf("Unhandled error in call to function %v", n.Fun),
	})
}

func returnsAnError(types *types.Tuple) bool {
	for i := 0; i < types.Len(); i++ {
		if types.At(i).Type().String() == "error" {
			return true
		}
	}
	return false
}
