package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// LintRangesRule lints given else constructs.
type LintRangesRule struct{}

// Apply applies the rule to given file.
func (r *LintRangesRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRanges{file, onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *LintRangesRule) Name() string {
	return "no-else-return"
}

type lintRanges struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintRanges) Visit(node ast.Node) ast.Visitor {
	rs, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}
	if rs.Value == nil {
		// for x = range m { ... }
		return w // single var form
	}
	if !isIdent(rs.Value, "_") {
		// for ?, y = range m { ... }
		return w
	}

	w.onFailure(lint.Failure{
		Failure:    fmt.Sprintf("should omit 2nd value from range; this loop is equivalent to `for %s %s range ...`", w.file.Render(rs.Key), rs.Tok),
		Confidence: 1,
		Node:       rs.Value,
	})

	// TODO: replacement?
	return w
}
