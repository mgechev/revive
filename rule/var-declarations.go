package rule

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"github.com/mgechev/revive/linter"
)

// VarDeclarationsRule lints given else constructs.
type VarDeclarationsRule struct{}

// Apply applies the rule to given file.
func (r *VarDeclarationsRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
	var failures []linter.Failure

	fileAst := file.GetAST()
	walker := &lintVarDeclarations{
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
func (r *VarDeclarationsRule) Name() string {
	return "blank-imports"
}

type lintVarDeclarations struct {
	fileAst   *ast.File
	file      *linter.File
	lastGen   *ast.GenDecl
	onFailure func(linter.Failure)
}

func (w *lintVarDeclarations) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.GenDecl:
		if v.Tok != token.CONST && v.Tok != token.VAR {
			return nil
		}
		w.lastGen = v
		return w
	case *ast.ValueSpec:
		if w.lastGen.Tok == token.CONST {
			return nil
		}
		if len(v.Names) > 1 || v.Type == nil || len(v.Values) == 0 {
			return nil
		}
		rhs := v.Values[0]
		// An underscore var appears in a common idiom for compile-time interface satisfaction,
		// as in "var _ Interface = (*Concrete)(nil)".
		if isIdent(v.Names[0], "_") {
			return nil
		}
		// If the RHS is a zero value, suggest dropping it.
		zero := false
		if lit, ok := rhs.(*ast.BasicLit); ok {
			zero = zeroLiteral[lit.Value]
		} else if isIdent(rhs, "nil") {
			zero = true
		}
		if zero {
			w.onFailure(linter.Failure{
				Confidence: 0.9,
				Node:       rhs,
				Failure:    fmt.Sprintf("should drop = %s from declaration of var %s; it is the zero value", render(rhs), v.Names[0]),
			})
			return nil
		}
	}
	return w
}

func render(x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, token.NewFileSet(), x); err != nil {
		panic(err)
	}
	return buf.String()
}
