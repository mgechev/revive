package defaultrule

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// ArgumentsLimitRule lints given else constructs.
type ArgumentsLimitRule struct{}

// Apply applies the rule to given file.
func (r *ArgumentsLimitRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	if len(arguments) != 1 {
		panic(`invalid configuration for "argument-limit"`)
	}
	total, err := strconv.ParseInt(arguments[0], 10, 32)
	if err != nil {
		panic(`invalid configuration for "argument-limit"`)
	}

	var failures []rule.Failure

	walker := lintArgsNum{
		total: total,
		onFailure: func(failure rule.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.GetAST())

	return failures
}

// Name returns the rule name.
func (r *ArgumentsLimitRule) Name() string {
	return "argument-limit"
}

type lintArgsNum struct {
	total     int64
	onFailure func(rule.Failure)
}

func (w lintArgsNum) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.FuncDecl)
	if ok {
		num := int64(len(node.Type.Params.List))
		if num > w.total {
			w.onFailure(rule.Failure{
				Failure: fmt.Sprintf("maximum number of arguments per function exceeded; max %d but got %d", w.total, num),
				Type:    rule.FailureTypeWarning,
				Node:    node.Type,
			})
			return nil
		}
	}
	return w
}
