package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// UnusedParamRule lints unused params in functions.
type UnusedParamRule struct{}

// Apply applies the rule to given file.
func (r *UnusedParamRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintUnusedParamRule{onFailure: onFailure}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (r *UnusedParamRule) Name() string {
	return "unused-parameter"
}

type lintUnusedParamRule struct {
	onFailure func(lint.Failure)
}

func (w lintUnusedParamRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		fv := newFuncVisitor(retrieveNamedParams(n.Type.Params.List))
		if n.Body != nil {
			ast.Walk(fv, n.Body)
			checkUnusedParams(w, fv.params, n)
		}
		return nil
	}

	return w
}

type scope struct {
	vars map[string]bool
}

func newScope() scope {
	return scope{make(map[string]bool, 0)}
}

func (s *scope) addVar(exps []ast.Expr) {
	for _, e := range exps {
		if id, ok := e.(*ast.Ident); ok {
			s.vars[id.Name] = true
		}
	}
}

type scopeStack struct {
	stk []scope
}

func (s *scopeStack) openScope() {
	s.stk = append(s.stk, newScope())
}

func (s *scopeStack) closeScope() {
	if len(s.stk) > 0 {
		s.stk = s.stk[:len(s.stk)-1]
	}
}

func (s *scopeStack) currentScope() scope {
	if len(s.stk) > 0 {
		return s.stk[len(s.stk)-1]
	}

	panic("no current scope")
}

func newScopeStack() scopeStack {
	return scopeStack{make([]scope, 0)}
}

type funcVisitor struct {
	sStk   scopeStack
	params map[string]bool
}

func newFuncVisitor(params map[string]bool) funcVisitor {
	return funcVisitor{sStk: newScopeStack(), params: params}
}

func walkStmtList(v ast.Visitor, list []ast.Stmt) {
	for _, s := range list {
		ast.Walk(v, s)
	}
}

func (v funcVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BlockStmt:
		v.sStk.openScope()
		walkStmtList(v, n.List)
		v.sStk.closeScope()
		return nil
	case *ast.AssignStmt:
		varSelector := func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			return ok && id.Obj != nil && id.Obj.Kind.String() == "var"
		}
		uses := pickFromExpList(n.Rhs, varSelector)
		for _, id := range uses {
			markParamAsUsed(id.(*ast.Ident), v)
		}
		cs := v.sStk.currentScope()
		cs.addVar(n.Lhs)
	case *ast.Ident:
		if n.Obj != nil {
			if n.Obj.Kind.String() == "var" {
				markParamAsUsed(n, v)
			}
		}
	}

	return v
}

func retrieveNamedParams(pl []*ast.Field) map[string]bool {
	result := make(map[string]bool, len(pl))
	for _, p := range pl {
		for _, n := range p.Names {
			if n.Name != "_" {
				result[n.Name] = true
			}
		}
	}
	return result
}

func checkUnusedParams(w lintUnusedParamRule, params map[string]bool, n *ast.FuncDecl) {
	for k, v := range params {
		if v {
			w.onFailure(lint.Failure{
				Confidence: 0.8, // confidence is not 1.0 because of shadow variables
				Node:       n,
				Category:   "bad practice",
				Failure:    fmt.Sprintf("parameter '%s' seems to be unused, consider removing or renaming it as _", k),
			})
		}
	}

}
func markParamAsUsed(id *ast.Ident, v funcVisitor) {
	for _, s := range v.sStk.stk {
		if s.vars[id.Name] {
			return
		}
	}

	if v.params[id.Name] {
		v.params[id.Name] = false
	}
}

type picker struct {
	fselect  func(n ast.Node) bool
	onSelect func(n ast.Node)
}

func pick(n ast.Node, fselect func(n ast.Node) bool) []interface{} {
	var result []interface{}
	onSelect := func(n ast.Node) {
		result = append(result, n)
	}
	p := picker{fselect: fselect, onSelect: onSelect}
	ast.Walk(p, n)
	return result
}

func pickFromExpList(l []ast.Expr, fselect func(n ast.Node) bool) []interface{} {
	result := make([]interface{}, 0)
	for _, e := range l {
		result = append(result, pick(e, fselect)...)
	}
	return result
}

func (p picker) Visit(node ast.Node) ast.Visitor {
	if p.fselect == nil {
		return nil
	}

	if p.fselect(node) {
		p.onSelect(node)
	}

	return p
}
