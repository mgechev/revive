package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UnexportedNamingRule lints wrongly named unexported symbols.
type UnexportedNamingRule struct{}

// Apply applies the rule to given file.
func (r *UnexportedNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	r.checkMethods(onFailure, file.AST.Decls)

	ba := &unexportablenamingLinter{onFailure}
	ast.Walk(ba, file.AST)

	return failures
}

func (r *UnexportedNamingRule) checkMethods(onFailure func(failure lint.Failure), decls []ast.Decl) {
	for _, decl := range decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Recv == nil {
			continue // it's a function not a method
		}

		methodName := funcDecl.Name.Name
		if !ast.IsExported(methodName) {
			continue // it's an unexported method
		}

		recvTypes := astutils.GetTypeNames(funcDecl.Recv)
		recvType := strings.TrimLeft(recvTypes[0], "*")
		if ast.IsExported(recvType) {
			continue // it's an exported receiver
		}

		onFailure(lint.Failure{
			Confidence: 0.8,
			Failure:    fmt.Sprintf("method %q is exported but it's attached to the unexported type %q", methodName, recvType),
			Node:       funcDecl.Name,
		})
	}
}

// Name returns the rule name.
func (*UnexportedNamingRule) Name() string {
	return "unexported-naming"
}

type unexportablenamingLinter struct {
	onFailure func(lint.Failure)
}

func (unl unexportablenamingLinter) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		unl.lintFunction(n.Type, n.Body)
		return nil
	case *ast.FuncLit:
		unl.lintFunction(n.Type, n.Body)

		return nil
	case *ast.AssignStmt:
		if n.Tok != token.DEFINE {
			return nil
		}

		ids := []*ast.Ident{}
		for _, e := range n.Lhs {
			id, ok := e.(*ast.Ident)
			if !ok {
				continue
			}
			ids = append(ids, id)
		}

		unl.lintIDs(ids)

	case *ast.DeclStmt:
		gd, ok := n.Decl.(*ast.GenDecl)
		if !ok {
			return nil
		}

		if len(gd.Specs) < 1 {
			return nil
		}

		vs, ok := gd.Specs[0].(*ast.ValueSpec)
		if !ok {
			return nil
		}

		unl.lintIDs(vs.Names)
	}

	return unl
}

func (unl unexportablenamingLinter) lintFunction(ft *ast.FuncType, body *ast.BlockStmt) {
	unl.lintFields(ft.Params)
	unl.lintFields(ft.Results)

	if body != nil {
		ast.Walk(unl, body)
	}
}

func (unl unexportablenamingLinter) lintFields(fields *ast.FieldList) {
	if fields == nil {
		return
	}

	ids := []*ast.Ident{}
	for _, field := range fields.List {
		ids = append(ids, field.Names...)
	}

	unl.lintIDs(ids)
}

func (unl unexportablenamingLinter) lintIDs(ids []*ast.Ident) {
	for _, id := range ids {
		if id.IsExported() {
			unl.onFailure(lint.Failure{
				Node:       id,
				Confidence: 1,
				Category:   lint.FailureCategoryNaming,
				Failure:    fmt.Sprintf("the symbol %s is local, its name should start with a lowercase letter", id.String()),
			})
		}
	}
}
