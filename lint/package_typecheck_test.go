package lint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestPackageTypeCheckWithEmptyGOROOT(t *testing.T) {
	t.Setenv("GOROOT", "")

	fset := token.NewFileSet()
	src := `package p
import "fmt"
func f() { fmt.Println("ok") }
`

	parsed, err := parser.ParseFile(fset, "p.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	pkg := &Package{
		fset: fset,
		files: map[string]*File{
			"p.go": {
				Name: "p.go",
				AST:  parsed,
				Pkg:  nil,
			},
		},
	}
	pkg.files["p.go"].Pkg = pkg

	if err := pkg.TypeCheck(); err != nil {
		t.Fatalf("type check failed with empty GOROOT: %v", err)
	}

	var printCall *ast.CallExpr
	ast.Inspect(parsed, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if ok {
			printCall = call
			return false
		}
		return true
	})
	if printCall == nil {
		t.Fatal("expected to find fmt.Println call")
	}
	if got := pkg.TypeOf(printCall.Fun); got == nil {
		t.Fatal("expected type info for fmt.Println call")
	}
}
