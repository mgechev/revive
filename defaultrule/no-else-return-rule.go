package defaultrule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// LintElseRule lints given else constructs.
type LintElseRule struct{}

// Apply applies the rule to given file.
func (r *LintElseRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	var failures []rule.Failure

	onFailure := func(failure rule.Failure) {
		failures = append(failures, failure)
	}

	w := lintElse{make(map[*ast.IfStmt]bool), onFailure}
	ast.Walk(w, file.GetAST())
	return failures
}

// Name returns the rule name.
func (r *LintElseRule) Name() string {
	return "no-else-return"
}

type lintElse struct {
	ignore    map[*ast.IfStmt]bool
	onFailure func(rule.Failure)
}

func (w lintElse) Visit(node ast.Node) ast.Visitor {
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
	lastStmt := ifStmt.Body.List[len(ifStmt.Body.List)-1]
	if _, ok := lastStmt.(*ast.ReturnStmt); ok {
		extra := ""
		if shortDecl {
			extra = " (move short variable declaration to its own line if necessary)"
		}
		w.onFailure(rule.Failure{
			Failure: "if block ends with a return statement, so drop this else and outdent its block" + extra,
			Type:    rule.FailureTypeWarning,
			Node:    ifStmt.Else,
		})
	}
	return w
}
