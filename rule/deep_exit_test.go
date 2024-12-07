package rule

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"
)

func TestLintDeepExit_isTestExample(t *testing.T) {
	tests := []struct {
		name       string
		funcDecl   string
		isTestFile bool
		want       bool
	}{
		{
			name:       "Package level example",
			funcDecl:   "func Example() {}",
			isTestFile: true,
			want:       true,
		},
		{
			name:       "Function example",
			funcDecl:   "func ExampleFunction() {}",
			isTestFile: true,
			want:       true,
		},
		{
			name:       "Method example",
			funcDecl:   "func ExampleType_Method() {}",
			isTestFile: true,
			want:       true,
		},
		{
			name:       "Wrong example function",
			funcDecl:   "func Examplemethod() {}",
			isTestFile: true,
			want:       false,
		},
		{
			name:       "Not an example",
			funcDecl:   "func NotAnExample() {}",
			isTestFile: true,
			want:       false,
		},
		{
			name:       "Example with parameters",
			funcDecl:   "func ExampleWithParams(a int) {}",
			isTestFile: true,
			want:       false,
		},
		{
			name:       "Not a test file",
			funcDecl:   "func Example() {}",
			isTestFile: false,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, err := parser.ParseFile(fs, "", "package main\n"+tt.funcDecl, parser.AllErrors)
			if err != nil {
				t.Fatal(err)
			}
			idx := slices.IndexFunc(node.Decls, func(decl ast.Decl) bool {
				_, ok := decl.(*ast.FuncDecl)
				return ok
			})
			fd := node.Decls[idx].(*ast.FuncDecl)

			w := &lintDeepExit{isTestFile: tt.isTestFile}
			got := w.isTestExample(fd)
			if got != tt.want {
				t.Errorf("isTestExample() = %v, want %v", got, tt.want)
			}
		})
	}
}
