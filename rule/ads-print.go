package rule

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ADSPrintRule lints program exit at functions other than main or init.
type ADSPrintRule struct{}

// Apply applies the rule to given file.
func (r *ADSPrintRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	if strings.HasSuffix(file.Name, "_test.go") {
		return failures
	}

	skipPackages := map[string]bool{"\"bdd\"": true}

	for _, ispec := range file.AST.Imports {
		iName := ispec.Path.Value
		if skipPackages[iName] {
			return failures
		}
	}

	var printFunctions = map[string]map[string]bool{
		"fmt": map[string]bool{
			"Print":   true,
			"Printf":  true,
			"Println": true,
		},
	}

	w := lintADSPrint{onFailure, printFunctions}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *ADSPrintRule) Name() string {
	return "ads-print"
}

type lintADSPrint struct {
	onFailure      func(lint.Failure)
	printFunctions map[string]map[string]bool
}

func (w lintADSPrint) Visit(node ast.Node) ast.Visitor {
	se, ok := node.(*ast.ExprStmt)
	if !ok {
		return w
	}
	ce, ok := se.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	fc, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}
	id, ok := fc.X.(*ast.Ident)
	if !ok {
		return w
	}

	fn := fc.Sel.Name
	pkg := id.Name
	if w.printFunctions[pkg] != nil && w.printFunctions[pkg][fn] { // it's a call to a print function
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   "bad practice",
			URL:        "#deep-exit",
			Failure:    fmt.Sprintf("do not call %s.%s, use logger", pkg, fn),
		})
	}

	return w
}
