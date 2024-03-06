package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// EmptyBlockRule lints given else constructs.
type EmptyBlockRule struct{}

// Apply applies the rule to given file.
func (*EmptyBlockRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintEmptyBlock{make(map[*ast.BlockStmt]bool), onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*EmptyBlockRule) Name() string {
	return "empty-block"
}

type lintEmptyBlock struct {
	ignore    map[*ast.BlockStmt]bool
	onFailure func(lint.Failure)
}

func (w lintEmptyBlock) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.ignore[n.Body] = true
		return w
	case *ast.FuncLit:
		w.ignore[n.Body] = true
		return w
	case *ast.SelectStmt:
		w.ignore[n.Body] = true
		return w
	case *ast.ForStmt:
		if len(n.Body.List) == 0 && n.Init == nil && n.Post == nil && n.Cond != nil {
			if _, isCall := n.Cond.(*ast.CallExpr); isCall {
				w.ignore[n.Body] = true
				return w
			}
		}
	case *ast.RangeStmt:
		if len(n.Body.List) == 0 {
			if n.Key != nil || n.Value != nil || !isRangeOverChannel(n.X) {
				w.onFailure(lint.Failure{
					Confidence: 0.9,
					Node:       n,
					Category:   "logic",
					Failure:    "this block is empty, you can remove it",
				})
			}
			return nil // skip visiting the range subtree (it will produce a duplicated failure)
		}
	case *ast.BlockStmt:
		if !w.ignore[n] && len(n.List) == 0 {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   "logic",
				Failure:    "this block is empty, you can remove it",
			})
		}
	}

	return w
}

// isRangeOverChannel implements a best-effort detection for expressions in the context of 'for range'
// whose result is a channel.
// To do this it assumes that channels can only come from call return values, assignments or variables.
func isRangeOverChannel(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return isChanCallExpr(e, 0)
	case *ast.Ident:
		return isChanIdent(e)
	}
	return false
}

func isChanIdent(ident *ast.Ident) bool {
	var typ ast.Expr
	switch decl := ident.Obj.Decl.(type) {
	case *ast.Field:
		// function parameter
		typ = decl.Type.(*ast.ChanType)
	case *ast.AssignStmt:
		return isChanAssignStmt(decl, ident.Name)
	case *ast.ValueSpec:
		// defined variable (var x type)
		if decl.Type != nil {
			typ = decl.Type
			break
		}
		// assigned variable (var x = ...)
		idx := -1
		for i, n := range decl.Names {
			if n.Name == ident.Name {
				idx = i
				break
			}
		}
		if idx < 0 || idx+1 > len(decl.Values) {
			return false
		}
		var isChan bool
		switch expr := decl.Values[idx].(type) {
		case *ast.Ident:
			isChan = isChanIdent(expr)
		case *ast.CallExpr:
			isChan = isChanCallExpr(expr, idx)
		}
		if isChan {
			return true
		}
	}
	_, isChan := typ.(*ast.ChanType)
	return isChan
}

func isChanAssignStmt(assign *ast.AssignStmt, name string) bool {
	idx := -1
	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name != name {
			continue
		}
		idx = i
		break
	}
	if idx < 0 {
		return false
	}

	var rhs ast.Expr
	if len(assign.Lhs) == len(assign.Rhs) {
		// assignment with equal sides: a, b := c, d()
		rhs = assign.Rhs[idx]
	} else {
		// assignment with uneven sides: a, b := c()
		rhs = assign.Rhs[0]
	}

	switch expr := rhs.(type) {
	case *ast.CallExpr:
		return isChanCallExpr(expr, idx)
	case *ast.Ident:
		return isChanIdent(expr)
	}
	return false
}

func isChanCallExpr(call *ast.CallExpr, returnValueIndex int) bool {
	var fieldList *ast.FieldList
	var typ ast.Expr

	switch f := call.Fun.(type) {
	case *ast.FuncLit:
		// inline function definition and call
		fieldList = f.Type.Results
	case *ast.Ident:
		if f.Name == "make" {
			// special handling for make builtin
			typ = call.Args[0]
			break
		}
		// normal function call
		decl, ok := f.Obj.Decl.(*ast.FuncDecl)
		if !ok {
			break
		}
		fieldList = decl.Type.Results
	}
	if fieldList != nil && len(fieldList.List) > returnValueIndex {
		typ = fieldList.List[returnValueIndex].Type
	}
	_, isChan := typ.(*ast.ChanType)
	return isChan
}
