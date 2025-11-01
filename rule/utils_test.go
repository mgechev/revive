package rule

import (
	"fmt"
	"go/ast"
	"testing"
)

func TestIsCallToExitFunction(t *testing.T) {
	tests := []struct {
		pkgName      string
		functionName string
		functionArgs []ast.Expr
		expected     bool
	}{
		{"os", "Exit", nil, true},
		{"syscall", "Exit", nil, true},
		{"log", "Fatal", nil, true},
		{"log", "Fatalf", nil, true},
		{"log", "Fatalln", nil, true},
		{"log", "Panic", nil, true},
		{"log", "Panicf", nil, true},
		{"flag", "Parse", nil, true},
		{"flag", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "flag"},
				Sel: &ast.Ident{Name: "ExitOnError"},
			},
		}, true},
		{"log", "Print", nil, false},
		{"fmt", "Errorf", nil, false},
		{"flag", "NewFlagSet", []ast.Expr{
			nil,
			&ast.SelectorExpr{
				X:   &ast.Ident{Name: "flag"},
				Sel: &ast.Ident{Name: "ContinueOnError"},
			},
		}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.pkgName, tt.functionName), func(t *testing.T) {
			if got := isCallToExitFunction(tt.pkgName, tt.functionName, tt.functionArgs); got != tt.expected {
				t.Errorf("isCallToExitFunction(%s, %s, %v) = %v; want %v", tt.pkgName, tt.functionName, tt.functionArgs, got, tt.expected)
			}
		})
	}
}
