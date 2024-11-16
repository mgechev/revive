// Package rule implements revive's linting rules.
package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

// MaxPublicStructsRule lints given else constructs.
type MaxPublicStructsRule struct {
	max int64
	sync.Mutex
}

const defaultMaxPublicStructs = 5

func (r *MaxPublicStructsRule) configure(arguments lint.Arguments) error {
	r.Lock()
	defer r.Unlock()
	if r.max == 0 {
		return nil // already configured
	}

	if len(arguments) < 1 {
		r.max = defaultMaxPublicStructs
		return nil
	}

	check := checkNumberOfArguments(1, arguments, r.Name())
	if check != nil {
		return check
	}

	maxStructs, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		return errors.New(`invalid value passed as argument number to the "max-public-structs" rule`)
	}
	r.max = maxStructs
	return nil
}

// Apply applies the rule to given file.
func (r *MaxPublicStructsRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	var failures []lint.Failure
	err := r.configure(arguments)
	if err != nil {
		return failures, err
	}

	fileAst := file.AST

	walker := &lintMaxPublicStructs{
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	if walker.current > r.max {
		walker.onFailure(lint.Failure{
			Failure:    "you have exceeded the maximum number of public struct declarations",
			Confidence: 1,
			Node:       fileAst,
			Category:   "style",
		})
	}

	return failures, nil
}

// Name returns the rule name.
func (*MaxPublicStructsRule) Name() string {
	return "max-public-structs"
}

type lintMaxPublicStructs struct {
	current   int64
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

func (w *lintMaxPublicStructs) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.TypeSpec:
		name := v.Name.Name
		first := string(name[0])
		if strings.ToUpper(first) == first {
			w.current++
		}
	}
	return w
}
