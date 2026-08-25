package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ExitAfterDeferRule spots calls to [os.Exit], [syscall.Exit] or [log.Fatal] that
// are reachable after a defer statement, since such calls stop the program without
// running the deferred functions.
type ExitAfterDeferRule struct{}

// Name returns the rule name.
func (*ExitAfterDeferRule) Name() string {
	return "exit-after-defer"
}

// Apply applies the rule to given file.
func (*ExitAfterDeferRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	ast.Inspect(file.AST, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			if fn.Body != nil {
				checkExitAfterDefer(fn.Body, onFailure)
			}
			return false
		case *ast.FuncLit:
			checkExitAfterDefer(fn.Body, onFailure)
			return false
		}
		return true
	})

	return failures
}

// checkExitAfterDefer reports exit calls that follow a defer statement within
// the same function scope. Function literals are analyzed on their own because
// they introduce a new defer scope.
func checkExitAfterDefer(body *ast.BlockStmt, onFailure func(lint.Failure)) {
	deferred := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			checkExitAfterDefer(n.Body, onFailure)
			return false
		case *ast.DeferStmt:
			deferred = true
			// A deferred function literal has its own defer scope.
			if lit, ok := n.Call.Fun.(*ast.FuncLit); ok {
				checkExitAfterDefer(lit.Body, onFailure)
			}
			return false
		case *ast.CallExpr:
			if !deferred {
				return true
			}
			if pkg, fn, ok := isUnrecoverableExitCall(n); ok {
				onFailure(lint.Failure{
					Confidence: 1,
					Node:       n,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    fmt.Sprintf("%s.%s after a defer statement prevents deferred calls from running", pkg, fn),
				})
			}
		}
		return true
	})
}

// isUnrecoverableExitCall reports whether call terminates the program without
// running deferred functions: [os.Exit], [syscall.Exit] or one of the [log.Fatal]
// functions (which call [os.Exit]).
//
// panic and [log.Panic] are intentionally excluded: deferred functions still run
// while a panic unwinds the stack.
func isUnrecoverableExitCall(call *ast.CallExpr) (pkg, fn string, ok bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", "", false
	}

	pkg, fn = id.Name, sel.Sel.Name
	switch pkg {
	case "os", "syscall":
		return pkg, fn, fn == "Exit"
	case "log":
		switch fn {
		case "Fatal", "Fatalf", "Fatalln":
			return pkg, fn, true
		}
	}
	return "", "", false
}
