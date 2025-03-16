package rule

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// TimeAfterLeak spots potential goroutine leaks caused by Time.After()
type TimeAfterLeak struct{}

// Apply applies the rule to given file.
func (r *TimeAfterLeak) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsAtLeastGoVersion(lint.Go123) {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeAfterLeak{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*TimeAfterLeak) Name() string {
	return "time-after-leak"
}

type lintTimeAfterLeak struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeAfterLeak) Visit(node ast.Node) ast.Visitor {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return w
	}

	cases := selectStmt.Body.List
	if len(cases) <= 1 {
		return w // even if the single case is a read from time.After there is nothing to warn about
	}

	for _, c := range cases {
		commClause, ok := c.(*ast.CommClause)
		if !ok { // we check even if it is not possible to have case that is not a CommClause
			continue
		}

		comm := commClause.Comm
		if comm == nil {
			continue // is the default select case
		}

		// case something

		exprStmt, ok := comm.(*ast.ExprStmt)
		if !ok {
			continue // is not an expression statement (it's, for example, an assignment)
		}

		expr, ok := exprStmt.X.(*ast.UnaryExpr)
		isChannelRead := ok && expr.Op == token.ARROW
		if !isChannelRead {
			continue // is not a channel read expression
		}

		// case <- expr

		call, ok := expr.X.(*ast.CallExpr)
		if !ok {
			continue // is not a read from a channel returned by a function call
		}

		// case <- f()

		if isPkgDot(call.Fun, "time", "After") {
			// case <- time.After(...)
			w.onFailure(lint.Failure{
				Confidence: 0.8,
				Node:       call.Fun,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "the underlying goroutine of time.After() is not garbage-collected until timer expiration, prefer NewTimer+Timer.Stop",
			})
		}
	}

	return w
}
