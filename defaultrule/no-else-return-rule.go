package defaultrule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

const (
	ruleName = "no-else-return"
	failure  = "if block ends with a return statement, so drop this else and outdent its block"
)

// LintElseRule lints given else constructs.
type LintElseRule struct {
	rule.AbstractRule
}

// Apply applies the rule to given file.
func (r *LintElseRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	r.File = file
	ast.Walk(lintElse{r}, file.GetAST())
	return r.Failures()
}

// Name returns the rule name.
func (r *LintElseRule) Name() string {
	return ruleName
}

type lintElse struct {
	r rule.Rule
}

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
			f.r.AddFailures(rule.Failure{
				Failure:  failure,
				Type:     rule.FailureTypeWarning,
				Position: f.r.Position(node.Else.Pos(), node.Else.End()),
			})
			return f
		}
	}
	return f
}
