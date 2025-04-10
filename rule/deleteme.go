package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// DeletemeRule is a sandbox rule to tests ideas
type DeletemeRule struct{}

// Apply applies the rule to given file.
func (r *DeletemeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeletemeRule{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*DeletemeRule) Name() string {
	return "deleteme"
}

type lintDeletemeRule struct {
	onFailure func(lint.Failure)
}

func (w *lintDeletemeRule) Visit(node ast.Node) ast.Visitor {
	// we visit the AST looking for "case <- time.After(...)"
	switch n := node.(type) {
	case *ast.IfStmt: // is a select case
		if n.Else == nil {
			return w
		}
		elseBlock, ok := n.Else.(*ast.BlockStmt)
		if !ok {
			return w
		}

		if len(n.Body.List) < 1 || len(n.Body.List) > 1 || len(elseBlock.List) > 1 {
			return w
		}

		if !isReturnBoolean(n.Body.List[0]) {
			return w
		}

		if !isReturnBoolean(elseBlock.List[0]) {
			return w
		}

		w.onFailure(lint.Failure{
			Confidence: 0.8,
			Node:       n,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "return condition",
		})
	}
	return w
}

func isReturnBoolean(stmt ast.Stmt) bool {
	returnStmt, ok := stmt.(*ast.ReturnStmt)
	if !ok {
		return false
	}

	if len(returnStmt.Results) > 1 || len(returnStmt.Results) < 1 {
		return false
	}

	result := gofmt(returnStmt.Results[0])

	return result == "true" || result == "false"
}
