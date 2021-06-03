package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

type NestedStructs struct{}

func (r *NestedStructs) Name() string {
	return "nested-structs"
}

func (r *NestedStructs) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if len(arguments) > 0 {
		panic(r.Name() + " doesn't take any arguments")
	}

	walker := &lintNestedStructs{
		fileAST: file.AST,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	if walker.count > 0 {
		walker.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("no nested structs are allowed, got %d", walker.count),
			Confidence: 1,
			Node:       file.AST,
			Category:   "style",
		})
	}

	return failures
}

type lintNestedStructs struct {
	count     int64
	fileAST   *ast.File
	onFailure func(lint.Failure)
}

func (w *lintNestedStructs) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.Field:
		if _, ok := v.Type.(*ast.StructType); ok {
			w.count++
			break
		}
	}
	return w
}
