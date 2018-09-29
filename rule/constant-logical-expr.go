package rule

import (
	"bytes"
	"fmt"
	"github.com/mgechev/revive/lint"
	"go/ast"
	"go/format"
	"go/token"
)

// ConstantLogicalExprRule warns on constant logical expressions.
type ConstantLogicalExprRule struct{}

// Apply applies the rule to given file.
func (r *ConstantLogicalExprRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	astFile := file.AST
	w := &lintConstantLogicalExpr{astFile, onFailure}
	ast.Walk(w, astFile)
	return failures
}

// Name returns the rule name.
func (r *ConstantLogicalExprRule) Name() string {
	return "constant-logical-expr"
}

type lintConstantLogicalExpr struct {
	file      *ast.File
	onFailure func(lint.Failure)
}

func (w *lintConstantLogicalExpr) Visit(node ast.Node) ast.Visitor {
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

func (w *lintConstantLogicalExpr) isOperatorWithLogicalResult(t token.Token) bool {
	switch t {
	case token.LAND, token.LOR, token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (w *lintConstantLogicalExpr) isInequalityOperator(t token.Token) bool {
	switch t {
	case token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func (w lintConstantLogicalExpr) areEqual(x, y ast.Expr) bool {
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

func (w lintConstantLogicalExpr) newFailure(node ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   "logic",
		Failure:    msg,
	})
}
