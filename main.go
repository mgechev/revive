package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/mgechev/golinter/syntaxvisitor"
)

type CustomLinter struct {
	syntaxvisitor.SyntaxVisitor
}

func (w *CustomLinter) VisitIdent(node *ast.Ident) {
	fmt.Println("Child", node.Name)
}

// This example demonstrates how to inspect the AST of a Go program.
func ExampleInspect() {
	// src is the input for which we want to inspect the AST.
	src := `
  package p
  const c = 1.0
  var X = f(3.14)*2 + c
  `

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	var visitor CustomLinter
	visitor.SyntaxVisitor.Impl = &visitor
	visitor.Visit(f)

	// output:
	// src.go:2:9:	p
	// src.go:3:7:	c
	// src.go:3:11:	1.0
	// src.go:4:5:	X
	// src.go:4:9:	f
	// src.go:4:11:	3.14
	// src.go:4:17:	2
	// src.go:4:21:	c
}

func main() {
	ExampleInspect()
}
