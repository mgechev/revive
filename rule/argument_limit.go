package rule

import (
	"fmt"
	"go/ast"
	"sync"

	"github.com/mgechev/revive/lint"
)

// ArgumentsLimitRule lints the number of arguments a function can receive.
type ArgumentsLimitRule struct {
	max int

	configureOnce sync.Once
}

const defaultArgumentsLimit = 8

func (r *ArgumentsLimitRule) configure(arguments lint.Arguments) {
	if len(arguments) < 1 {
		r.max = defaultArgumentsLimit
		return
	}

	maxArguments, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		panic(`invalid value passed as argument number to the "argument-limit" rule`)
	}
	r.max = int(maxArguments)
}

// Apply applies the rule to given file.
func (r *ArgumentsLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configureOnce.Do(func() { r.configure(arguments) })

	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		numParams := 0
		for _, l := range funcDecl.Type.Params.List {
			numParams += len(l.Names)
		}

		if numParams > r.max {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("maximum number of arguments per function exceeded; max %d but got %d", r.max, numParams),
				Node:       funcDecl.Type,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*ArgumentsLimitRule) Name() string {
	return "argument-limit"
}
