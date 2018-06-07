package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// SuperfluousElse lints given else constructs.
type SuperfluousElse struct{}

// Apply applies the rule to given file.
func (r *SuperfluousElse) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintSuperfluousElse{make(map[*ast.IfStmt]bool), onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *SuperfluousElse) Name() string {
	return "superfluous-else"
}

type lintSuperfluousElse struct {
	ignore    map[*ast.IfStmt]bool
	onFailure func(lint.Failure)
}

func (w lintSuperfluousElse) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok || ifStmt.Else == nil {
		return w
	}
	if w.ignore[ifStmt] {
		return w
	}
	if elseif, ok := ifStmt.Else.(*ast.IfStmt); ok {
		w.ignore[elseif] = true
		return w
	}
	if _, ok := ifStmt.Else.(*ast.BlockStmt); !ok {
		// only care about elses without conditions
		return w
	}
	if len(ifStmt.Body.List) == 0 {
		return w
	}
	shortDecl := false // does the if statement have a ":=" initialization statement?
	if ifStmt.Init != nil {
		if as, ok := ifStmt.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
			shortDecl = true
		}
	}
	extra := ""
	if shortDecl {
		extra = " (move short variable declaration to its own line if necessary)"
	}

	lastStmt := ifStmt.Body.List[len(ifStmt.Body.List)-1]
	if stmt, ok := lastStmt.(*ast.BranchStmt); ok {
		token := stmt.Tok.String()
		if token != "fallthrough" {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ifStmt.Else,
				Category:   "indent",
				URL:        "#indent-error-flow",
				Failure:    "if block ends with a " + token + " statement, so drop this else and outdent its block" + extra,
			})
		}
	}

	return w
}
