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
		fv := funcVisitor{params: retrieveNamedParams(n.Type.Params.List)}
		if n.Body != nil {
			ast.Walk(fv, n.Body)
			checkUnusedParams(w, fv.params, n)
		}
		return nil
	}

	return w
}

type funcVisitor struct {
	params map[string]bool
}

func (v funcVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		if n.Obj != nil {
			if n.Obj.Kind.String() == "var" {
				markParamAsUsed(n, v.params)
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
func markParamAsUsed(id *ast.Ident, params map[string]bool) {
	if params[id.Name] {
		params[id.Name] = false
	}
}
