package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// UselessFallthroughRule warns on useless fallthroughs in switch case clauses.
type UselessFallthroughRule struct{}

// Apply applies the rule to given file.
func (*UselessFallthroughRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	commentsMap := file.CommentMap()

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUselessFallthrough{onFailure: onFailure, commentsMap: commentsMap}
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
	onFailure   func(lint.Failure)
	commentsMap ast.CommentMap
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
	for i := range casesCount - 1 {
		caseClause := switchStmt.Body.List[i].(*ast.CaseClause)
		caseBody := caseClause.Body

		if len(caseBody) != 1 {
			continue // skip if body is not exactly one statement
		}

		branchStmt, ok := caseBody[0].(*ast.BranchStmt)
		if !ok || branchStmt.Tok != token.FALLTHROUGH {
			continue // not a fallthrough
		}

		nextCaseClause := switchStmt.Body.List[i+1].(*ast.CaseClause)
		if nextCaseClause.List == nil {
			// The next clause is 'default:', and this is a valid pattern.
			// Skip reporting this fallthrough.
			continue
		}

		if _, ok := w.commentsMap[branchStmt]; ok {
			// The fallthrough has a comment, still report with lower confidence.
			w.onFailure(lint.Failure{
				Confidence: 0.5,
				Node:       branchStmt,
				Category:   lint.FailureCategoryCodeStyle,
				Failure:    `this "fallthrough" can be removed by consolidating this case clause with the next one`,
			})
		} else {
			w.onFailure(lint.Failure{
				Confidence: 1.0,
				Node:       branchStmt,
				Category:   lint.FailureCategoryCodeStyle,
				Failure:    `this "fallthrough" can be removed by consolidating this case clause with the next one`,
			})
		}

		ast.Walk(w, caseClause)
	}

	return nil
}
