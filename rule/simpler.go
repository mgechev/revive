package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// SimplerRule lints code to suggest simplifications.
type SimplerRule struct{}

// Apply applies the rule to given file.
func (r *SimplerRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	var branchingFunctions = map[string]map[string]bool{
		"os": map[string]bool{"Exit": true},
		"log": map[string]bool{
			"Fatal":   true,
			"Fatalf":  true,
			"Fatalln": true,
			"Panic":   true,
			"Panicf":  true,
			"Panicln": true,
		},
	}

	w := lintSimplerRule{onFailure, branchingFunctions}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *SimplerRule) Name() string {
	return "simpler"
}

type lintSimplerRule struct {
	onFailure          func(lint.Failure)
	branchingFunctions map[string]map[string]bool
}

const (
	trueName  = "true"
	falseName = "false"
)

func (w lintSimplerRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if n.Op != token.EQL && n.Op != token.NEQ {
			return w
		}

		w.checkBinaryExpr(n)

	case *ast.FuncDecl:
		if n.Body == nil || n.Type.Results != nil {
			return w
		}
		stmts := n.Body.List
		if len(stmts) == 0 {
			return w
		}

		lastStmt := stmts[len(stmts)-1]
		rs, ok := lastStmt.(*ast.ReturnStmt)
		if !ok {
			return w
		}

		if len(rs.Results) == 0 {
			w.newFailure(lastStmt, "omit unnecessary return statement")
		}

	case *ast.SwitchStmt:
		w.checkSwitchBody(n.Body)
	case *ast.TypeSwitchStmt:
		w.checkSwitchBody(n.Body)
	case *ast.CaseClause:
		if n.Body == nil {
			return w
		}
		stmts := n.Body
		if len(stmts) == 0 {
			return w
		}

		lastStmt := stmts[len(stmts)-1]
		rs, ok := lastStmt.(*ast.BranchStmt)
		if !ok {
			return w
		}

		if rs.Tok == token.BREAK && rs.Label == nil {
			w.newFailure(lastStmt, "omit unnecessary break at the end of case clause")
		}
	}

	return w
}

func (w lintSimplerRule) checkSwitchBody(b *ast.BlockStmt) {
	cases := b.List
	if len(cases) != 1 {
		return
	}

	_, ok := cases[0].(*ast.CaseClause)
	if !ok {
		return
	}

	w.newFailure(b, "switch with only one case can be replaced by an if-then")
}

func isExprABooleanLit(n ast.Node) bool {
	oper, ok := n.(*ast.Ident)
	return ok && (oper.Name == trueName || oper.Name == falseName)
}

func (w lintSimplerRule) checkBinaryExpr(be *ast.BinaryExpr) {
	if isExprABooleanLit(be.X) {
		w.newFailure(be, "omit comparison with boolean constants")
	}

	if isExprABooleanLit(be.Y) {
		w.newFailure(be, "omit comparison with boolean constants")
	}
}

func (w lintSimplerRule) newFailure(node ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   "style",
		Failure:    msg,
	})
}
