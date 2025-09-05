package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// ForbiddenCallInWgGoRule spots calls to panic or wg.Done when using WaitGroup.Go.
type ForbiddenCallInWgGoRule struct{}

// Apply applies the rule to given file.
func (*ForbiddenCallInWgGoRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go125) {
		return nil // skip analysis if Go version < 1.25
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintForbiddenCallInWgGo{
		onFailure: onFailure,
	}

	// Iterate over declarations looking for function declarations
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue // not a function
		}

		if fn.Body == nil {
			continue // external (no-Go) function
		}

		// Analyze the function body
		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*ForbiddenCallInWgGoRule) Name() string {
	return "forbidden-call-in-wg-go"
}

type lintForbiddenCallInWgGo struct {
	onFailure func(lint.Failure)
}

func (w *lintForbiddenCallInWgGo) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w // not a call of statements
	}

	if !astutils.IsPkgDotName(call.Fun, "wg", "Go") {
		return w // not a call to wg.Go
	}

	if len(call.Args) != 1 {
		return nil // no argument (impossible)
	}

	funcLit, ok := call.Args[0].(*ast.FuncLit)
	if !ok {
		return nil // the argument is not a function literal
	}

	// search a wg.Done or panic in the body of the function literal
	wgDoneOrPanic := astutils.SeekNode[*ast.CallExpr](funcLit.Body, w.wgDoneOrPanicPicker)
	if wgDoneOrPanic == nil {
		return nil // there is no a call to wg.Done or panic in the call to wg.Do
	}

	msg := fmt.Sprintf("do not call %s inside wg.Go", w.getCallee(wgDoneOrPanic))
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       wgDoneOrPanic,
		Category:   lint.FailureCategoryErrors,
		Failure:    msg,
	})

	return nil
}

// getCallee yields the name of the callee of the given call expression.
// Hypothesis: the call expression is either a call to panic or a call to wg.Done.
func (w *lintForbiddenCallInWgGo) getCallee(callToWgDoneOrPanic *ast.CallExpr) string {
	_, ok := callToWgDoneOrPanic.Fun.(*ast.Ident)
	if !ok {
		return "wg.Done"
	}

	return "panic"
}

// function used when calling astutils.SeekNode that search for calls to wg.Done.
func (*lintForbiddenCallInWgGo) wgDoneOrPanicPicker(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	return astutils.IsIdent(call.Fun, "panic") ||
		astutils.IsPkgDotName(call.Fun, "wg", "Done")
}
