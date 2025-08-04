package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// UselessFallthroughRule warns on useless fallthroughts in switch case clauses.
type UselessFallthroughRule struct{}

// Apply applies the rule to given file.
func (*UselessFallthroughRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUselessFallthrough{onFailure: onFailure}
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
func (*UselessFallthroughRule) Name() string {
	return "useless-fallthrough"
}

type lintUselessFallthrough struct {
	onFailure func(lint.Failure)
}

func (w *lintUselessFallthrough) Visit(node ast.Node) ast.Visitor {
	switchStmt, ok := node.(*ast.SwitchStmt)
	if !ok { // not a switch statement, keep walking the AST
		return w
	}

	if switchStmt.Tag == nil {
		return w // Not interested in un-tagged switches
	}

	casesCount := len(switchStmt.Body.List)
	for i := 0; i < casesCount-1; i++ {
		caseClause := switchStmt.Body.List[i].(*ast.CaseClause)
		caseBody := caseClause.Body

		if len(caseBody) != 1 {
			continue // skip body if is not just only one statement
		}

		branchStmt, ok := caseBody[0].(*ast.BranchStmt)
		if !ok || branchStmt.Tok != token.FALLTHROUGH {
			continue // not a fallthrought
		}

		confidence := 1.0
		if nextCaseClause := switchStmt.Body.List[i+1].(*ast.CaseClause); nextCaseClause.List == nil {
			// the next case clause is the default clause, report with lower confidence.
			confidence = 0.8
		}

		w.onFailure(lint.Failure{
			Confidence: confidence,
			Node:       branchStmt,
			Category:   lint.FailureCategoryCodeStyle,
			Failure:    `this "fallthrough" can be removed by consolidating this case clause with the next one`,
		})

		ast.Walk(w, caseClause)
	}

	return nil
}
