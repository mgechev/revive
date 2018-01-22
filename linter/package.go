package linter

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/gcexportdata"
)

// Package represents a package in the project.
type Package struct {
	fset  *token.FileSet
	files map[string]*File

	TypesPkg  *types.Package
	TypesInfo *types.Info

	// sortable is the set of types in the package that implement sort.Interface.
	Sortable map[string]bool
	// main is whether this is a "main" package.
	main int
}

var newImporter = func(fset *token.FileSet) types.ImporterFrom {
	return gcexportdata.NewImporter(fset, make(map[string]*types.Package))
}

var (
	trueValue  = 1
	falseValue = 2
	notSet     = 3
)

// IsMain returns if that's the main package.
func (p *Package) IsMain() bool {
	if p.main == trueValue {
		return true
	} else if p.main == falseValue {
		return false
	}
	for _, f := range p.files {
		if f.isMain() {
			p.main = trueValue
			return true
		}
	}
	p.main = falseValue
	return false
}

// TypeCheck performs type checking for given package.
func (p *Package) TypeCheck() error {
	config := &types.Config{
		// By setting a no-op error reporter, the type checker does as much work as possible.
		Error:    func(error) {},
		Importer: newImporter(p.fset),
	}
	info := &types.Info{
		Types:  make(map[ast.Expr]types.TypeAndValue),
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
		Scopes: make(map[ast.Node]*types.Scope),
	}
	var anyFile *File
	var astFiles []*ast.File
	for _, f := range p.files {
		anyFile = f
		astFiles = append(astFiles, f.AST)
	}
	typesPkg, err := config.Check(anyFile.AST.Name.Name, p.fset, astFiles, info)
	// Remember the typechecking info, even if config.Check failed,
	// since we will get partial information.
	p.TypesPkg = typesPkg
	p.TypesInfo = info
	return err
}

// TypeOf returns the type of an expression.
func (p *Package) TypeOf(expr ast.Expr) types.Type {
	if p.TypesInfo == nil {
		return nil
	}
	return p.TypesInfo.TypeOf(expr)
}

func (p *Package) lint(rules []Rule, config RulesConfig) []Failure {
	var failures []Failure
	p.TypeCheck()
	for _, file := range p.files {
		failures = append(failures, file.lint(rules, config)...)
	}
	return failures
}
