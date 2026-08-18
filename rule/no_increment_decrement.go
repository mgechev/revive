package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// NoIncrementDecrementRule suggests replacing `i++` and `i--` with `i += 1` and `i -= 1`,
// except when they are used as the post statement of a for loop (i.e. as loop counters).
// It is the opposite of the increment-decrement rule.
type NoIncrementDecrementRule struct{}

// Apply applies the rule to given file.
func (*NoIncrementDecrementRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := &lintNoIncrementDecrement{
		file:      file,
		loopPosts: map[ast.Stmt]struct{}{},
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoIncrementDecrementRule) Name() string {
	return "no-increment-decrement"
}

type lintNoIncrementDecrement struct {
	file      *lint.File
	loopPosts map[ast.Stmt]struct{}
	onFailure func(lint.Failure)
}

func (w *lintNoIncrementDecrement) Visit(n ast.Node) ast.Visitor {
	switch stmt := n.(type) {
	case *ast.ForStmt:
		if stmt.Post != nil {
			// The post statement of a for loop is the idiomatic place for a loop
			// counter, so it is exempted from the rule.
			w.loopPosts[stmt.Post] = struct{}{}
		}
	case *ast.IncDecStmt:
		if _, isLoopPost := w.loopPosts[stmt]; isLoopPost {
			return w
		}

		var replacement string
		switch stmt.Tok {
		case token.INC:
			replacement = "+= 1"
		case token.DEC:
			replacement = "-= 1"
		default:
			return w
		}

		w.onFailure(lint.Failure{
			Confidence: 0.8,
			Node:       stmt,
			Category:   lint.FailureCategoryUnaryOp,
			Failure:    fmt.Sprintf("should replace %s with %s %s", w.file.Render(stmt), w.file.Render(stmt.X), replacement),
		})
	}

	return w
}
