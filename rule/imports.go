package rule

import (
	"go/ast"

	"github.com/mgechev/revive/linter"
)

// ImportsRule lints given else constructs.
type ImportsRule struct{}

// Apply applies the rule to given file.
func (r *ImportsRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
	var failures []linter.Failure

	fileAst := file.AST
	walker := lintImports{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure linter.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ImportsRule) Name() string {
	return "imports"
}

type lintImports struct {
	file      *linter.File
	fileAst   *ast.File
	onFailure func(linter.Failure)
}

func (w lintImports) Visit(n ast.Node) ast.Visitor {
	for _, is := range w.fileAst.Imports {
		if is.Name != nil && is.Name.Name == "." && !isTest(w.file) {
			w.onFailure(linter.Failure{
				Confidence: 1,
				Failure:    "should not use dot imports",
				Node:       is,
			})
		}

	}
	return nil
}
