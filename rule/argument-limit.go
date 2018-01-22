package rule

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/mgechev/revive/linter"
)

// ArgumentsLimitRule lints given else constructs.
type ArgumentsLimitRule struct{}

// Apply applies the rule to given file.
func (r *ArgumentsLimitRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
	if len(arguments) != 1 {
		panic(`invalid configuration for "argument-limit"`)
	}
	total, err := strconv.ParseInt(arguments[0], 10, 32)
	if err != nil {
		panic(`invalid configuration for "argument-limit"`)
	}

	var failures []linter.Failure

	walker := lintArgsNum{
		total: total,
		onFailure: func(failure linter.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *ArgumentsLimitRule) Name() string {
	return "argument-limit"
}

type lintArgsNum struct {
	total     int64
	onFailure func(linter.Failure)
}

func (w lintArgsNum) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.FuncDecl)
	if ok {
		num := int64(len(node.Type.Params.List))
		if num > w.total {
			w.onFailure(linter.Failure{
				Failure: fmt.Sprintf("maximum number of arguments per function exceeded; max %d but got %d", w.total, num),
				Type:    linter.FailureTypeWarning,
				Node:    node.Type,
			})
			return w
		}
	}
	return w
}
