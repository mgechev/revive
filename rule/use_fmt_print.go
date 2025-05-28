package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// UseFmtPrintRule lints calls to print and println.
type UseFmtPrintRule struct{}

// Apply applies the rule to given file.
func (*UseFmtPrintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintUseFmtPrint{onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseFmtPrintRule) Name() string {
	return "use-fmt-print"
}

type lintUseFmtPrint struct {
	onFailure func(lint.Failure)
}

func (w lintUseFmtPrint) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w // nothing to do, the node is not a call
	}

	id, ok := (ce.Fun).(*ast.Ident)
	if !ok {
		return nil
	}

	name := id.Name
	switch name {
	default:
		return nil // nothing to do, the call is not println(...) nor print(...)
	case "println", "print":

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       node,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf(`avoid using built-in function %q, use fmt.F%s(os.Stderr, ...) instead`, name, name),
		})

		return w
	}
}
