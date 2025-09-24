package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// RedundantAssignmentRule detects unnecessary "self-assignments" of range variables.
type RedundantAssignmentRule struct{}

// Apply applies the rule to given file.
func (*RedundantAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go122) {
		return nil
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantAssignmentRule{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*RedundantAssignmentRule) Name() string {
	return "redundant-assignment"
}

type lintRedundantAssignmentRule struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantAssignmentRule) Visit(node ast.Node) ast.Visitor {
	// Only interested in  range statements
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w // not a range statement
	}

	body := rangeStmt.Body
	if body == nil || len(body.List) == 0 {
		return nil // empty body
	}

	stmt := body.List[0]

	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok || assignStmt.Tok != token.DEFINE || len(assignStmt.Lhs) != 1 {
		return w // not an assignment statement of the form x := something
	}

	lhsStr, _ := astutils.GetIdentName(assignStmt.Lhs[0])
	rhsStr, _ := astutils.GetIdentName(assignStmt.Rhs[0])
	if lhsStr != rhsStr {
		return w // not an assignment of the form x := x
	}

	keyStr, _ := astutils.GetIdentName(rangeStmt.Key)
	valueStr, _ := astutils.GetIdentName(rangeStmt.Value)

	if lhsStr != keyStr && lhsStr != valueStr {
		return w // the assigned variable is not from the range statement
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       assignStmt,
		Category:   lint.FailureCategoryOptimization,
		Failure:    fmt.Sprintf("redundant assignment of range variable %q, use it directly", lhsStr),
	})

	return w
}
