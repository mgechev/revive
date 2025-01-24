package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// RedundantTestMainExitRule suggests removing Exit call in TestMain function for test files.
type RedundantTestMainExitRule struct{}

// Apply applies the rule to given file.
func (*RedundantTestMainExitRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantTestMainExit{onFailure: onFailure, isTestFile: file.IsTest()}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*RedundantTestMainExitRule) Name() string {
	return "redundant-test-main-exit"
}

type lintRedundantTestMainExit struct {
	onFailure  func(lint.Failure)
	isTestFile bool
}

func (w *lintRedundantTestMainExit) Visit(node ast.Node) ast.Visitor {
	if !w.isTestFile {
		return nil // skip analysis of this file if it is not a test file
	}

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

	pkg := id.Name
	fn := fc.Sel.Name
	if isCallToExitFunction(pkg, fn) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("redundant call to %s.%s in TestMain function, the test runner will handle it automatically as of Go 1.15", pkg, fn),
		})
	}

	return w
}

func (w *lintRedundantTestMainExit) isTestMain(fd *ast.FuncDecl) bool {
	return w.isTestFile && fd.Name.Name == "TestMain"
}
