package rule

import (
	"bytes"
	"fmt"
	"github.com/mgechev/revive/lint"
	"go/ast"
	"go/format"
	"go/token"
)

// SuspiciousLogicalExprRule warns on suspicious logical expressions.
type SuspiciousLogicalExprRule struct{}

// Apply applies the rule to given file.
func (r *SuspiciousLogicalExprRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	astFile := file.AST
	w := &lintSuspiciousLogicalExpr{astFile, onFailure}
	ast.Walk(w, astFile)
	return failures
}

// Name returns the rule name.
func (r *SuspiciousLogicalExprRule) Name() string {
	return "suspicious-logical-expr"
}

type lintSuspiciousLogicalExpr struct {
	file      *ast.File
	onFailure func(lint.Failure)
}

func (w *lintSuspiciousLogicalExpr) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if !w.isOperatorWithLogicalResult(n.Op) {
			return w
		}

		if !w.areEqual(n.X, n.Y) {
			return w
		}

		if n.Op == token.EQL {
			w.newFailure(n, "expression always evaluates to true")
			return w
		}

		if w.isInequalityOperator(n.Op) {
			w.newFailure(n, "expression always evaluates to false")
			return w
		}

		w.newFailure(n, "left and right hand-side sub-expressions are the same")
	}

	return w
}

func (w *lintSuspiciousLogicalExpr) isOperatorWithLogicalResult(t token.Token) bool {
	switch t {
	case token.LAND, token.LOR, token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (w *lintSuspiciousLogicalExpr) isInequalityOperator(t token.Token) bool {
	switch t {
	case token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (w lintSuspiciousLogicalExpr) areEqual(x, y ast.Expr) bool {
	fset := token.NewFileSet()
	var buf1 bytes.Buffer
	if err := format.Node(&buf1, fset, x); err != nil {
		return false // keep going in case of error
	}

	var buf2 bytes.Buffer
	if err := format.Node(&buf2, fset, y); err != nil {
		return false // keep going in case of error
	}

	return fmt.Sprintf("%s", buf1.Bytes()) == fmt.Sprintf("%s", buf2.Bytes())
}

func (w lintSuspiciousLogicalExpr) newFailure(node ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   "logic",
		Failure:    msg,
	})
}
