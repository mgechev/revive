package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ReturnArgumentsLimitRule lints given else constructs.
type ReturnArgumentsLimitRule struct{}

// Apply applies the rule to given file.
func (r *ReturnArgumentsLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	if len(arguments) != 1 {
		panic(`invalid configuration for "return-argument-limit"`)
	}

	total, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		panic(`invalid value passed as return argument number to the "return-argument-limit" rule`)
	}

	var failures []lint.Failure

	walker := lintReturnArgumentsNum{
		total: int(total),
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *ReturnArgumentsLimitRule) Name() string {
	return "return-argument-limit"
}

type lintReturnArgumentsNum struct {
	total     int
	onFailure func(lint.Failure)
}

func (w lintReturnArgumentsNum) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.FuncDecl)
	if ok {
		num := 0
		if node.Type.Results != nil {
			num = node.Type.Results.NumFields()
		}
		if num > w.total {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("maximum number of return argument per function exceeded; max %d but got %d", w.total, num),
				Node:       node.Type,
			})
			return w
		}
	}
	return w
}
