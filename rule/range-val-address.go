package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// RangeValAddress lints
type RangeValAddress struct{}

// Apply applies the rule to given file.
func (r *RangeValAddress) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := rangeValAddress{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *RangeValAddress) Name() string {
	return "range-val-address"
}

type rangeValAddress struct {
	onFailure func(lint.Failure)
}

func (w rangeValAddress) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.RangeStmt)
	if ok {
		rangeValue := w.getNameFromExpr(n.Value)
		if rangeValue == "" {
			return w
		}

		fselect := func(n ast.Node) bool {
			asgmt, ok := n.(*ast.AssignStmt)
			if ok {
				for _, exp := range asgmt.Lhs {
					e, ok := exp.(*ast.IndexExpr)
					if ok {
						u, ok := e.Index.(*ast.UnaryExpr) // e.g. a[&value]...
						if ok && u.Op == token.AND && w.getNameFromExpr(u.X) == rangeValue {
							return true
						}
					}
				}

				for _, exp := range asgmt.Rhs {
					switch e := exp.(type) {
					case *ast.UnaryExpr: // e.g. ...&value
						if e.Op == token.AND && w.getNameFromExpr(e.X) == rangeValue {
							return true
						}
					case *ast.CallExpr: // e.g. ...append(arr, &value)
						for _, v := range e.Args {
							u, ok := v.(*ast.UnaryExpr)
							if ok && u.Op == token.AND && w.getNameFromExpr(u.X) == rangeValue {
								return true
							}
						}
					}
				}
			}
			return false
		}

		assignmentsToReceiver := pick(n.Body, fselect, nil)
		for _, assignment := range assignmentsToReceiver {
			w.onFailure(lint.Failure{
				Node:       assignment,
				Confidence: 1,
				Failure:    "suspicious assignment in range-loop. variables always have the same address",
			})
		}
	}
	return w
}

func (rangeValAddress) getNameFromExpr(ie ast.Expr) string {
	ident, ok := ie.(*ast.Ident)
	if !ok {
		return ""
	}

	return ident.Name
}
