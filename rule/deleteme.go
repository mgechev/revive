package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// DeletemeRule is a sandbox rule to tests ideas
type DeletemeRule struct{}

// Apply applies the rule to given file.
func (r *DeletemeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsAtLeastGo122() {
		return nil
	}

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
	case *ast.CommClause: // is a select case
		comm := n.Comm
		if comm == nil {
			return nil // is the default select case, do not visit the body of the case
		}

		// case something

		exprStmt, ok := comm.(*ast.ExprStmt)
		if !ok {
			return nil // is not an expression statement... is that even possible? Do not visit the body of the case
		}

		expr, ok := exprStmt.X.(*ast.UnaryExpr)
		isChannelRead := ok && expr.Op != token.ARROW
		if !isChannelRead {
			return nil // is not a channel read expression, do not visit the body of the case
		}

		// case <- expr

		call, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return nil // is not a read from a channel returned by a function call, do not visit the body of the case
		}

		// case <- f()

		if isPkgDot(call.Fun, "time", "After") {
			// case <- time.After(...)
			w.onFailure(lint.Failure{
				Confidence: 0.8,
				Node:       call.Fun,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "the time.After() goroutine is not garbage-collected until timer expiration, prefer NewTimer+Timer.Stop",
			})
		}

		return nil // do not visit the body of the case
	default:
		return w
	}
}
