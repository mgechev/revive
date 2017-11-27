package defaultrule

import (
	"go/ast"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// ImportsRule lints given else constructs.
type ImportsRule struct{}

// Apply applies the rule to given file.
func (r *ImportsRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	var failures []rule.Failure

	fileAst := file.GetAST()
	walker := lintImports{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure rule.Failure) {
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
	file      *file.File
	fileAst   *ast.File
	onFailure func(rule.Failure)
}

func (w lintImports) Visit(n ast.Node) ast.Visitor {
	for _, is := range w.fileAst.Imports {
		if is.Name != nil && is.Name.Name == "." && !isTest(w.file) {
			w.onFailure(rule.Failure{
				Confidence: 1,
				Failure:    "should not use dot imports",
				Node:       is,
			})
		}

	}
	return nil
}
