package defaultrules

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/golinter/file"
	"github.com/mgechev/golinter/rules"
	"github.com/mgechev/golinter/visitors"
)

const ruleName = "no-else-return"

// LintElseRule lints given else constructs.
type LintElseRule struct {
	rules.Rule
}

// Apply applies the rule to given file.
func (r *LintElseRule) Apply(file *file.File, arguments rules.RuleArguments) []rules.Failure {
	res := &lintElseVisitor{}
	visitors.Setup(res, rules.RuleConfig{Name: ruleName, Arguments: arguments}, file)
	res.Visit(file.GetAST())
	return res.GetFailures()
}

type lintElseVisitor struct {
	visitors.RuleVisitor
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
		w.AddFailure(rules.Failure{
			Failure:  "if block ends with a return statement, so drop this else and outdent its block",
			Type:     rules.FailureTypeWarning,
			Position: w.GetPosition(node.Else.Pos(), node.Else.End()),
		})
	}
}
