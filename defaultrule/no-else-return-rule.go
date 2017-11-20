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
	ast.Walk(lintElse(func(failure rule.Failure) {
		failures = append(failures, failure)
	}), file.GetAST())
	return failures
}

// Name returns the rule name.
func (r *LintElseRule) Name() string {
	return "no-else-return"
}

type lintElse func(rule.Failure)

func (f lintElse) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.IfStmt)
	if ok {
		if node.Else == nil {
			return f
		}
		if _, ok := node.Else.(*ast.BlockStmt); !ok {
			// only care about elses without conditions
			return f
		}
		if len(node.Body.List) == 0 {
			return f
		}
		// shortDecl := false // does the if statement have a ":=" initialization statement?
		if node.Init != nil {
			if as, ok := node.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
				// shortDecl = true
			}
		}
		lastStmt := node.Body.List[len(node.Body.List)-1]
		if _, ok := lastStmt.(*ast.ReturnStmt); ok {
			f(rule.Failure{
				Failure: "if block ends with a return statement, so drop this else and outdent its block",
				Type:    rule.FailureTypeWarning,
				Node:    node.Else,
			})
			return f
		}
	}
	return f
}
