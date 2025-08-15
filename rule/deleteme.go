package rule

import (
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// DeletemeRule is a sandbox rule to tests ideas
type DeletemeRule struct{}

// Apply applies the rule to given file.
func (r *DeletemeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeletemeRule{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*DeletemeRule) Name() string {
	return "deleteme"
}

type lintDeletemeRule struct {
	onFailure func(lint.Failure)
}

func (w *lintDeletemeRule) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	for i, stmt := range block.List {
		//println(">>>> block stmt:", astutils.GoFmt(stmt))
		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := expr.X.(*ast.CallExpr)
		notACallToWgAdd := !ok || !astutils.IsPkgDotName(call.Fun, "wg", "Add")
		if notACallToWgAdd {
			//println(">>>> CONTINUE main")
			continue
		}
		//println(">>>> FOUND wg.Add lets check next statements")
		for i++; i < len(block.List); i++ {
			stmt := block.List[i]
			//println(">>>> next stmt:", astutils.GoFmt(stmt))
			goStmt, ok := stmt.(*ast.GoStmt)
			if !ok {
				//println(">>>> CONTINUE not a go stmt")
				continue
			}
			funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
			//println(">>>> funclit:", astutils.GoFmt(funcLit))
			if !ok {
				continue
			}
			picker := func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				result := ok && astutils.IsPkgDotName(call.Fun, "wg", "Done")
				//println(">>>> picking on:", astutils.GoFmt(n), result)
				return result
			}

			found := astutils.PickNodes(funcLit.Body, picker)
			if len(found) > 0 {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       call,
					Category:   lint.FailureCategoryCodeStyle,
					Failure:    "replace wg.Add()...go {...wg.Done()...} with wg.Go(...)",
				})
			}
		}
	}

	return w
}
