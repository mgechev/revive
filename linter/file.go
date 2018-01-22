package linter

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// File abstraction used for representing files.
type File struct {
	Name    string
	pkg     *Package
	Content []byte
	ast     *ast.File
}

// NewFile creates a new file
func NewFile(name string, content []byte, pkg *Package) (*File, error) {
	f, err := parser.ParseFile(pkg.Fset, name, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &File{
		Name:    name,
		Content: content,
		pkg:     pkg,
		ast:     f,
	}, nil
}

// ToPosition returns line and column for given position.
func (f *File) ToPosition(pos token.Pos) token.Position {
	return f.pkg.Fset.Position(pos)
}

// GetAST returns the AST of the file
func (f *File) GetAST() *ast.File {
	return f.ast
}

func (f *File) isMain() bool {
	if f.GetAST().Name.Name == "main" {
		return true
	}
	return false
}
