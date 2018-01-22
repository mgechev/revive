package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/linter"
)

// LintRangesRule lints given else constructs.
type LintRangesRule struct{}

// Apply applies the rule to given file.
func (r *LintRangesRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
	var failures []linter.Failure

	onFailure := func(failure linter.Failure) {
		failures = append(failures, failure)
	}

	astFile := file.GetAST()
	w := &lintRanges{astFile, onFailure}
	ast.Walk(w, astFile)
	return failures
}

// Name returns the rule name.
func (r *LintRangesRule) Name() string {
	return "no-else-return"
}

type lintRanges struct {
	file      *ast.File
	onFailure func(linter.Failure)
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

	w.onFailure(linter.Failure{
		Failure:    fmt.Sprintf("should omit 2nd value from range; this loop is equivalent to `for %s %s range ...`", render(rs.Key), rs.Tok),
		Confidence: 1,
		Node:       rs.Value,
	})

	// TODO: replacement?
	return w
}
