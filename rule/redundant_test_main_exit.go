package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// RedundantTestMainExitRule suggests removing redundant [os.Exit] or [syscall.Exit] calls in TestMain function.
type RedundantTestMainExitRule struct{}

// Apply applies the rule to given file.
func (*RedundantTestMainExitRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if !file.IsTest() || !file.Pkg.IsAtLeastGoVersion(lint.Go115) {
		// skip analysis for non-test files or for Go versions before 1.15
		return failures
	}

	for _, decl := range file.AST.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Body == nil || !astutils.FuncSignatureIs(fd, "TestMain", []string{"*testing.M"}, nil) || len(fd.Type.Params.List[0].Names) != 1 {
			continue
		}

		w := &lintRedundantTestMainExit{
			mName:  fd.Type.Params.List[0].Names[0].Name,
			writes: map[string][]ast.Expr{},
			onFailure: func(failure lint.Failure) {
				failures = append(failures, failure)
			},
		}
		ast.Inspect(fd.Body, w.collectWrites)
		ast.Inspect(fd.Body, w.checkExitCalls)
	}

	return failures
}

// Name returns the rule name.
func (*RedundantTestMainExitRule) Name() string {
	return "redundant-test-main-exit"
}

type lintRedundantTestMainExit struct {
	mName     string                // name of the *testing.M parameter
	writes    map[string][]ast.Expr // values written to each variable; nil for writes with an unknown value (e.g. code++ or &code)
	onFailure func(lint.Failure)
}

// collectWrites records every write to a variable: only a variable written exactly once, from m.Run, holds its result.
func (w *lintRedundantTestMainExit) collectWrites(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		for i, lhs := range n.Lhs {
			var value ast.Expr
			isPlainAssign := n.Tok == token.ASSIGN || n.Tok == token.DEFINE
			if isPlainAssign && len(n.Lhs) == len(n.Rhs) {
				value = n.Rhs[i]
			}
			w.recordWrite(lhs, value)
		}
	case *ast.ValueSpec:
		for i, name := range n.Names {
			if len(n.Values) == 0 {
				continue // declaration without a value is not a write
			}
			var value ast.Expr
			if len(n.Names) == len(n.Values) {
				value = n.Values[i]
			}
			w.recordWrite(name, value)
		}
	case *ast.IncDecStmt:
		w.recordWrite(n.X, nil)
	case *ast.UnaryExpr:
		if n.Op == token.AND {
			w.recordWrite(n.X, nil)
		}
	case *ast.RangeStmt:
		w.recordWrite(n.Key, nil)
		w.recordWrite(n.Value, nil)
	}

	return true
}

func (w *lintRedundantTestMainExit) recordWrite(target, value ast.Expr) {
	id, ok := ast.Unparen(target).(*ast.Ident)
	if !ok || id.Name == "_" {
		return
	}

	w.writes[id.Name] = append(w.writes[id.Name], value)
}

// checkExitCalls reports [os.Exit] and [syscall.Exit] calls whose argument is the result of m.Run.
func (w *lintRedundantTestMainExit) checkExitCalls(node ast.Node) bool {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}

	var pkg string
	switch {
	case astutils.IsPkgDotName(ce.Fun, "os", "Exit"):
		pkg = "os"
	case astutils.IsPkgDotName(ce.Fun, "syscall", "Exit"):
		pkg = "syscall"
	default:
		return true
	}

	if len(ce.Args) != 1 || !w.isRunResult(ce.Args[0]) {
		return true
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("redundant call to %s.Exit in TestMain function, the test runner will handle it automatically as of Go 1.15", pkg),
	})

	return true
}

// isRunResult returns true if expr is a call to m.Run or a variable whose only write is the result of m.Run.
func (w *lintRedundantTestMainExit) isRunResult(expr ast.Expr) bool {
	expr = ast.Unparen(expr)
	id, ok := expr.(*ast.Ident)
	if !ok {
		return w.isRunCall(expr)
	}

	writes := w.writes[id.Name]
	return len(writes) == 1 && writes[0] != nil && w.isRunCall(writes[0])
}

func (w *lintRedundantTestMainExit) isRunCall(expr ast.Expr) bool {
	ce, ok := ast.Unparen(expr).(*ast.CallExpr)
	if !ok {
		return false
	}

	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok || se.Sel.Name != "Run" {
		return false
	}

	return astutils.IsIdent(ast.Unparen(se.X), w.mName)
}
