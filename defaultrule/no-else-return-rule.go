package defaultrule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
	"github.com/mgechev/revive/visitor"
)

const (
	ruleName = "no-else-return"
	failure  = "if block ends with a return statement, so drop this else and outdent its block"
)

// LintElseRule lints given else constructs.
type LintElseRule struct {
	rule.Rule
}

// Apply applies the rule to given file.
func (r *LintElseRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	res := &lintElseVisitor{}
	visitor.Setup(res, rule.Config{Name: ruleName, Arguments: arguments}, file)
	res.Visit(file.GetAST())
	return res.GetFailures()
}

// GetName returns the rule name.
func (r *LintElseRule) GetName() string {
	return ruleName
}

type lintElseVisitor struct {
	visitor.RuleVisitor
}

func (w *lintElseVisitor) VisitIfStmt(node *ast.IfStmt) {
	if node.Else == nil {
		return
	}
	if _, ok := node.Else.(*ast.BlockStmt); !ok {
		// only care about elses without conditions
		return
	}
	if len(node.Body.List) == 0 {
		return
	}
	// shortDecl := false // does the if statement have a ":=" initialization statement?
	if node.Init != nil {
		if as, ok := node.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
			// shortDecl = true
		}
	}
	lastStmt := node.Body.List[len(node.Body.List)-1]
	if _, ok := lastStmt.(*ast.ReturnStmt); ok {
		w.AddFailure(rule.Failure{
			Failure:  failure,
			Type:     rule.FailureTypeWarning,
			Position: w.GetPosition(node.Else.Pos(), node.Else.End()),
		})
	}
}
