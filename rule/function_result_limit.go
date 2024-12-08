package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"sync"

	"github.com/mgechev/revive/lint"
)

// FunctionResultsLimitRule limits the maximum number of results a function can return.
type FunctionResultsLimitRule struct {
	max int

	configureOnce sync.Once
}

// Apply applies the rule to given file.
func (r *FunctionResultsLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var configureErr error
	r.configureOnce.Do(func() { configureErr = r.configure(arguments) })

	if configureErr != nil {
		return []lint.Failure{lint.NewInternalFailure(configureErr.Error())}
	}

	var failures []lint.Failure
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		num := 0
		hasResults := funcDecl.Type.Results != nil
		if hasResults {
			num = funcDecl.Type.Results.NumFields()
		}

		if num <= r.max {
			continue
		}

		failures = append(failures, lint.Failure{
			Confidence: 1,
			Failure:    fmt.Sprintf("maximum number of return results per function exceeded; max %d but got %d", r.max, num),
			Node:       funcDecl.Type,
		})
	}

	return failures
}

// Name returns the rule name.
func (*FunctionResultsLimitRule) Name() string {
	return "function-result-limit"
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
