package rule

import (
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseWaitGroupGoRule spots Go idoms that might be rewritten using WaitGroup.Go.
type UseWaitGroupGoRule struct{}

// Apply applies the rule to given file.
func (*UseWaitGroupGoRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go125) {
		return nil // skip analysis if Go version < 1.25
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseWaitGroupGo{
		onFailure: onFailure,
	}

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Body == nil {
			continue // external (no-Go) function
		}

		ast.Walk(w, fn.Body)
		//w.analyzeBlock(fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*UseWaitGroupGoRule) Name() string {
	return "use-waitgroup-go"
}

type lintUseWaitGroupGo struct {
	onFailure func(lint.Failure)
}

var wgDonePicker = func(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	result := ok && astutils.IsPkgDotName(call.Fun, "wg", "Done")
	return result
}

func (w *lintUseWaitGroupGo) analyzeBlock(b *ast.BlockStmt) {
	// we will iterate over all statements in search for wg.Add()
	stmts := b.List
	for i := 0; i < len(stmts); i++ {
		stmt := stmts[i]

		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue // not an expression statements thus not a function call
		}

		// Lets check if the expression statement is a call to wg.Add
		call, ok := expr.X.(*ast.CallExpr)
		notACallToWgAdd := !ok || !astutils.IsPkgDotName(call.Fun, "wg", "Add")
		if notACallToWgAdd {
			continue
		}

		// Here we have identified a call to wg.Add
		// Let's iterate over the statements that follow the wg.Add
		// to see if there is a go statement that runs a goroutine with a wg.Done
		//
		// wg.Add is the i-th statement of block.List
		// we will iterate from the (i+1)-th statement up to the last statement of block.List
		for i++; i < len(stmts); i++ {
			stmt := stmts[i]
			// looking for a go statement
			goStmt, ok := stmt.(*ast.GoStmt)
			if !ok {
				continue // not a go statement
			}

			// here we found a the go statement
			// now let's check is the go statement is applied to a function literal that contains a wg.Done
			funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
			if !ok {
				continue // the go statements runs a function defined elsewhere
			}

			// here we found a go statement running a function literal
			// now we will look for a wg.Done inside the body of the function literal
			found := astutils.PickNodes(funcLit.Body, wgDonePicker) // TODO use SeekNode
			if len(found) > 0 {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       call,
					Category:   lint.FailureCategoryCodeStyle,
					Failure:    "replace wg.Add()...go {...wg.Done()...} with wg.Go(...)",
				})
				break
			}
		}
	}
}

// This visitor searches AST subtrees with the following form
// wg.Add(...)
// ...
//
//	go func (...) {
//	   ...
//	   wg.Done // or defer wg.Done
//	   ...
//	}
//
// Warning: the analysis only looks for exactly wg.Add and wg.Done, that means
// calls to Add and Done on a WaitGroup struct within a variable named differently than wg will be ignored
// This simplification avoids requiring type information while still makes the rule work in most of the cases.
// (Who names a wait group differently than wg ?!)
func (w *lintUseWaitGroupGo) Visit(node ast.Node) ast.Visitor {
	// Only interested in blocks of statements
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w // not a block of statements
	}

	// Once in a block of statements
	// we will iterate over all statements in search for wg.Add()
	for i := 0; i < len(block.List); i++ {
		stmt := block.List[i]
		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue // not an expression statements thus not a function call
		}

		// Lets check if the expression statement is a call to wg.Add
		call, ok := expr.X.(*ast.CallExpr)
		notACallToWgAdd := !ok || !astutils.IsPkgDotName(call.Fun, "wg", "Add")
		if notACallToWgAdd {
			continue
		}

		// Here we have identified a call to wg.Add
		// Let's iterate over the statements that follow the wg.Add
		// to see if there is a go statement that runs a goroutine with a wg.Done
		//
		// wg.Add is the i-th statement of block.List
		// we will iterate from the (i+1)-th statement up to the last statement of block.List
		for i++; i < len(block.List); i++ {
			stmt := block.List[i]
			// looking for a go statement
			goStmt, ok := stmt.(*ast.GoStmt)
			if !ok {
				continue // not a go statement
			}

			// here we found a the go statement
			// now let's check is the go statement is applied to a function literal that contains a wg.Done
			funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
			if !ok {
				continue // the go statements runs a function defined elsewhere
			}

			// here we found a go statement running a function literal
			// now we will look for a wg.Done inside the body of the function literal
			found := astutils.PickNodes(funcLit.Body, wgDonePicker) // TODO use SeekNode
			if len(found) > 0 {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       call,
					Category:   lint.FailureCategoryCodeStyle,
					Failure:    "replace wg.Add()...go {...wg.Done()...} with wg.Go(...)",
				})
			}

			return w
		}
	}

	return nil
}
