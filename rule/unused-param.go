package rule

import (
	"fmt"
	"go/ast"
	"go/token"

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

func (s *scope) addVars(exps []ast.Expr) {
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
	varSelector := func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		return ok && id.Obj != nil && id.Obj.Kind.String() == "var"
	}
	switch n := node.(type) {
	case *ast.BlockStmt:
		v.sStk.openScope()
		walkStmtList(v, n.List)
		v.sStk.closeScope()
		return nil
	case *ast.AssignStmt:
		var uses []ast.Node
		if isOpAssign(n.Tok) { // Case of id += expr
			uses = append(uses, pickFromExpList(n.Lhs, varSelector, nil)...)
		} else { // Case of id[expr] = expr
			indexSelector := func(n ast.Node) bool {
				_, ok := n.(*ast.IndexExpr)
				return ok
			}
			f := func(n ast.Node) []ast.Node {
				ie, ok := n.(*ast.IndexExpr)
				if !ok { // not possible
					return nil
				}

				return pick(ie.Index, varSelector, nil)
			}

			uses = append(uses, pickFromExpList(n.Lhs, indexSelector, f)...)
		}

		uses = append(uses, pickFromExpList(n.Rhs, varSelector, nil)...)

		markParamListAsUsed(uses, v)
		cs := v.sStk.currentScope()
		cs.addVars(n.Lhs)
	case *ast.Ident:
		if n.Obj != nil {
			if n.Obj.Kind.String() == "var" {
				markParamAsUsed(n, v)
			}
		}
	case *ast.ForStmt:
		v.sStk.openScope()
		if n.Init != nil {
			ast.Walk(v, n.Init)
		}
		uses := pickFromExpList([]ast.Expr{n.Cond}, varSelector, nil)
		markParamListAsUsed(uses, v)
		ast.Walk(v, n.Body)
		v.sStk.closeScope()
		return nil
	case *ast.SwitchStmt:
		v.sStk.openScope()
		if n.Init != nil {
			ast.Walk(v, n.Init)
		}
		uses := pickFromExpList([]ast.Expr{n.Tag}, varSelector, nil)
		markParamListAsUsed(uses, v)
		// Analyze cases (they are not BlockStmt but a list of Stmt)
		cases := n.Body.List
		for _, c := range cases {
			cc, ok := c.(*ast.CaseClause)
			if !ok {
				continue
			}
			uses := pickFromExpList(cc.List, varSelector, nil)
			markParamListAsUsed(uses, v)
			v.sStk.openScope()
			for _, stmt := range cc.Body {
				ast.Walk(v, stmt)
			}
			v.sStk.closeScope()
		}

		v.sStk.closeScope()
		return nil
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

func markParamListAsUsed(ids []ast.Node, v funcVisitor) {
	for _, id := range ids {
		markParamAsUsed(id.(*ast.Ident), v)
	}
}

func markParamAsUsed(id *ast.Ident, v funcVisitor) { // TODO: constraint parameters to receive just a list of params and a scope stack
	for _, s := range v.sStk.stk {
		if s.vars[id.Name] {
			return
		}
	}

	if v.params[id.Name] {
		v.params[id.Name] = false
	}
}

func isOpAssign(aTok token.Token) bool {
	return aTok == token.ADD_ASSIGN || aTok == token.AND_ASSIGN ||
		aTok == token.MUL_ASSIGN || aTok == token.OR_ASSIGN ||
		aTok == token.QUO_ASSIGN || aTok == token.REM_ASSIGN ||
		aTok == token.SHL_ASSIGN || aTok == token.SHR_ASSIGN ||
		aTok == token.SUB_ASSIGN || aTok == token.XOR_ASSIGN
}
