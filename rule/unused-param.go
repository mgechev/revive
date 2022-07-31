package rule

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/mgechev/revive/lint"
)

// UnusedParamRule lints unused params in functions.
type UnusedParamRule struct{}

// Apply applies the rule to given file.
func (*UnusedParamRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := lintUnusedParamRule{
		typesInfo: file.Pkg.TypesInfo(),
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UnusedParamRule) Name() string {
	return "unused-parameter"
}

type lintUnusedParamRule struct {
	typesInfo *types.Info
	onFailure func(lint.Failure)
}

func (w lintUnusedParamRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		params := retrieveNamedParams(n.Type.Params)
		if len(params) < 1 {
			return nil // skip, func without parameters
		}

		if n.Body == nil {
			return nil // skip, is a function prototype
		}

		structFields := map[*ast.Object]bool{}

		// inspect the func body looking for references to parameters
		// except struct field keys.
		fselect := func(n ast.Node) bool {
			if lit, ok := n.(*ast.CompositeLit); ok {
				isStruct := false

				switch lit.Type.(type) {
				case *ast.StructType:
					isStruct = true
				case *ast.Ident:
					_, isStruct = w.typesInfo.TypeOf(lit.Type).Underlying().(*types.Struct)
				}

				if isStruct {
					for _, e := range lit.Elts {
						if kv, ok := e.(*ast.KeyValueExpr); ok {
							if ident, ok := kv.Key.(*ast.Ident); ok {
								structFields[ident.Obj] = true
							}
						}
					}
				}
			}

			ident, isAnID := n.(*ast.Ident)
			if !isAnID {
				return false
			}

			if params[ident.Obj] && !structFields[ident.Obj] {
				params[ident.Obj] = false // mark as used
			}

			return false
		}
		_ = pick(n.Body, fselect, nil)

		for _, p := range n.Type.Params.List {
			for _, n := range p.Names {
				if params[n.Obj] {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Node:       n,
						Category:   "bad practice",
						Failure:    fmt.Sprintf("parameter '%s' seems to be unused, consider removing or renaming it as _", n.Name),
					})
				}
			}
		}

		return nil // full method body already inspected
	}

	return w
}

func retrieveNamedParams(params *ast.FieldList) map[*ast.Object]bool {
	result := map[*ast.Object]bool{}
	if params.List == nil {
		return result
	}

	for _, p := range params.List {
		for _, n := range p.Names {
			if n.Name == "_" {
				continue
			}

			result[n.Obj] = true
		}
	}

	return result
}
