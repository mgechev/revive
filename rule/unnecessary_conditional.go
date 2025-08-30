package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UnnecessaryConditionalRule warns on if...else statements with both branches being the same.
type UnnecessaryConditionalRule struct{}

// Apply applies the rule to given file.
func (*UnnecessaryConditionalRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryConditional{onFailure: onFailure}
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*UnnecessaryConditionalRule) Name() string {
	return "unnecessary-conditional"
}

type lintUnnecessaryConditional struct {
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryConditional) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	if ifStmt.Else == nil {
		return w // if without else
	}

	elseBranch, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok { // if-else-if construction, the rule only copes with single if...else statements
		return w
	}

	thenStmts := ifStmt.Body.List
	elseStmts := elseBranch.List
	if len(thenStmts) != 1 || len(elseStmts) != 1 {
		return w
	}

	replacement := ""
	switch thenStmt := thenStmts[0].(type) {
	case *ast.ReturnStmt:
		thenBool, ok := w.resultValueIsBooleanLiteral(thenStmt.Results)
		if !ok {
			return w
		}

		elseStmt, ok := elseStmts[0].(*ast.ReturnStmt)
		if !ok {
			return w
		}

		_, ok = w.resultValueIsBooleanLiteral(elseStmt.Results)
		if !ok {
			return w
		}

		cond := astutils.GoFmt(ifStmt.Cond)
		if thenBool != "true" {
			cond = "!(" + cond + ")"
		}

		replacement = "return " + cond
	case *ast.AssignStmt:
		thenBool, ok := w.isBooleanLiteral(thenStmt.Rhs)
		if !ok {
			return w
		}

		thenLhs := astutils.GoFmt(thenStmt.Lhs[0])

		elseStmt, ok := elseStmts[0].(*ast.AssignStmt)
		if !ok {
			return w
		}

		elseLhs := astutils.GoFmt(elseStmt.Lhs[0])
		if thenLhs != elseLhs {
			return w
		}

		_, ok = w.isBooleanLiteral(elseStmt.Rhs)
		if !ok {
			return w
		}

		cond := astutils.GoFmt(ifStmt.Cond)
		if thenBool != "true" {
			cond = "!(" + cond + ")"
		}

		replacement = fmt.Sprintf("%s %s %s", thenLhs, thenStmt.Tok.String(), cond)
	default:
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1.0,
		Node:       ifStmt,
		Category:   lint.FailureCategoryLogic,
		Failure:    "replace this conditional by: " + replacement,
	})

	return nil
}

func (w *lintUnnecessaryConditional) resultValueIsBooleanLiteral(results []ast.Expr) (string, bool) {
	if len(results) != 1 {
		return "", false
	}

	ident, ok := results[0].(*ast.Ident)
	if !ok {
		return "", false
	}

	return ident.Name, (ident.Name == "true" || ident.Name == "false")
}

func (*lintUnnecessaryConditional) isBooleanLiteral(exprs []ast.Expr) (string, bool) {
	if len(exprs) != 1 {
		return "", false
	}

	ident, ok := exprs[0].(*ast.Ident)
	if !ok {
		return "", false
	}

	return ident.Name, (ident.Name == "true" || ident.Name == "false")
}
