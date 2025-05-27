package rule

import (
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseFmtPrint lints calls to print and println.
type UseFmtPrint struct{}

// Apply applies the rule to given file.
func (*UseFmtPrint) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintUseFmtPrint{onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseFmtPrint) Name() string {
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

	if !astutils.IsIdent(ce.Fun, "println") && !astutils.IsIdent(ce.Fun, "print") {
		return nil // nothing to do, the call is not println(...) nor print(...)
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "avoid using built in printing functions, use fmt.Print or fmt.Println",
	})

	return w
}
