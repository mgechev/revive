package rule

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"

	"github.com/mgechev/revive/lint"
)

// AtomicRule lints given else constructs.
type AtomicRule struct{}

// Apply applies the rule to given file.
func (r *AtomicRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	walker := atomic{
		file: file,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (r *AtomicRule) Name() string {
	return "atomic"
}

type atomic struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w atomic) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	if len(n.Lhs) != len(n.Rhs) {
		return w
	}
	if len(n.Lhs) == 1 && n.Tok == token.DEFINE {
		return w
	}

	for i, right := range n.Rhs {
		call, ok := right.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkgIdent, _ := sel.X.(*ast.Ident)
		if w.file.Pkg.TypesInfo != nil {
			pkgName, ok := w.file.Pkg.TypesInfo.Uses[pkgIdent].(*types.PkgName)
			if !ok || pkgName.Imported().Path() != "sync/atomic" {
				continue
			}
		}

		switch sel.Sel.Name {
		case "AddInt32", "AddInt64", "AddUint32", "AddUint64", "AddUintptr":
			left := n.Lhs[i]
			if len(call.Args) != 2 {
				continue
			}
			arg := call.Args[0]
			broken := false

			if uarg, ok := arg.(*ast.UnaryExpr); ok && uarg.Op == token.AND {
				broken = w.gofmt(left) == w.gofmt(uarg.X)
			} else if star, ok := left.(*ast.StarExpr); ok {
				broken = w.gofmt(star.X) == w.gofmt(arg)
			}

			if broken {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Failure:    fmt.Sprintf("direct assignment to atomic value"),
					Node:       n,
				})
			}
		}
	}
	return w
}

// gofmt returns a string representation of the expression.
func (w atomic) gofmt(x ast.Expr) string {
	buf := bytes.Buffer{}
	fs := token.NewFileSet()
	printer.Fprint(&buf, fs, x)
	return buf.String()
}
