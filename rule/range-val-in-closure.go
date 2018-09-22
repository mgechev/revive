package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// RangeValInClosureRule looks for iteration vars used in closures.
type RangeValInClosureRule struct{}

// Apply applies the rule to given file.
func (r *RangeValInClosureRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintRangeValInClosureRule{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *RangeValInClosureRule) Name() string {
	return "range-val-in-closure"
}

type lintRangeValInClosureRule struct {
	params    map[string]bool
	onFailure func(lint.Failure)
}

func (lintRangeValInClosureRule) retrieveParamNames(pl []*ast.Field) map[string]bool {
	result := make(map[string]bool, len(pl))
	for _, p := range pl {
		for _, n := range p.Names {
			result[n.Name] = true
		}
	}
	return result
}

// checkId checks if the given id is referenced in a closure inside the given statement block
func (w lintRangeValInClosureRule) checkId(id string, blk *ast.BlockStmt) {
	if id == "_" {
		return
	}

	for _, stmt := range blk.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		if len(assign.Lhs) == 1 && len(assign.Rhs) == 1 && isIdent(assign.Lhs[0], id) && isIdent(assign.Rhs[0], id) {
			return // the range body has an assignment of the form id := id, thus using it in a closure is safe.
		}
	}

	fselect := func(n ast.Node) bool { // picks go statements
		_, ok := n.(*ast.GoStmt)
		return ok
	}

	goStmts := pick(blk, fselect, nil)

GoStmtIter:
	for _, gs := range goStmts {
		gs, _ := gs.(*ast.GoStmt)

		cf := gs.Call.Fun
		fLit, ok := cf.(*ast.FuncLit)
		if !ok {
			continue
		}

		// check if the range value (id) is passed as argument
		for _, arg := range gs.Call.Args {
			ident, ok := arg.(*ast.Ident)
			if ok {
				if ident.Name == id {
					continue GoStmtIter
				}
			}
		}

		fselect := func(n ast.Node) bool { // picks reference to the range value (id)
			ident, ok := n.(*ast.Ident)
			return ok && ident.Name == id
		}

		ref2Id := pick(fLit.Body, fselect, nil)

		if len(ref2Id) > 0 {
			w.onFailure(lint.Failure{
				Confidence: 0.8,
				Node:       fLit,
				Category:   "logic",
				Failure:    fmt.Sprintf("range value '%s' seems to be referenced inside the closure", id),
			})
		}
	}
}

func (w lintRangeValInClosureRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.RangeStmt:
		// check the range value
		ident, ok := n.Value.(*ast.Ident)
		if ok {
			w.checkId(ident.Name, n.Body)
		}

		// check the range key
		ident, ok = n.Key.(*ast.Ident)
		if ok {
			w.checkId(ident.Name, n.Body)
		}
	}

	return w
}
