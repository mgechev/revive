package rule

import (
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// WaitGroupDoneInWaitGroupGoRule spots Go idioms that might be rewritten using WaitGroup.Go.
type WaitGroupDoneInWaitGroupGoRule struct{}

// Apply applies the rule to given file.
func (*WaitGroupDoneInWaitGroupGoRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go125) {
		return nil // skip analysis if Go version < 1.25
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintWaitGroupDoneInWaitGroupGo{
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
func (*WaitGroupDoneInWaitGroupGoRule) Name() string {
	return "waitgroup-done-in-waitgroup-go"
}

type lintWaitGroupDoneInWaitGroupGo struct {
	onFailure func(lint.Failure)
}

func (w *lintWaitGroupDoneInWaitGroupGo) Visit(node ast.Node) ast.Visitor {
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

	// search a wg.Done in the body of the function literal
	wgDone := astutils.SeekNode[*ast.CallExpr](funcLit.Body, wgDonePicker)
	if wgDone == nil {
		return nil // there is no a call to wg.Done in the call to wg.Do
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       wgDone,
		Category:   lint.FailureCategoryErrors,
		Failure:    "do not call wg.Done inside wg.Go",
	})

	return nil
}
