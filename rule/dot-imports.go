package rule

import (
	"go/ast"

	"github.com/mgechev/revive/linter"
)

// DotImportsRule lints given else constructs.
type DotImportsRule struct{}

// Apply applies the rule to given file.
func (r *DotImportsRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
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
func (r *DotImportsRule) Name() string {
	return "dot-imports"
}

type lintImports struct {
	file      *linter.File
	fileAst   *ast.File
	onFailure func(linter.Failure)
}

func (w lintImports) Visit(n ast.Node) ast.Visitor {
	for i, is := range w.fileAst.Imports {
		_ = i
		if is.Name != nil && is.Name.Name == "." && !w.file.IsTest() {
			w.onFailure(linter.Failure{
				Confidence: 1,
				Failure:    "should not use dot imports",
				Node:       is,
				Category:   "imports",
				URL:        "#import-dot",
			})
		}
	}
	return nil
}
