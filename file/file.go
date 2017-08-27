package file

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// File abstraction used for representing files.
type File struct {
	Name    string
	files   *token.FileSet
	Content []byte
	ast     *ast.File
}

// New creates a new file
func New(name string, content []byte, files *token.FileSet) (*File, error) {
	f, err := parser.ParseFile(files, name, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &File{
		Name:    name,
		Content: content,
		files:   files,
		ast:     f,
	}, nil
}

// ToPosition returns line and column for given position.
func (f *File) ToPosition(pos token.Pos) token.Position {
	return f.files.Position(pos)
}

// GetAST returns the AST of the file
func (f *File) GetAST() *ast.File {
	return f.ast
}
