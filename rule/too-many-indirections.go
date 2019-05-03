package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// TooManyIndirectionsRule lints given else constructs.
type TooManyIndirectionsRule struct{}

const defaultLimit int = 3

// Apply applies the rule to given file.
func (r *TooManyIndirectionsRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	limit := defaultLimit
	if len(args) > 0 {
		var ok bool
		confLimit, ok := args[0].(int64) // Alt. non panicking version
		if !ok {
			panic(fmt.Sprintf("Invalid argument for too-many-indirections expected an int64, got %T", args[0]))
		}
		limit = int(confLimit)
	}

	walker := &lintTooManyIndirections{
		limit: limit,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *TooManyIndirectionsRule) Name() string {
	return "too-many-indirections"
}

type lintTooManyIndirections struct {
	limit     int
	onFailure func(lint.Failure)
}

// Visit looks for selector expressions.
// If the expression has too many indirections (.) a failure will be created.
func (w *lintTooManyIndirections) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		c := w.countIndirections(n)
		if c > w.limit {
			w.onFailure(lint.Failure{
				Category:   "bad practice",
				Confidence: 1,
				Node:       n,
				Failure:    fmt.Sprintf("Too many %d (>%d) indirections in expression", c, w.limit),
			})
		}

		return nil // skip analysis of this subtree
	}

	return w
}

// countIndirections returns the length of an indirection chain
func (w lintTooManyIndirections) countIndirections(node ast.Node) int {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		return 1 + w.countIndirections(n.X)
	case *ast.IndexExpr:
		return w.countIndirections(n.X)
	case *ast.CallExpr:
		return w.countIndirections(n.Fun)
	}

	return 0
}
