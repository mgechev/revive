package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"sync"

	"github.com/mgechev/revive/lint"
)

// FunctionResultsLimitRule lints given else constructs.
type FunctionResultsLimitRule struct {
	max int
}

const defaultResultsLimit = 3

func (r *FunctionResultsLimitRule) configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.max = defaultResultsLimit
		return nil
	}

	maxResults, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		return fmt.Errorf(`invalid value passed as return results number to the "function-result-limit" rule; need int64 but got %T`, arguments[0])
	}
	if maxResults < 0 {
		return errors.New(`the value passed as return results number to the "function-result-limit" rule cannot be negative`)
	}

	r.max = int(maxResults)
	return nil
}

// Apply applies the rule to given file.
func (r *FunctionResultsLimitRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	check := sync.OnceValue(func() error { return r.configure(arguments) })
	if err := check(); err != nil {
		return nil, err
	}

	var failures []lint.Failure

	walker := lintFunctionResultsNum{
		max: r.max,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures, nil
}

// Name returns the rule name.
func (*FunctionResultsLimitRule) Name() string {
	return "function-result-limit"
}

type lintFunctionResultsNum struct {
	max       int
	onFailure func(lint.Failure)
}

func (w lintFunctionResultsNum) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.FuncDecl)
	if ok {
		num := 0
		hasResults := node.Type.Results != nil
		if hasResults {
			num = node.Type.Results.NumFields()
		}
		if num > w.max {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("maximum number of return results per function exceeded; max %d but got %d", w.max, num),
				Node:       node.Type,
			})
		}

		return nil // skip visiting function's body
	}

	return w
}
