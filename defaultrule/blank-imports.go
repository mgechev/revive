package defaultrule

import (
	"go/ast"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// BlankImportsRule lints given else constructs.
type BlankImportsRule struct{}

// Apply applies the rule to given file.
func (r *BlankImportsRule) Apply(file *file.File, arguments rule.Arguments) []rule.Failure {
	var failures []rule.Failure

	fileAst := file.GetAST()
	walker := lintBlankImports{
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
func (r *BlankImportsRule) Name() string {
	return "blank-imports"
}

type lintBlankImports struct {
	fileAst   *ast.File
	file      *file.File
	onFailure func(rule.Failure)
}

func (w lintBlankImports) Visit(n ast.Node) ast.Visitor {
	// In package main and in tests, we don't complain about blank imports.
	if w.fileAst.Name.Name == "main" || isTest(w.file) {
		return nil
	}

	// The first element of each contiguous group of blank imports should have
	// an explanatory comment of some kind.
	for i, imp := range w.fileAst.Imports {
		pos := w.file.ToPosition(imp.Pos())

		if !isBlank(imp.Name) {
			continue // Ignore non-blank imports.
		}
		if i > 0 {
			prev := w.fileAst.Imports[i-1]
			prevPos := w.file.ToPosition(prev.Pos())
			if isBlank(prev.Name) && prevPos.Line+1 == pos.Line {
				continue // A subsequent blank in a group.
			}
		}

		// This is the first blank import of a group.
		if imp.Doc == nil && imp.Comment == nil {
			w.onFailure(rule.Failure{
				Node:       imp,
				Failure:    "a blank import should be only in a main or test package, or have a comment justifying it",
				Confidence: 1,
			})
		}
	}
	return nil
}
