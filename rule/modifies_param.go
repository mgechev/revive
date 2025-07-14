package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// ModifiesParamRule warns on assignments to function parameters.
type ModifiesParamRule struct{}

// Apply applies the rule to given file.
func (*ModifiesParamRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintModifiesParamRule{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*ModifiesParamRule) Name() string {
	return "modifies-parameter"
}

type lintModifiesParamRule struct {
	params    map[string]bool
	onFailure func(lint.Failure)
}

func retrieveParamNames(pl []*ast.Field) map[string]bool {
	result := make(map[string]bool, len(pl))
	for _, p := range pl {
		for _, n := range p.Names {
			if n.Name == "_" {
				continue
			}

			result[n.Name] = true
		}
	}
	return result
}

func (w lintModifiesParamRule) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.FuncDecl:
		w.params = retrieveParamNames(v.Type.Params.List)
	case *ast.IncDecStmt:
		if id, ok := v.X.(*ast.Ident); ok {
			checkParam(id, &w)
		}
	case *ast.AssignStmt:
		lhs := v.Lhs
		for i, e := range lhs {
			id, ok := e.(*ast.Ident)
			if ok {
				if i < len(v.Rhs) {
					if callExpr, ok := v.Rhs[i].(*ast.CallExpr); ok && isSlicesDelete(callExpr) {
						w.checkSlicesDelete(callExpr)
						continue
					}
				}
				checkParam(id, &w)
			}
		}
	case *ast.ExprStmt:
		if callExpr, ok := v.X.(*ast.CallExpr); ok && isSlicesDelete(callExpr) {
			w.checkSlicesDelete(callExpr)
		}
	}

	return w
}

func checkParam(id *ast.Ident, w *lintModifiesParamRule) {
	if w.params[id.Name] {
		w.onFailure(lint.Failure{
			Confidence: 0.5, // confidence is low because of shadow variables
			Node:       id,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("parameter '%s' seems to be modified", id),
		})
	}
}

func isSlicesDelete(callExpr *ast.CallExpr) bool {
	return astutils.IsPkgDotName(callExpr.Fun, "slices", "Delete") ||
		astutils.IsPkgDotName(callExpr.Fun, "slices", "DeleteFunc")
}

func (w *lintModifiesParamRule) checkSlicesDelete(callExpr *ast.CallExpr) {
	if len(callExpr.Args) == 0 {
		return
	}

	if id, ok := callExpr.Args[0].(*ast.Ident); ok && w.params[id.Name] {
		funcName := "function"
		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			funcName = fmt.Sprintf("%s.%s", sel.X, sel.Sel.Name)
		}

		w.onFailure(lint.Failure{
			Confidence: 1, // slices.Delete/DeleteFunc always modifies
			Node:       callExpr,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("parameter '%s' is modified by %s", id.Name, funcName),
		})
	}
}
