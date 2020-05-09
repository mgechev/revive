package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// UnconditionalRecursionRule lints given else constructs.
type UnconditionalRecursionRule struct{}

// Apply applies the rule to given file.
func (r *UnconditionalRecursionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	// err := file.Pkg.TypeCheck()
	// if err != nil {
	// 	panic(fmt.Sprintf("Error while type-checking file %s: %v", file.Name, err))
	// }
	w := lintUnconditionalRecursionRule{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *UnconditionalRecursionRule) Name() string {
	return "unconditional-recursion"
}

type funcDesc struct {
	reciverId *ast.Ident
	id        *ast.Ident
}

func (fd *funcDesc) equal(other *funcDesc) bool {
	receiversAreEqual := (fd.reciverId == nil && other.reciverId == nil) || fd.reciverId != nil && other.reciverId != nil && fd.reciverId.Name == other.reciverId.Name
	idsAreEqual := (fd.id == nil && other.id == nil) || fd.id.Name == other.id.Name

	return receiversAreEqual && idsAreEqual
}

type funcStatus struct {
	funcDesc            *funcDesc
	seenConditionalExit bool
}

type lintUnconditionalRecursionRule struct {
	onFailure   func(lint.Failure)
	currentFunc *funcStatus
}

func (w lintUnconditionalRecursionRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt:
		w.updateFuncStatus(n.Body)
		w.updateFuncStatus(n.Else)
		return nil
	case *ast.SelectStmt:
		w.updateFuncStatus(n.Body)
		return nil
	case *ast.RangeStmt:
		w.updateFuncStatus(n.Body)
		return nil
	case *ast.TypeSwitchStmt:
		w.updateFuncStatus(n.Body)
		return nil
	case *ast.SwitchStmt:
		w.updateFuncStatus(n.Body)
		return nil
	case *ast.GoStmt:
		for _, a := range n.Call.Args {
			ast.Walk(w, a) // check if arguments have a recursive call
		}
		return nil // recursive async call is not an issue
	case *ast.ForStmt:
		if n.Cond != nil {
			return nil
		}
		// unconditional loop
		return w
	case *ast.FuncDecl:
		var rec *ast.Ident
		switch {
		case n.Recv == nil || n.Recv.NumFields() < 1 || len(n.Recv.List[0].Names) < 1:
			rec = nil
		default:
			rec = n.Recv.List[0].Names[0]
		}

		w.currentFunc = &funcStatus{&funcDesc{rec, n.Name}, false}
	case *ast.CallExpr:
		var funcId *ast.Ident
		var selector *ast.Ident
		switch c := n.Fun.(type) {
		case *ast.Ident:
			selector = nil
			funcId = c
		case *ast.SelectorExpr:
			var ok bool
			selector, ok = c.X.(*ast.Ident)
			if !ok { // a.b....Foo()
				return nil
			}
			funcId = c.Sel
		default:
			return w
		}

		if w.currentFunc != nil && // not in a func body
			!w.currentFunc.seenConditionalExit && // there is a conditional exit in the function
			w.currentFunc.funcDesc.equal(&funcDesc{selector, funcId}) {
			w.onFailure(lint.Failure{
				Category:   "logic",
				Confidence: 1,
				Node:       n,
				Failure:    "unconditional recursive call",
			})
		}
	}

	return w
}

func (w *lintUnconditionalRecursionRule) updateFuncStatus(node ast.Node) {
	if node == nil || w.currentFunc == nil || w.currentFunc.seenConditionalExit {
		return
	}

	w.currentFunc.seenConditionalExit = w.hasControlExit(node)
}

func (w *lintUnconditionalRecursionRule) hasControlExit(node ast.Node) bool {
	filter := func(n ast.Node) bool {
		_, ok := n.(*ast.ReturnStmt)
		return ok
	}

	return len(pick(node, filter, nil)) != 0
}
