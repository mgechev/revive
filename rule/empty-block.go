package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// EmptyBlockRule lints given else constructs.
type EmptyBlockRule struct{}

// Apply applies the rule to given file.
func (r *EmptyBlockRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintEmptyBlock{make(map[*ast.IfStmt]bool), onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *EmptyBlockRule) Name() string {
	return "empty-block"
}

type lintEmptyBlock struct {
	ignore    map[*ast.IfStmt]bool
	onFailure func(lint.Failure)
}

func (w lintEmptyBlock) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	if len(block.List) == 0 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       block,
			Category:   "logic",
			URL:        "#empty-block",
			Failure:    "this block is empty, you can remove it",
		})
	}

	return w
}
