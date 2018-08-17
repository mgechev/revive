package rule

import (
	"github.com/mgechev/revive/lint"
	"go/ast"
	"go/token"
)

// BoolCompRule warns on suspicious or overcomplex boolean comparisons.
type BoolCompRule struct{}

// Apply applies the rule to given file.
func (r *BoolCompRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	astFile := file.AST
	w := &lintBoolComp{astFile, onFailure}
	ast.Walk(w, astFile)
	return failures
}

// Name returns the rule name.
func (r *BoolCompRule) Name() string {
	return "bool-comp"
}

type lintBoolComp struct {
	file      *ast.File
	onFailure func(lint.Failure)
}

func (w *lintBoolComp) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if !isBoolOp(n.Op) {
			return w
		}

		w.checkBinaryExpr(n)
	}

	return w
}

const (
	trueName  = "true"
	falseName = "false"
)

func isBoolOp(t token.Token) bool {
	switch t {
	case token.LAND, token.LOR, token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		return true
	}

	return false
}

func isExprABooleanLit(n ast.Node) bool {
	oper, ok := n.(*ast.Ident)
	return ok && (oper.Name == trueName || oper.Name == falseName)
}

func (w lintBoolComp) checkBinaryExpr(be *ast.BinaryExpr) {
	if isExprABooleanLit(be.X) || isExprABooleanLit(be.Y) {
		w.newFailure(be, "omit comparison with boolean constants", "style")
	}

	if w.areEqual(be.X, be.Y) {
		w.newFailure(be, "operands are the same on both sides of the binary expression", "logic")
	}
}

func (w lintBoolComp) areEqual(x, y ast.Expr) bool {
	// just check the simple case when both expressions are identifiers
	left, ok1 := x.(*ast.Ident)
	right, ok2 := y.(*ast.Ident)

	return ok1 && ok2 && left.Name == right.Name
}

func (w lintBoolComp) newFailure(node ast.Node, msg string, cat string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   cat,
		Failure:    msg,
	})
}
